import { useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useFetchGameStatus } from '../../providers/game/gameDataApi';
import GameResult from '../GameResult/GameResult';
import Live from '../Live/Live';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

export default function ViewGame() {
  const navigate = useNavigate();
  const { eventId: eventIdParam, gameId: gameIdParam } = useParams<{
    eventId: string;
    gameId: string;
  }>();
  const eventId = eventIdParam ? parseInt(eventIdParam, 10) : 0;
  const gameId = gameIdParam ? parseInt(gameIdParam, 10) : 0;

  const { status, isLoading, isError, isForbidden, refetch } =
    useFetchGameStatus(eventId, gameId);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  if (!isLoading && !isError && !isForbidden) {
    if (status === 'live') {
      return <Live />;
    }
    if (status === 'ended') {
      return <GameResult />;
    }
  }

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate('/')}
    >
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-3 px-3">
        <p className="text-center text-3xl">Game not available</p>
        <p className="text-center text-default-500">
          This game is not live or finished.
        </p>
      </div>
    </LoadingTemplate>
  );
}
