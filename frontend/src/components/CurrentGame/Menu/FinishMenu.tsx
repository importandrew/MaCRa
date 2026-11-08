import { Button, useOverlayState } from '@heroui/react';
import { PostRequest, authFetch, eventGamePath } from '../../Request';
import { useSubmitOnce } from '../../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';

export default function FinishMenu({
  gameId,
  eventId,
  canFinish = false,
}: {
  gameId: number | null;
  // eventId is the event of the shown game: the URL's on the live view, the
  // selected one at the table.
  eventId: number | null;
  canFinish?: boolean;
}) {
  const navigate = useNavigate();
  const { isPending: isFinishing, run } = useSubmitOnce();
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  const handleFinish = async () => {
    if (!gameId || !eventId) return;

    const finishRequest = PostRequest(
      eventGamePath(eventId, gameId, '/finish')
    );
    const data = await authFetch<null>(finishRequest);

    if (!data) return;

    showToast({
      title: 'Game finished',
      variant: 'warning',
    });

    navigate('/last');
  };

  if (!canFinish) {
    return null;
  }

  return (
    <>
      <Button
        variant="primary"
        onPress={onOpen}
        isDisabled={!gameId || !eventId || isFinishing}
        className="w-1/2"
      >
        <WarningIcon className="text-white" />
        {isFinishing ? 'Finishing...' : 'Finish Game'}
      </Button>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          void run(handleFinish);
        }}
        isOpen={isOpen}
        headerText="Confirm finish game"
        cancelText="Cancel"
        confirmText="Finish Game"
        prompt="Are you sure you want to finish this game? Game will be counted. This action will finalize the result."
        confirmVariant="primary"
      />
    </>
  );
}
