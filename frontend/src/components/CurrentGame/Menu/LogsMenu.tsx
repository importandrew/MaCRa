import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import { useState } from 'react';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import GameLogs from '../../GameResult/GameLogs';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';
import { UndoIcon } from '../../Icons/UndoIcon';
import {
  GameLogsState,
  undoLastRound,
  useFetchGameLogsData,
} from '../../../providers/game/gameDataApi';
import PendingButton from '../../PendingButton';
import { Outcome } from '../../../models/GameData';
import { GameRoundLogs } from '../../../models/GameResult';

// LogsMenu shows a game's round log with the Undo button. A table screen
// passes its own log and size (tableLogs, tableLogSize from useTableLogs), so
// the menu shows and undoes the state the table shows; without them (the
// result page), the menu loads the log itself.
export default function LogsMenu({
  gameId,
  eventId,
  canUndoLastRound = false,
  onUndoSuccess,
  tableLogs,
  tableLogSize,
  onStale,
}: {
  gameId: number | null;
  eventId?: number;
  canUndoLastRound?: boolean;
  onUndoSuccess?: () => void;
  tableLogs?: GameLogsState;
  tableLogSize?: number;
  // onStale reloads the table screen after an undo refused because the log
  // has changed.
  onStale?: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const { isPending: undoLoading, run } = useSubmitOnce();
  const ownLogs = useFetchGameLogsData(tableLogs ? null : gameId, eventId);
  const logs = tableLogs ?? ownLogs;
  const logSize = tableLogSize ?? logs.data?.round_logs.length;

  // Undo sends the log size the referee sees, so the server refuses it (409)
  // if someone else has entered or undone an entry since; the table screen
  // (onStale) and the log then reload, and the message says what changed, as
  // for a round entry.
  // A 409 with the log unchanged means the final round can't be undone: a
  // player has started another game (the server's message names them). A
  // second tap is ignored.
  const handleUndoConfirm = () =>
    run(async () => {
      if (gameId == null || logSize == null) return;
      const seenLastId = lastEntryId(logs.data);
      setConfirmOpen(false);
      const { status, error } = await undoLastRound(gameId, logSize, eventId);
      if (status === 200) {
        showToast({
          title: t('game.lastRoundUndoneToast'),
          variant: 'warning',
        });
        // A table screen reloads its log with the table.
        if (!tableLogs) void logs.fetchData();
        onUndoSuccess?.();
        return;
      }
      if (status !== 409) return;

      // The table first, so the log loaded next matches it. An undo and a
      // new entry can bring the log back to the same size, so compare the
      // last entry too.
      await onStale?.();
      const fresh = await logs.fetchData();
      const isChanged =
        fresh != null &&
        (fresh.round_logs.length !== logSize ||
          lastEntryId(fresh) !== seenLastId);
      if (fresh && isChanged) {
        showToast({
          title: changeTitle(fresh, logSize, t),
          description: t('game.checkLogAndRetry'),
          variant: 'warning',
        });
        return;
      }
      showToast({
        title: t('game.cannotUndoFinal'),
        description: error,
        variant: 'warning',
      });
    });

  return (
    <div className="flex flex-col justify-center mb-3 gap-2">
      <GameLogs
        isChartVisible={false}
        isLogVisible={true}
        gameId={gameId}
        eventId={eventId}
        logs={logs}
      ></GameLogs>
      {canUndoLastRound && (
        <>
          <PendingButton
            className="w-1/2 self-center"
            variant="tertiary"
            size="sm"
            onPress={() => setConfirmOpen(true)}
            isPending={undoLoading}
            isDisabled={gameId == null || !logSize}

            icon={<UndoIcon />}
          >
            {t('game.undoLastRound')}
          </PendingButton>
          <Confirm
            isOpen={confirmOpen}
            onClose={() => setConfirmOpen(false)}
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

function lastEntryId(log: GameRoundLogs | null): number | undefined {
  return log?.round_logs[log.round_logs.length - 1]?.id;
}
