import { useTranslation } from 'react-i18next';
import { Button, Card, useOverlayState } from '@heroui/react';
import { useState } from 'react';
import {
  SEATING_METHODS,
  type SeatingMethod,
  type TourState,
  type TourTable,
} from '../../models/Tour';
import SegmentedPicker from '../SegmentedPicker';
import ListStepper from '../ListStepper';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  cancelTour,
  generateTour,
  tourAction,
} from '../../providers/event/tournamentApi';
import Confirm from '../Confirm';
import AddTimeModal from './AddTimeModal';
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

  // hasPlan: a seating is prepared for all tours (the "fixed" method).
  const hasPlan = state.plan != null;
  // picked: the seating method the referee picked; null until they do.
  const [picked, setPicked] = useState<SeatingMethod | null>(null);
  // shownMethod: the method the next generation uses: the one picked, else
  // the prepared seating when it has the next tour, else random (also when
  // "fixed" was picked and the plan is gone since).
  const shownMethod: SeatingMethod =
    picked === 'fixed' && !hasPlan
      ? 'random'
      : (picked ??
        (state.plan?.tours.includes(state.next) ? 'fixed' : 'random'));
  // step: an interval seating's step, 1 to the number of tables; null
  // means the number of tables (the default).
  const tableCount = Math.max(1, Math.floor(state.active / 4));
  const [step, setStep] = useState<number | null>(null);
  const methodPicker = (
    <>
      <SegmentedPicker
        ariaLabel={t('tours.methodLabel')}
        // Four one-word labels don't fit a row at 320 px in Russian.
        orientation="vertical"
        options={SEATING_METHODS.filter(
          (value) => value !== 'fixed' || hasPlan
        ).map((value) => ({
          value,
          label: t(`tours.method.${value}`),
        }))}
        selectedKey={shownMethod}
        onSelectionChange={setPicked}
      />
      {shownMethod === 'interval' && tableCount > 1 && (
        <>
          <p className="text-sm">{t('tours.step')}</p>
          <ListStepper
            label={t('tours.step')}
            values={Array.from({ length: tableCount }, (_, i) => i + 1)}
            value={step ?? tableCount}
            onChange={setStep}
          />
        </>
      )}
    </>
  );

  const generate = () =>
    runGenerate(async () => {
      if (
        await generateTour(
          eventId,
          shownMethod,
          shownMethod === 'interval' && step != null && step !== tableCount
            ? step
            : undefined
        )
      ) {
        showToast({
          title: t('tours.generated', { n: state.next }),
          variant: 'success',
        });
      }
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

  // Why the next tour can't be published yet, if it can't.
  const blocker = !isDrawn
    ? t('tours.needsDraw')
    : state.missing > 0
      ? t('tours.missing', { count: state.missing })
      : null;

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
              {current.seating_step
                ? t('tours.methodWithStep', {
                    method: t(`tours.method.${current.seating}`),
                    step: current.seating_step,
                  })
                : t('tours.methodUsed', {
                    method: t(`tours.method.${current.seating}`),
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
            <p className="text-sm text-muted">
              {blocker ??
                t(`tours.methodHint.${shownMethod}`, {
                  players: state.active,
                  n: state.next,
                })}
            </p>
            <PendingButton
              variant="primary"
              isPending={isGenerating}
              isDisabled={blocker != null}
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
                {canRun && isLive && table.status === 'live' && (
                  <Button
                    size="sm"
                    variant="tertiary"
                    className="float-end"
                    onPress={() =>
                      setTimeFor({ gameId: table.game_id, number: table.table })
                    }
                  >
                    {t('tours.addTime')}
                  </Button>
                )}
                <ul>
                  {table.seats.map((seat) => (
                    <li key={seat.user_id} className="flex items-center gap-2">
                      <span className="w-5 shrink-0 text-center">
                        {WIND_KANJI[seat.wind]}
                      </span>
                      <span className="w-7 shrink-0 text-end tabular-nums text-muted">
                        {seat.number}
                      </span>
                      <span className="min-w-0 truncate">{seat.name}</span>
                    </li>
                  ))}
                </ul>
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
