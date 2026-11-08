import { Button, useOverlayState } from '@heroui/react';
import { DeleteRequest, authFetch, eventGamePath } from '../../Request';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';

export default function AbortMenu({
  gameId,
  eventId,
  canAbort = false,
}: {
  gameId: number | null;
  // eventId is the event of the shown game: the URL's on the live view, the
  // selected one at the table.
  eventId: number | null;
  canAbort?: boolean;
}) {
  const navigate = useNavigate();
  const { isPending: isAborting, run } = useSubmitOnce();
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  if (!canAbort) {
    return null;
  }

  const handleAbort = async () => {
    if (!gameId || !eventId) return;

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
        isDisabled={!gameId || !eventId || isAborting}
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
