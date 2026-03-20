import { Button, Listbox, ListboxItem, Spinner, cn } from '@heroui/react';
import { useCallback } from 'react';
import { HistoryProps } from '../../models/GameData';
import { UserData } from '../../models/UserData';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import {
  useFetchAllGameResults,
  useFetchPlayerGamesHistory,
} from '../../providers/game/gameDataApi';
import { ChevronRightIcon } from '../Icons/ChevronRightIcon';
import { PlayCircleIcon } from '../Icons/PlayCircleIcon';
import { useNavigate } from 'react-router-dom';
import { UsersIcon } from '../Icons/UsersIcon';

export const ListboxWrapper = ({ children }: { children: React.ReactNode }) => (
  <div className="w-full max-w-[600px] md:border-small p-3 rounded-small border-default-200 dark:border-default-100">
    {children}
  </div>
);

export const IconWrapper = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => (
  <div
    className={cn(
      className,
      'flex items-center rounded-small justify-center w-7 h-7'
    )}
  >
    {children}
  </div>
);

export const ItemCounter = () => (
  <div className="flex items-center gap-1 text-default-400">
    <span className="text-small">{}</span>
    <ChevronRightIcon className="text-xl" />
  </div>
);

const History = (props: HistoryProps) => {
  const navigate = useNavigate();
  const currentPlayer = JSON.parse(
    localStorage.getItem(AUTH_USER_DATA) || '{}'
  ) as UserData;
  const { playerId = currentPlayer.id } = props;
  const { data, isLoading, isError, refetch } =
    useFetchPlayerGamesHistory(playerId);

  const {
    data: resultsMap,
    isLoading: resultsLoading,
    isError: resultsError,
  } = useFetchAllGameResults(data?.game_ids || []);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleAction = (id: number) => {
    navigate(`/game/${id}/history`);
  };

  if (isLoading || resultsLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  if (isError || resultsError) {
    return (
      <div className="flex w-screen h-full justify-center items-center">
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  if (!data || !resultsMap) {
    return (
      <div className="flex flex-col h-full w-screen justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const getPlayerResult = (gameId: number) => {
    const result = resultsMap?.[gameId];
    if (!result) return null;
    const playerResult = result.player_results.find(
      (pr) => pr.id === currentPlayer.id
    );
    return playerResult;
  };

  const items = data.game_ids.map((id) => {
    const playerResult = getPlayerResult(id);
    const result = resultsMap?.[id];
    const label = playerResult
      ? `TP: ${playerResult.tp}, GP: ${playerResult.gp}`
      : resultsLoading
        ? 'Loading...'
        : resultsError
          ? 'Error loading results'
          : 'No result data';

    return {
      key: id,
      label,
      finished: !result?.can_undo_last_round,
      date: formatDate(result?.finished_at || ''),
    };
  });

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] justify-center items-center md:my-3">
      <ListboxWrapper>
        <Listbox
          aria-label="Played games history"
          variant="flat"
          color="primary"
          items={items}
          onAction={(key) => handleAction(Number(key))}
        >
          {(item) => (
            <ListboxItem
              key={item.key}
              textValue={item.label}
              startContent={item.finished ? <UsersIcon /> : <PlayCircleIcon />}
              endContent={<ItemCounter />}
            >
              <div className="flex flex-col gap-1">
                <span>{item.label}</span>
                <div className="px-2 py-1 w-auto rounded-small bg-primary-100">
                  <span className="text-tiny text-default-600">
                    {item.date}
                  </span>
                </div>
              </div>
            </ListboxItem>
          )}
        </Listbox>
      </ListboxWrapper>
    </div>
  );
};

export default History;
