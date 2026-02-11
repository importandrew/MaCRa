import { Button, useDisclosure } from '@heroui/react';
import { PostRequest, globalFetch } from '../../Request';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';
import { showToast } from '../../Toast';

export default function FinishMenu({ gameId }: { gameId: number | null }) {
  const navigate = useNavigate();
  const [isFinishing, setIsFinishing] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleFinish = async () => {
    if (!gameId) return;

    setIsFinishing(true);
    const finishRequest = PostRequest(`/api/game/${gameId}/finish`);
    const data = await globalFetch<null>(finishRequest);

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
        prompt="Are you sure you want to finish this game? This action will finalize the result."
      />
    </>
  );
}

