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

  // Undo sends the size of the log shown above, so the server refuses it if
  // someone else has entered or undone an entry since (409, the server's
  // message is shown); the log then reloads. A second tap is ignored too.
  const handleUndoConfirm = () =>
    run(async () => {
      if (gameId == null || logSize == null) return;
      setConfirmOpen(false);
      const status = await undoLastRound(gameId, logSize, eventId);
      if (status === 200) {
        showToast({ title: 'Last round undone', variant: 'warning' });
        void logs.fetchData();
        onUndoSuccess?.();
        return;
      }
      if (status === 409) {
        void logs.fetchData();
      }
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
