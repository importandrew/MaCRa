import { Button, Link } from '@heroui/react';
import { useCallback, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useFetchGameResultData } from '../../providers/game/gameDataApi';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import LogsMenu from '../CurrentGame/Menu/LogsMenu';
import GameLogs from './GameLogs';
import GameResultTable from './GameResultTable';

export default function GameResult() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const gameId = id ? parseInt(id, 10) : 0;
  const { data, isError, isLoading, fetchData } =
    useFetchGameResultData(gameId);
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
      onRefresh={onRefresh}
      hasNoData={!data}
    >
      {data ? (
        <div className="flex flex-col grow justify-center w-screen sm:max-w-[600px] h-full px-3 gap-3">
          <GameResultTable data={data} />
          <div className="grid grid-cols-2 sm:flex sm:flex-row justify-center gap-3">
            <Button
              type="button"
              onPress={showChart}
              color="default"
              variant="flat"
            >
              {isChartVisible ? 'Hide chart' : 'Show chart'}
            </Button>
            <Button
              type="button"
              onPress={showLogs}
              color="secondary"
              variant="flat"
            >
              {isLogsVisible ? 'Hide logs' : 'Show logs'}
            </Button>
            <Button
              type="button"
              onPress={playAgain}
              color="success"
              variant="flat"
            >
              Play Again
            </Button>
            <Button href="/" as={Link} color="primary" variant="flat">
              Home
            </Button>
          </div>
          {isLogsVisible && (
            <LogsMenu
              gameId={data.id}
              canUndoLastRound={data.can_undo_last_round}
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
            />
          )}
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
