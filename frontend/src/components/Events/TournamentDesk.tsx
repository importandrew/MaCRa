import { useTranslation } from 'react-i18next';
import { Button, Card, Chip, useOverlayState } from '@heroui/react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { eventStatus } from '../../models/Event';
import type { TournamentParticipant } from '../../models/TournamentSettings';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useEventPage } from '../../providers/event/eventLayoutContext';
import {
  drawNumbers,
  useFetchParticipants,
  useFetchTours,
  withdrawParticipant,
} from '../../providers/event/tournamentApi';
import ToursCard from './ToursCard';
import Confirm from '../Confirm';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

// TournamentDesk is /event/:eventId/tournament, the tournament's desk for its
// staff: the participants with their numbers and, for a referee, the draw
// (before the first tour starts) and withdrawals (after). Tours come here
// later.
export default function TournamentDesk() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event, hasCapability, canView } =
    useEventPage('view_tournament');
  const canRun = hasCapability('run_tournament');
  const { data, isLoading, isError, isForbidden, refetch } =
    useFetchParticipants(canView ? eventId : null);
  const tours = useFetchTours(canView ? eventId : null);
  // An action on the desk may change both: e.g. a draw lets a tour be
  // published.
  const refetchAll = () => Promise.all([refetch(), tours.refetch()]);
  const { isPending: isDrawing, run: runDraw } = useSubmitOnce();
  const { run: runWithdraw } = useSubmitOnce();
  const drawConfirm = useOverlayState();
  const withdrawConfirm = useOverlayState();
  const [toWithdraw, setToWithdraw] = useState<TournamentParticipant | null>(
    null
  );

  const status = event ? eventStatus(event) : null;
  const participants = data?.participants ?? [];
  // A published seating goes by the numbers: it is cancelled before a new
  // draw.
  const isSeated = tours.data?.current?.status === 'published';
  const canDraw =
    canRun && status === 'preparation' && participants.length > 0 && !isSeated;
  const canWithdraw = canRun && status === 'in_progress';

  const draw = () =>
    runDraw(async () => {
      drawConfirm.close();
      if (await drawNumbers(eventId)) {
        showToast({ title: t('tournamentDesk.drawn'), variant: 'success' });
      }
      await refetchAll();
    });

  const withdraw = () =>
    runWithdraw(async () => {
      withdrawConfirm.close();
      if (!toWithdraw) return;
      if (await withdrawParticipant(eventId, toWithdraw.user_id)) {
        showToast({
          title: t('tournamentDesk.withdrawn', { name: toWithdraw.name }),
          variant: 'success',
        });
      }
      await refetch();
    });

  return (
    <LoadingTemplate
      isLoading={isLoading || tours.isLoading}
      isError={isError || tours.isError}
      isForbidden={!canView || isForbidden}
      onRefresh={refetchAll}
      onGoHome={() => navigate(`/event/${eventId}`)}
    >
      {event && data && (
        <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
          <h1 className="text-2xl font-semibold truncate min-w-0">
            {t('tournamentDesk.title', { name: event.name })}
          </h1>
          {tours.data && (
            <ToursCard
              eventId={eventId}
              state={tours.data}
              isDrawn={data.drawn}
              canRun={canRun}
              onChanged={refetchAll}
            />
          )}
          <Card>
            <Card.Content className="gap-3">
              <h2 className="font-medium">
                {t('tournamentDesk.participants', {
                  n: participants.length,
                })}
              </h2>
              <p className="text-sm text-muted">
                {participants.length === 0
                  ? t('tournamentDesk.noPlayers')
                  : status !== 'preparation'
                    ? t('tournamentDesk.closedHint')
                    : data.drawn
                      ? t('tournamentDesk.drawnHint')
                      : t('tournamentDesk.notDrawnHint')}
              </p>
              {canDraw && (
                <PendingButton
                  variant={data.drawn ? 'secondary' : 'primary'}
                  isPending={isDrawing}
                  // A redraw replaces everyone's numbers: confirm it.
                  onPress={data.drawn ? drawConfirm.open : draw}
                  className="w-3/4 self-center"
                >
                  {data.drawn
                    ? t('tournamentDesk.redraw')
                    : t('tournamentDesk.draw')}
                </PendingButton>
              )}
              {participants.length > 0 && (
                <ul aria-label={t('tournamentDesk.participantsList')}>
                  {participants.map((participant) => (
                    <li
                      key={participant.user_id}
                      className="flex items-center gap-3 border-b border-separator px-2 py-2 last:border-b-0"
                    >
                      <span className="w-8 shrink-0 text-end tabular-nums text-muted">
                        {participant.number ?? '—'}
                      </span>
                      <span
                        className={`min-w-0 flex-1 truncate ${
                          participant.withdrawn ? 'text-muted' : ''
                        }`}
                      >
                        {participant.name}
                      </span>
                      {participant.withdrawn ? (
                        <Chip size="sm" variant="soft">
                          {t('tournamentDesk.withdrewMark')}
                        </Chip>
                      ) : (
                        canWithdraw && (
                          <Button
                            size="sm"
                            variant="danger-soft"
                            onPress={() => {
                              setToWithdraw(participant);
                              withdrawConfirm.open();
                            }}
                          >
                            {t('tournamentDesk.withdraw')}
                          </Button>
                        )
                      )}
                    </li>
                  ))}
                </ul>
              )}
            </Card.Content>
          </Card>
        </div>
      )}
      <Confirm
        isOpen={drawConfirm.isOpen}
        onClose={drawConfirm.close}
        onConfirm={draw}
        headerText={t('tournamentDesk.redrawTitle')}
        prompt={t('tournamentDesk.redrawPrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('tournamentDesk.redraw')}
        confirmVariant="primary"
      />
      <Confirm
        isOpen={withdrawConfirm.isOpen}
        onClose={withdrawConfirm.close}
        onConfirm={withdraw}
        headerText={t('tournamentDesk.withdrawTitle')}
        prompt={t('tournamentDesk.withdrawPrompt', {
          name: toWithdraw?.name ?? '',
        })}
        cancelText={t('common.cancel')}
        confirmText={t('tournamentDesk.withdraw')}
        confirmVariant="danger"
      />
    </LoadingTemplate>
  );
}
