import { Button } from '@heroui/react';
import { useState } from 'react';
import GameLogs from '../../GameResults/GameLogs';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';
import { UndoIcon } from '../../Icons/UndoIcon';
import { undoLastRound } from '../../../providers/game/gameDataApi';

export default function LogsMenu({
  gameId,
  canUndoLastRound = false,
  onUndoSuccess,
}: {
  gameId: number | null;
  canUndoLastRound?: boolean;
  onUndoSuccess?: () => void;
}) {
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [undoLoading, setUndoLoading] = useState(false);

  const handleUndoConfirm = async () => {
    if (gameId == null) return;
    setConfirmOpen(false);
    setUndoLoading(true);
    const ok = await undoLastRound(gameId);
    setUndoLoading(false);
    if (ok) {
      showToast({ title: 'Last round undone', color: 'warning' });
      onUndoSuccess?.();
    }
  };

  return (
    <div className="flex flex-col justify-center mb-3 gap-2">
      <GameLogs
        isChartVisible={false}
        isLogVisible={true}
        gameId={gameId}
      ></GameLogs>
      {canUndoLastRound && (
        <>
          <Button
            className="w-1/2 self-center"
            color="default"
            variant="flat"
            size="sm"
            startContent={<UndoIcon />}
            onPress={() => setConfirmOpen(true)}
            isLoading={undoLoading}
            isDisabled={gameId == null}
          >
            Undo last round
          </Button>
          <Confirm
            isOpen={confirmOpen}
            onClose={() => setConfirmOpen(false)}
            onConfirm={handleUndoConfirm}
            prompt="Do you want to undo the last round? This cannot be undone."
            headerText="Undo last round"
            cancelText="Cancel"
            confirmText="Undo"
            confirmButtonColor="danger"
          />
        </>
      )}
    </div>
  );
}
