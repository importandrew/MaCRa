import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import { useOverlayState } from '@heroui/react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import GameLogs, { LogsState } from '../GameResult/GameLogs';
import Confirm from '../Confirm';
import { showToast } from '../Toast';
import { UndoIcon } from '../Icons/UndoIcon';
import { undoLastRound } from '../../providers/game/gameDataApi';
import PendingButton from '../PendingButton';
import { Outcome } from '../../models/GameData';
import { apiErrorText } from '../apiErrorText';
import { GameRoundLogs } from '../../models/GameResult';

// LogsMenu shows a game's round log with the Undo button. The caller passes
// the log: a table screen the one from its table response, so the menu shows
// and undoes the state the table shows; the result page the one it loaded.
export default function LogsMenu({
  gameId,
  eventId,
  canUndoLastRound = false,
  onUndoSuccess,
  logs,
  onStale,
}: {
  gameId: number | null;
  eventId: number | null;
  canUndoLastRound?: boolean;
  onUndoSuccess?: () => void;
  logs: LogsState;
  // onStale reloads the screen after an undo refused because the log has
  // changed, and returns the fresh log (null if it can't be loaded).
  onStale: () => Promise<GameRoundLogs | null>;
}) {
  const { t } = useTranslation();
  const confirm = useOverlayState();
  const { isPending: undoLoading, run } = useSubmitOnce();
  const logSize = logs.data?.round_logs.length;

  // Undo sends the log size the referee sees, so the server refuses it
  // (log_changed) if someone else has entered or undone an entry since; the
  // screen then reloads, and the message says what changed, as for a round
  // entry. Any other refusal (e.g. the final round can't be undone because a
  // player has started another game) is told as is. A second tap is ignored.
  const handleUndoConfirm = () =>
    run(async () => {
      if (eventId == null || gameId == null || logSize == null) return;
      confirm.close();
      const { status, error } = await undoLastRound(eventId, gameId, logSize);
      if (status === 200) {
        showToast({
          title: t('game.lastRoundUndoneToast'),
          variant: 'success',
        });
        onUndoSuccess?.();
        return;
      }
      if (status !== 409) return;

      if (error?.error === 'log_changed') {
        const fresh = await onStale();
        showToast({
          title: fresh
            ? changeTitle(fresh, logSize, t)
            : t('apiErrors.log_changed'),
          description: t('game.checkLogAndRetry'),
          variant: 'warning',
        });
        return;
      }
      showToast({
        title: t('game.cannotUndoFinal'),
        description: apiErrorText(error, t),
        variant: 'warning',
      });
    });

  return (
    <div className="flex flex-col justify-center mb-3 gap-2">
      <GameLogs isChartVisible={false} isLogVisible={true} logs={logs} />
      {canUndoLastRound && (
        <>
          <PendingButton
            className="min-w-1/2 self-center"
            variant="tertiary"
            size="sm"
            onPress={confirm.open}
            isPending={undoLoading}
            isDisabled={gameId == null || !logSize}

            icon={<UndoIcon />}
          >
            {t('game.undoLastRound')}
          </PendingButton>
          <Confirm
            isOpen={confirm.isOpen}
            onClose={confirm.close}
            onConfirm={handleUndoConfirm}
            prompt={t('game.undoConfirmPrompt')}
            headerText={t('game.undoLastRound')}
            cancelText={t('common.cancel')}
            confirmText={t('game.undo')}
            confirmVariant="danger"
          />
        </>
      )}
    </div>
  );
}

// changeTitle says what changed in the log since the referee saw logSize
// entries, in the words the round entry uses.
function changeTitle(
  fresh: GameRoundLogs,
  logSize: number,
  t: TFunction
): string {
  const entries = fresh.round_logs;
  if (entries.length < logSize) {
    return t('game.lastEntryAlreadyUndone');
  }
  const last = entries[entries.length - 1];
  const isPenalty =
    last.outcome === Outcome.Penalty || last.outcome === Outcome.FalseMahjong;
  return isPenalty ? t('game.penaltyJustEntered') : t('game.roundJustEntered');
}
