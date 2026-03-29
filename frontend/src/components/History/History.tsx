import { Listbox, ListboxItem, cn } from '@heroui/react';
import { useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { HistoryProps } from '../../models/GameData';
import { UserData } from '../../models/UserData';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import {
  useFetchBatchGameResults,
  useFetchPlayerGamesHistory,
} from '../../providers/game/gameDataApi';
import { getDateTimeFormatted } from '../../utils/date-tools';
import { ChevronRightIcon } from '../Icons/ChevronRightIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

export const ListboxWrapper = ({ children }: { children: React.ReactNode }) => (
  <div className="w-full max-w-[600px] md:border-small rounded-small border-default-200 dark:border-default-100">
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

  const gameIds = useMemo(() => data?.game_ids || [], [data?.game_ids]);

  const {
    data: resultsMap,
    isLoading: resultsLoading,
    isError: resultsError,
  } = useFetchBatchGameResults(gameIds);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleAction = (id: number) => {
    navigate(`/game/${id}/result`);
  };

  const isLoadingState = isLoading || resultsLoading;
  const isErrorState = isError || resultsError;

  const getPlayerResult = (gameId: number) => {
    const result = resultsMap?.[gameId];
    if (!result) return null;
    const playerResult = result.player_results.find(
      (pr) => pr.id === currentPlayer.id
    );
    return playerResult;
  };

  const items = data
    ? data.game_ids.map((id) => {
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
        date: getDateTimeFormatted(result?.finished_at || ''),
        isLastInList: id === data.game_ids[data.game_ids.length - 1],
      };
    })
    : [];

  return (
    <LoadingTemplate
      isLoading={isLoadingState}
      isError={isErrorState}
      onRefresh={onRefresh}
    >
      <div className="flex flex-col h-full w-screen md:w-[600px] justify-start items-center md:my-3 px-3 md:px-0">
        <ListboxWrapper>
          <Listbox
            emptyContent={<p className="text-center text-3xl">No games played yet</p>}
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
                showDivider={!item.isLastInList}
                endContent={<ItemCounter />}
              >
                <div className="flex flex-col gap-1">
                  <span>{item.label}</span>
                  <div className="w-auto rounded-small">
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
    </LoadingTemplate>
  );
};

export default History;
