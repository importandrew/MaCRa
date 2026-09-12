import { Button, useDisclosure } from '@heroui/react';
import { DeleteRequest, authFetch, eventGamePath } from '../../Request';
import { getStoredEventId } from '../../../providers/event/eventConstants';
import { useState } from 'react';
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
  const [isAborting, setIsAborting] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  if (!canAbort) {
    return null;
  }

  const handleAbort = async () => {
    if (!gameId) return;

    setIsAborting(true);
    const eventId = getStoredEventId();
    if (!eventId) {
      setIsAborting(false);
      return;
    }

    const abortRequest = DeleteRequest(
      eventGamePath(eventId, gameId, '/abort')
    );
    const data = await authFetch<null>(abortRequest);

    if (!data) {
      setIsAborting(false);
      return;
    }

    setIsAborting(false);
    showToast({
      title: 'Game aborted',
      color: 'warning',
    });

    navigate('/');
  };

  return (
    <>
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
        prompt="Are you sure you want to abort this game? Game will NOT be counted. This action cannot be undone."
        confirmButtonColor="danger"
      />
    </>
  );
}
