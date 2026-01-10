import { toast } from 'react-toastify';
import { Button } from '@heroui/react';
import { DeleteRequest, globalFetch } from '../../Request';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { WarningIcon } from '../../Icons/WarningIcon';

export default function AbortMenu({ gameId }: { gameId: number | null }) {
  const navigate = useNavigate();
  const [isAborting, setIsAborting] = useState(false);

  const handleAbort = async () => {
    if (!gameId) return;

    const confirmed = window.confirm(
      'Are you sure you want to abort this game? This action cannot be undone.'
    );
    if (!confirmed) return;

    setIsAborting(true);
    const abortRequest = DeleteRequest(`/api/game/${gameId}/abort`);
    const data = await globalFetch<{}>(abortRequest);

    if (!data) {
      setIsAborting(false);
      return;
    }

    setIsAborting(false);
    toast.warn('Game aborted');

    navigate('/');
  };

  return (
    <div className="mt-auto pt-4 flex justify-center">
      <Button
        color="danger"
        variant="solid"
        onPress={handleAbort}
        isLoading={isAborting}
        isDisabled={!gameId || isAborting}
        className="w-1/2"
        startContent={<WarningIcon className="text-white" />}
      >
        {isAborting ? 'Aborting...' : 'Abort Game'}
      </Button>
    </div>
  );
}
