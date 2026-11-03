import { Button } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useContext, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { getStoredEventId } from '../../providers/event/eventConstants';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import { useFetchGameResultData } from '../../providers/game/gameDataApi';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import LogsMenu from '../CurrentGame/Menu/LogsMenu';
import GameLogs from './GameLogs';
import GameResultTable from './GameResultTable';

export default function GameResult() {
  const navigate = useNavigate();
  const { isAuthenticated } = useContext(AuthContext);
  const { eventId: eventIdParam, gameId: gameIdParam } = useParams<{
    eventId?: string;
    gameId?: string;
  }>();

  const isEventGameRoute = Boolean(eventIdParam && gameIdParam);
  const eventId = eventIdParam ? parseInt(eventIdParam, 10) : undefined;
  const gameId = gameIdParam ? parseInt(gameIdParam, 10) : undefined;
  const storedEventId = getStoredEventId();
  const resolvedEventId = eventId ?? storedEventId;
  const { hasCapability } = useEventCapabilities(resolvedEventId);
  const canCreateGame = hasCapability('create_game');

  const { data, isError, isForbidden, isLoading, fetchData } =
    useFetchGameResultData(
      isEventGameRoute && eventId && gameId ? { eventId, gameId } : undefined
    );
  const canUndoLastRound =
    hasCapability('undo_round') && Boolean(data?.can_undo_last_round);
  const [isChartVisible, setIsChartVisible] = useState(false);
  const [isLogsVisible, setIsLogsVisible] = useState(false);

  const showChart = () => {
    setIsChartVisible(!isChartVisible);
    setIsLogsVisible(false);
  };
  const showLogs = () => {
    setIsLogsVisible(!isLogsVisible);
    setIsChartVisible(false);
  };

  const onRefresh = useCallback(() => {
    fetchData();
  }, [fetchData]);

  const playAgain = () => {
    if (data) {
      const playerIds = data.player_results.map((player) => player.id);
      navigate('/new', { state: { playerIds } });
    }
  };

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate('/')}
    >
      {data ? (
        <div className="flex flex-col grow justify-center w-screen sm:max-w-[600px] h-full px-3 gap-3">
          <GameResultTable data={data} />
          <div className="flex flex-wrap gap-3 [&>*]:flex-1 [&>*]:basis-[calc(50%-0.375rem)] sm:[&>*]:basis-0">
            <Button
              type="button"
              onPress={showChart}
              variant="tertiary"
              fullWidth
            >
              {isChartVisible ? 'Hide chart' : 'Show chart'}
            </Button>
            <Button
              type="button"
              onPress={showLogs}
              variant="tertiary"
              fullWidth
            >
              {isLogsVisible ? 'Hide logs' : 'Show logs'}
            </Button>
            {!isEventGameRoute && isAuthenticated && canCreateGame && (
              <Button
                type="button"
                onPress={playAgain}
                variant="primary"
                fullWidth
              >
                Play Again
              </Button>
            )}
            <LinkButton to="/" variant="secondary" fullWidth>
              Home
            </LinkButton>
          </div>
          {isLogsVisible && (
            <LogsMenu
              gameId={data.id}
              eventId={resolvedEventId ?? undefined}
              canUndoLastRound={canUndoLastRound}
              onUndoSuccess={() => {
                navigate('/now');
              }}
            />
          )}
          {isChartVisible && (
            <GameLogs
              isChartVisible={true}
              isLogVisible={false}
              gameId={data.id}
              eventId={eventId ?? storedEventId ?? undefined}
            />
          )}
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
