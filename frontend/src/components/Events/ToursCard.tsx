import { useTranslation } from 'react-i18next';
import { Card, useOverlayState } from '@heroui/react';
import type { TourState } from '../../models/Tour';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { cancelTour, publishTour } from '../../providers/event/tournamentApi';
import Confirm from '../Confirm';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

// WIND_KANJI: the starting seats as the table shows them.
const WIND_KANJI = { east: '東', south: '南', west: '西', north: '北' };

// ToursCard is the tours part of the tournament's desk: the current tour and
// its tables, and for a referee publishing the next tour's seating (random
// for now) or cancelling a published one before it starts.
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
  const { isPending: isPublishing, run: runPublish } = useSubmitOnce();
  const { isPending: isCancelling, run: runCancel } = useSubmitOnce();
  const cancelConfirm = useOverlayState();
  const current = state.current;
  const isPublished = current?.status === 'published';

  const publish = () =>
    runPublish(async () => {
      if (await publishTour(eventId, 'random')) {
        showToast({
          title: t('tours.published', { n: state.next }),
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

        {canRun && state.next > 0 && (
          <>
            <p className="text-sm text-muted">
              {blocker ??
                t('tours.publishHint', {
                  n: state.next,
                  players: state.active,
                })}
            </p>
            <PendingButton
              variant="primary"
              isPending={isPublishing}
              isDisabled={blocker != null}
              onPress={publish}
              className="w-3/4 self-center"
            >
              {t('tours.publish', { n: state.next })}
            </PendingButton>
          </>
        )}

        {canRun && isPublished && (
          <PendingButton
            variant="danger-soft"
            isPending={isCancelling}
            onPress={cancelConfirm.open}
            className="w-3/4 self-center"
          >
            {t('tours.cancel')}
          </PendingButton>
        )}

        {state.tables.length > 0 && (
          <ul aria-label={t('tours.tables')} className="flex flex-col gap-3">
            {state.tables.map((table) => (
              <li key={table.table} className="rounded-xl bg-default/50 p-2">
                <p className="mb-1 text-sm font-medium">
                  {t('tours.table', { n: table.table })}
                </p>
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
              </li>
            ))}
          </ul>
        )}
      </Card.Content>
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
