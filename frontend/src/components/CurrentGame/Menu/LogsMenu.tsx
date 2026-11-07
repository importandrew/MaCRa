import { useState } from 'react';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import GameLogs from '../../GameResult/GameLogs';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';
import { UndoIcon } from '../../Icons/UndoIcon';
import { undoLastRound } from '../../../providers/game/gameDataApi';
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

  // A second Undo tap would remove one more round, so run it once.
  const handleUndoConfirm = () =>
    run(async () => {
      if (gameId == null) return;
      setConfirmOpen(false);
      const ok = await undoLastRound(gameId, eventId);
      if (ok) {
        showToast({ title: 'Last round undone', variant: 'warning' });
        onUndoSuccess?.();
      }
    });

  return (
    <div className="flex flex-col justify-center mb-3 gap-2">
      <GameLogs
        isChartVisible={false}
        isLogVisible={true}
        gameId={gameId}
        eventId={eventId}
      ></GameLogs>
      {canUndoLastRound && (
        <>
          <PendingButton
            className="w-1/2 self-center"
            variant="tertiary"
            size="sm"
            onPress={() => setConfirmOpen(true)}
            isPending={undoLoading}
            isDisabled={gameId == null}

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
