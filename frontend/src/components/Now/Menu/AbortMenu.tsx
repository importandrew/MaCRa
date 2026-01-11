import { addToast } from '@heroui/react';
import { Button, useDisclosure } from '@heroui/react';
import { DeleteRequest, globalFetch } from '../../Request';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';
import Confirm from '../../Confirm';

export default function AbortMenu({ gameId }: { gameId: number | null }) {
  const navigate = useNavigate();
  const [isAborting, setIsAborting] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleAbort = async () => {
    if (!gameId) return;

    setIsAborting(true);
    const abortRequest = DeleteRequest(`/api/game/${gameId}/abort`);
    const data = await globalFetch<null>(abortRequest);

    if (!data) {
      setIsAborting(false);
      return;
    }

    setIsAborting(false);
    addToast({
      title: 'Game aborted',
      color: 'warning',
    });

    navigate('/');
  };

  return (
    <>
      <div className="mt-auto pt-4 flex justify-center">
        <Button
          color="danger"
          variant="solid"
          onPress={onOpen}
          isDisabled={!gameId || isAborting}
          className="w-1/2"
          startContent={<WarningIcon className="text-white" />}
        >
          {isAborting ? 'Aborting...' : 'Abort Game'}
        </Button>
      </div>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          handleAbort();
        }}
        isOpen={isOpen}
        headerText="Confirm abort game"
        cancelText="Cancel"
        confirmText="Abort Game"
        prompt="Are you sure you want to abort this game? This action cannot be undone."
      />
    </>
  );
}
