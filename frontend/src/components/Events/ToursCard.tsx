import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  ToggleButton,
  ToggleButtonGroup,
  useOverlayState,
} from '@heroui/react';
import { useState } from 'react';
import type { Seating, TourState, TourTable } from '../../models/Tour';
import SeatingFields from './SeatingFields';
import { seatingLabel } from './seatingLabel';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  cancelTour,
  generateTour,
  tourAction,
} from '../../providers/event/tournamentApi';
import Confirm from '../Confirm';
import AddTimeModal from './AddTimeModal';
import SeatModal from './SeatModal';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import { getTimeFormatted } from '../../utils/date-tools';

// seatName is the name of a player at the table.
function seatName(table: TourTable, userId: number): string {
  return table.seats.find((seat) => seat.user_id === userId)?.name ?? '?';
}

// WIND_KANJI: the starting seats as the table shows them.
const WIND_KANJI = { east: '東', south: '南', west: '西', north: '北' };

// ToursCard is the tours part of the tournament's desk: the current tour and
// its tables, and for a referee generating the next tour's seating by the
// method picked (or from the seating prepared for all tours), publishing it,
// and running the tour.
export default function ToursCard({
  eventId,
  state,
  isDrawn,
  canRun,
  onChanged,
}: {
  eventId: number;
  state: TourState;
  // isDrawn: the participant numbers are drawn for the current players.
  isDrawn: boolean;
  canRun: boolean;
  onChanged: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const { isPending: isGenerating, run: runGenerate } = useSubmitOnce();
  const { isPending: isCancelling, run: runCancel } = useSubmitOnce();
  const cancelConfirm = useOverlayState();
  const current = state.current;
  const isDraft = current?.status === 'draft';
  const isPublished = current?.status === 'published';
  const isLive = current?.status === 'live';
  // playing: the live tour's tables whose game hasn't ended.
  const playing = state.tables.filter(
    (table) => table.status !== 'ended'
  ).length;
  const { isPending: isActing, run: runAct } = useSubmitOnce();
  const finishConfirm = useOverlayState();
  const abortConfirm = useOverlayState();

  // tourEnds: when the live tour's time is up, before time added to a
  // table.
  const tourEnds =
    isLive && current?.started_at && current.minutes
      ? new Date(
          new Date(current.started_at).getTime() + current.minutes * 60_000
        ).toISOString()
      : null;
  // timeFor: the table a referee adds time to (the modal is open).
  // seatFor: the seat a referee changes (the modal is open).
  const [seatFor, setSeatFor] = useState<{
    table: TourTable;
    seat: TourTable['seats'][number];
  } | null>(null);
  // canSeat: a referee changes who sits at a table of a published or live
  // tour whose game isn't over.
  const canSeat = (table: TourTable) =>
    canRun &&
    (isPublished || isLive) &&
    (table.status === 'created' || table.status === 'live');
  const [timeFor, setTimeFor] = useState<{
    gameId: number;
    number: number;
  } | null>(null);

  const act = (action: 'publish' | 'start' | 'finish' | 'abort') =>
    runAct(async () => {
      finishConfirm.close();
      abortConfirm.close();
      if (await tourAction(eventId, action)) {
        showToast({
          title: t(`tours.done.${action}`, { n: current?.number ?? 0 }),
          variant: 'success',
        });
      }
      await onChanged();
    });

  // nextPlan: the plan of the tour to seat next, if any.
  const nextPlan = state.plan.find((plan) => plan.number === state.next);
  const hasFile = (nextPlan?.file_tables ?? 0) > 0;
  const tableCount = Math.max(1, Math.floor(state.active / 4));
  // picked: the seating the referee picked for the next tour; null until
  // they do: the plan's then (fixed only while the tour has tables from a
  // file).
  const [picked, setPicked] = useState<Seating | null>(null);
  const planned: Seating = nextPlan
    ? {
        seating: nextPlan.seating,
        step: nextPlan.seating_step ?? undefined,
        avoid_repeats: nextPlan.avoid_repeats,
      }
    : { seating: 'random', avoid_repeats: false };
  const chosen = picked ?? planned;
  const seating: Seating =
    chosen.seating === 'fixed' && !hasFile
      ? { seating: 'random', avoid_repeats: false }
      : chosen;
  const methodPicker = (
    <SeatingFields
      value={seating}
      onChange={setPicked}
      tableCount={tableCount}
      hasFile={hasFile}
    />
  );

  const generate = () =>
    runGenerate(async () => {
      const step =
        seating.seating === 'interval' &&
        seating.step != null &&
        seating.step < tableCount
          ? seating.step
          : undefined;
      if (
        await generateTour(eventId, { ...seating, step }, chosenSubstitutes)
      ) {
        showToast({
          title: t('tours.generated', { n: state.next }),
          variant: 'success',
        });
      }
      // The next tour starts from its plan again.
      setPicked(null);
      setPickedSubstitutes(null);
      await onChanged();
    });

  const cancel = () =>
    runCancel(async () => {
      cancelConfirm.close();
      if (await cancelTour(eventId)) {
        showToast({ title: t('tours.cancelled'), variant: 'success' });
      }
      await onChanged();
    });

  // Substitutes complete the tables when players are missing: those who
  // sit nowhere, or at a table of the draft (generating again replaces it).
  const inDraft = new Set(
    isDraft
      ? state.tables.flatMap((table) => table.seats.map((s) => s.user_id))
      : []
  );
  const freeSubstitutes = state.substitutes.filter(
    (sub) => sub.table === 0 || inDraft.has(sub.id)
  );
  // pickedSubstitutes: the referee's choice; null until they choose: the
  // draft's substitutes then.
  const [pickedSubstitutes, setPickedSubstitutes] = useState<number[] | null>(
    null
  );
  const chosenSubstitutes = (
    pickedSubstitutes ??
    freeSubstitutes.map((sub) => sub.id).filter((id) => inDraft.has(id))
  )
    .filter((id) => freeSubstitutes.some((sub) => sub.id === id))
    .slice(0, state.missing);

  // Why the next tour can't be seated yet, if it can't: no draw, or fewer
  // players and free substitutes than whole tables need.
  const blocker = !isDrawn
    ? t('tours.needsDraw')
    : state.missing > freeSubstitutes.length
      ? t('tours.missing', { count: state.missing - freeSubstitutes.length })
      : null;
  // Still to pick: substitutes for the missing players.
  const toPick = blocker == null ? state.missing - chosenSubstitutes.length : 0;
  const substitutePicker = blocker == null && state.missing > 0 && (
    <>
      <p className="text-sm">{t('tours.substitutesLabel')}</p>
      <ToggleButtonGroup
        aria-label={t('tours.substitutesLabel')}
        selectionMode="multiple"
        orientation="vertical"
        fullWidth
        className="w-full"
        selectedKeys={chosenSubstitutes.map(String)}
        onSelectionChange={(keys) =>
          setPickedSubstitutes([...keys].map(Number).slice(0, state.missing))
        }
      >
        {freeSubstitutes.map((sub, i) => (
          <ToggleButton
            key={sub.id}
            id={String(sub.id)}
            isDisabled={toPick === 0 && !chosenSubstitutes.includes(sub.id)}
            className="h-auto min-h-11 w-full py-2 not-data-selected:text-foreground"
          >
            {i > 0 && <ToggleButtonGroup.Separator />}
            {sub.name}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
      <p className={`text-sm ${toPick > 0 ? 'text-warning' : 'text-muted'}`}>
        {toPick > 0
          ? t('tours.pickSubstitutes', { count: toPick })
          : t('tours.substitutesHint')}
      </p>
    </>
  );

  return (
    <Card>
      <Card.Content className="gap-3">
        <h2 className="font-medium">
          {current
            ? t('tours.current', {
                n: current.number,
                status: t(`tours.status.${current.status}`),
              })
            : t('tours.title')}
        </h2>
        {current && current.status !== 'finished' && (
          <div className="text-sm">
            <p>
              {t('tours.methodUsed', {
                method: seatingLabel(
                  t,
                  current.seating,
                  current.seating_step,
                  current.avoid_repeats
                ),
              })}
            </p>
            <p className={state.repeats > 0 ? 'text-warning' : 'text-muted'}>
              {state.repeats > 0
                ? t('tours.repeatsSummary', {
                    count: state.repeats,
                    maxMet: state.max_met,
                    players: state.players_with_repeats,
                  })
                : t('tours.noRepeats')}
            </p>
          </div>
        )}

        {canRun && state.next > 0 && !isDraft && (
          <>
            {blocker == null && methodPicker}
            {substitutePicker}
            <p className="text-sm text-muted">
              {blocker ??
                t(`tours.methodHint.${seating.seating}`, {
                  players: state.active,
                  n: state.next,
                })}
            </p>
            <PendingButton
              variant="primary"
              isPending={isGenerating}
              isDisabled={blocker != null || toPick > 0}
              onPress={generate}
              className="w-3/4 self-center"
            >
              {t('tours.generate', { n: state.next })}
            </PendingButton>
          </>
        )}

        {canRun && isDraft && (
          <>
            <p className="text-sm text-muted">{t('tours.draftHint')}</p>
            {methodPicker}
            {substitutePicker}
            <PendingButton
              variant="primary"
              isPending={isActing}
              onPress={() => act('publish')}
              className="w-3/4 self-center"
            >
              {t('tours.publish')}
            </PendingButton>
            <PendingButton
              variant="secondary"
              isPending={isGenerating}
              isDisabled={blocker != null || toPick > 0}
              onPress={generate}
              className="w-3/4 self-center"
            >
              {t('tours.regenerate')}
            </PendingButton>
            <PendingButton
              variant="danger-soft"
              isPending={isCancelling}
              onPress={cancelConfirm.open}
              className="w-3/4 self-center"
            >
              {t('tours.cancel')}
            </PendingButton>
          </>
        )}

        {canRun && isPublished && (
          <>
            <PendingButton
              variant="primary"
              isPending={isActing}
              onPress={() => act('start')}
              className="w-3/4 self-center"
            >
              {t('tours.start', { n: current.number })}
            </PendingButton>
            <PendingButton
              variant="danger-soft"
              isPending={isCancelling}
              onPress={cancelConfirm.open}
              className="w-3/4 self-center"
            >
              {t('tours.cancel')}
            </PendingButton>
          </>
        )}

        {canRun && isLive && (
          <>
            {tourEnds && (
              <p className="text-sm">
                {t('tours.endsAt', { time: getTimeFormatted(tourEnds) })}
              </p>
            )}
            <p className="text-sm text-muted">
              {playing > 0
                ? t('tours.playing', { count: playing })
                : t('tours.allEnded')}
            </p>
            <PendingButton
              variant="primary"
              isPending={isActing}
              isDisabled={playing > 0}
              onPress={finishConfirm.open}
              className="w-3/4 self-center"
            >
              {t('tours.finish', { n: current.number })}
            </PendingButton>
            <PendingButton
              variant="danger-soft"
              isPending={isActing}
              onPress={abortConfirm.open}
              className="w-3/4 self-center"
            >
              {t('tours.abort')}
            </PendingButton>
          </>
        )}

        {state.tables.length > 0 && (
          <ul aria-label={t('tours.tables')} className="flex flex-col gap-3">
            {state.tables.map((table) => (
              <li key={table.table} className="rounded-xl bg-default/50 p-2">
                <p className="mb-1 flex justify-between gap-2 text-sm font-medium">
                  <span>
                    {t('tours.table', { n: table.table })}
                    {table.places && (
                      <span className="font-normal text-muted">
                        {' · '}
                        {t('tours.places', {
                          from: table.places.from,
                          to: table.places.to,
                        })}
                      </span>
                    )}
                  </span>
                  {isLive && (
                    <span className="font-normal text-muted">
                      {t(`tours.tableStatus.${table.status}`)}
                      {table.extra_minutes > 0 &&
                        ` · ${t('tours.extra', { n: table.extra_minutes })}`}
                    </span>
                  )}
                </p>
                <ul>
                  {table.seats.map((seat) => (
                    <li key={seat.user_id} className="flex items-center gap-2">
                      <span className="w-5 shrink-0 text-center">
                        {WIND_KANJI[seat.wind]}
                      </span>
                      <span className="w-7 shrink-0 text-end tabular-nums text-muted">
                        {seat.substitute ? '—' : seat.number}
                      </span>
                      <span className="min-w-0 flex-1">
                        <span className="block truncate">
                          {seat.sitter ? seat.sitter.name : seat.name}
                        </span>
                        {/* A substitute's seat: whose it is, and whether the
                            player gets 0 / 0. */}
                        {(seat.sitter || seat.forfeit) && (
                          <span className="block truncate text-xs text-muted">
                            {[
                              seat.sitter && t('seat.for', { name: seat.name }),
                              // A substitute's own seat: "substitute", not
                              // "0 / 0" (it counts for nobody anyway).
                              seat.substitute
                                ? t('seat.substitute')
                                : seat.forfeit && t('seat.forfeitMark'),
                            ]
                              .filter(Boolean)
                              .join(' · ')}
                          </span>
                        )}
                      </span>
                      {canSeat(table) && (
                        <Button
                          size="sm"
                          variant="tertiary"
                          aria-label={t('seat.change', { name: seat.name })}
                          onPress={() => setSeatFor({ table, seat })}
                        >
                          {t('seat.changeShort')}
                        </Button>
                      )}
                    </li>
                  ))}
                </ul>
                {/* Under the seats, so a long name at the first one isn't
                    cut by it. */}
                {canRun && isLive && table.status === 'live' && (
                  <Button
                    size="sm"
                    variant="tertiary"
                    className="mt-1 ml-auto flex"
                    onPress={() =>
                      setTimeFor({ gameId: table.game_id, number: table.table })
                    }
                  >
                    {t('tours.addTime')}
                  </Button>
                )}
                {table.repeats > 0 && (
                  <p className="mt-1 text-sm text-warning">
                    {t('tours.metBefore', {
                      pairs: table.repeat_pairs
                        .map((pair) =>
                          t('tours.pair', {
                            a: seatName(table, pair.a),
                            b: seatName(table, pair.b),
                            count: pair.met,
                          })
                        )
                        .join(', '),
                    })}
                  </p>
                )}
              </li>
            ))}
          </ul>
        )}
      </Card.Content>
      <SeatModal
        eventId={eventId}
        target={seatFor}
        substitutes={state.substitutes}
        onClose={() => setSeatFor(null)}
        onDone={onChanged}
      />
      <AddTimeModal
        eventId={eventId}
        table={timeFor}
        onClose={() => setTimeFor(null)}
        onAdded={onChanged}
      />
      <Confirm
        isOpen={finishConfirm.isOpen}
        onClose={finishConfirm.close}
        onConfirm={() => act('finish')}
        headerText={t('tours.finishTitle')}
        prompt={t('tours.finishPrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('tours.finish', { n: current?.number ?? 0 })}
        confirmVariant="primary"
      />
      <Confirm
        isOpen={abortConfirm.isOpen}
        onClose={abortConfirm.close}
        onConfirm={() => act('abort')}
        headerText={t('tours.abortTitle')}
        prompt={t('tours.abortPrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('tours.abort')}
        confirmVariant="danger"
      />
      <Confirm
        isOpen={cancelConfirm.isOpen}
        onClose={cancelConfirm.close}
        onConfirm={cancel}
        headerText={t('tours.cancelTitle')}
        prompt={t('tours.cancelPrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('tours.cancel')}
        confirmVariant="danger"
      />
    </Card>
  );
}
