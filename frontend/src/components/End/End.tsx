import { Button, Link, Spinner } from '@heroui/react';
import { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetchGameResultData } from '../../providers/game/gameDataApi';
import GameLogs from './GameLogs';
import GameResultsTable from './GameResultsTable';
import LogsMenu from '../Now/Menu/LogsMenu';

export default function End() {
  const navigate = useNavigate();
  const { data, isError, isLoading, fetchData } = useFetchGameResultData(true);
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

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex w-screen h-full justify-center items-center">
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex flex-col h-full w-screen justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  return (
    <div className="flex flex-col grow justify-center w-screen sm:max-w-[600px] h-full px-3 gap-3">
      <GameResultsTable data={data} />
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
        <GameLogs isChartVisible={true} isLogVisible={false} gameId={data.id} />
      )}
    </div>
  );
}
