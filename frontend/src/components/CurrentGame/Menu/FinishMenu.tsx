import { Button, useDisclosure } from '@heroui/react';
import { PostRequest, authFetch, eventGamePath } from '../../Request';
import { getStoredEventId } from '../../../providers/event/eventConstants';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';

export default function FinishMenu({
  gameId,
  canFinish = false,
}: {
  gameId: number | null;
  canFinish?: boolean;
}) {
  const navigate = useNavigate();
  const [isFinishing, setIsFinishing] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleFinish = async () => {
    if (!gameId) return;

    const eventId = getStoredEventId();
    if (!eventId) return;

    setIsFinishing(true);
    const finishRequest = PostRequest(
      eventGamePath(eventId, gameId, '/finish')
    );
    const data = await authFetch<null>(finishRequest);

    if (!data) {
      setIsFinishing(false);
      return;
    }

    setIsFinishing(false);
    showToast({
      title: 'Game finished',
      color: 'warning',
    });

    navigate('/last');
  };

  if (!canFinish) {
    return null;
  }

  return (
    <>
      <Button
        color="primary"
        variant="solid"
        onPress={onOpen}
        isDisabled={!gameId || isFinishing}
        className="w-1/2"
        startContent={<WarningIcon className="text-white" />}
      >
        {isFinishing ? 'Finishing...' : 'Finish Game'}
      </Button>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          handleFinish();
        }}
        isOpen={isOpen}
        headerText="Confirm finish game"
        cancelText="Cancel"
        confirmText="Finish Game"
        prompt="Are you sure you want to finish this game? Game will be counted. This action will finalize the result."
        confirmButtonColor="primary"
      />
    </>
  );
}
