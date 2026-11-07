import { useState } from 'react';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import GameLogs from '../../GameResult/GameLogs';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';
import { UndoIcon } from '../../Icons/UndoIcon';
import {
  undoLastRound,
  useFetchGameLogsData,
} from '../../../providers/game/gameDataApi';
import PendingButton from '../../PendingButton';
import { Outcome } from '../../../models/GameData';
import { GameRoundLogs } from '../../../models/GameResult';

export default function LogsMenu({
  gameId,
  eventId,
  canUndoLastRound = false,
  onUndoSuccess,
}: {
  gameId: number | null;
  eventId?: number;
  canUndoLastRound?: boolean;
  onUndoSuccess?: () => void;
}) {
  const [confirmOpen, setConfirmOpen] = useState(false);
  const { isPending: undoLoading, run } = useSubmitOnce();
  const logs = useFetchGameLogsData(gameId, eventId);
  const logSize = logs.data?.round_logs.length;

  // Undo sends the size of the log shown above, so the server refuses it
  // (409) if someone else has entered or undone an entry since; the log then
  // reloads and the message says what changed, as for a round entry. A 409
  // with the log unchanged means the final round can't be undone: a player
  // has started another game (the server's message names them). A second
  // tap is ignored.
  const handleUndoConfirm = () =>
    run(async () => {
      if (gameId == null || logSize == null) return;
      const seenLastId = lastEntryId(logs.data);
      setConfirmOpen(false);
      const { status, error } = await undoLastRound(gameId, logSize, eventId);
      if (status === 200) {
        showToast({ title: 'Last round undone', variant: 'warning' });
        void logs.fetchData();
        onUndoSuccess?.();
        return;
      }
      if (status !== 409) return;

      // An undo and a new entry can bring the log back to the same size, so
      // compare the last entry too.
      const fresh = await logs.fetchData();
      const isChanged =
        fresh != null &&
        (fresh.round_logs.length !== logSize ||
          lastEntryId(fresh) !== seenLastId);
      if (fresh && isChanged) {
        showToast({
          title: changeTitle(fresh, logSize),
          description: 'Check the log and undo again if needed.',
          variant: 'warning',
        });
        return;
      }
      showToast({
        title: 'Cannot undo the final round',
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
            Undo last round
          </PendingButton>
          <Confirm
            isOpen={confirmOpen}
            onClose={() => setConfirmOpen(false)}
            onConfirm={handleUndoConfirm}
            prompt="Do you want to undo the last round? This cannot be undone."
            headerText="Undo last round"
            cancelText="Cancel"
            confirmText="Undo"
            confirmVariant="danger"
          />
        </>
      )}
    </div>
  );
}

// changeTitle says what changed in the log since the referee saw logSize
// entries, in the words the round entry uses.
function changeTitle(fresh: GameRoundLogs, logSize: number): string {
  const entries = fresh.round_logs;
  if (entries.length < logSize) {
    return 'The last entry has already been undone';
  }
  const last = entries[entries.length - 1];
  const isPenalty =
    last.outcome === Outcome.Penalty || last.outcome === Outcome.FalseMahjong;
  return isPenalty
    ? 'A penalty has just been entered'
    : 'A round has just been entered';
}

function lastEntryId(log: GameRoundLogs | null): number | undefined {
  return log?.round_logs[log.round_logs.length - 1]?.id;
}
