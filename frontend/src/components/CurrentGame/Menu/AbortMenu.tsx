import { Button, useOverlayState } from '@heroui/react';
import { DeleteRequest, authFetch, eventGamePath } from '../../Request';
import { getStoredEventId } from '../../../providers/event/eventConstants';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';

export default function AbortMenu({
  gameId,
  canAbort = false,
}: {
  gameId: number | null;
  canAbort?: boolean;
}) {
  const navigate = useNavigate();
  const { isPending: isAborting, run } = useSubmitOnce();
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  if (!canAbort) {
    return null;
  }

  const handleAbort = async () => {
    if (!gameId) return;

    const eventId = getStoredEventId();
    if (!eventId) return;

    const abortRequest = DeleteRequest(
      eventGamePath(eventId, gameId, '/abort')
    );
    const data = await authFetch<null>(abortRequest);

    if (!data) return;

    showToast({
      title: 'Game aborted',
      variant: 'warning',
    });

    navigate('/');
  };

  return (
    <>
      <Button
        variant="danger"
        onPress={onOpen}
        isDisabled={!gameId || isAborting}
        className="w-1/2"
      >
        <WarningIcon className="text-white" />
        {isAborting ? 'Aborting...' : 'Abort Game'}
      </Button>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          void run(handleAbort);
        }}
        isOpen={isOpen}
        headerText="Confirm abort game"
        cancelText="Cancel"
        confirmText="Abort Game"
        prompt="Are you sure you want to abort this game? Game will NOT be counted. This action cannot be undone."
        confirmVariant="danger"
      />
    </>
  );
}
