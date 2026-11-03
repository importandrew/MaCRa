import { useCallback, useMemo } from 'react';
import { HistoryProps } from '../../models/GameData';
import { UserData } from '../../models/UserData';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { getStoredEventId } from '../../providers/event/eventConstants';
import {
  useFetchBatchGameResults,
  useFetchPlayerGamesHistory,
} from '../../providers/game/gameDataApi';
import { getDateTimeFormatted } from '../../utils/date-tools';
import NavList from '../NavList';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

const History = (props: HistoryProps) => {
  const currentPlayer = JSON.parse(
    localStorage.getItem(AUTH_USER_DATA) || '{}'
  ) as UserData;
  const { playerId = currentPlayer.id } = props;
  const { data, isLoading, isError, refetch } =
    useFetchPlayerGamesHistory(playerId);

  const eventId = getStoredEventId() ?? 0;
  const gameIds = useMemo(() => data?.game_ids || [], [data?.game_ids]);

  const {
    data: resultsMap,
    isLoading: resultsLoading,
    isError: resultsError,
  } = useFetchBatchGameResults(gameIds, eventId);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

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
        <NavList
          ariaLabel="Played games history"
          emptyContent={
            <p className="text-center text-3xl">No games played yet</p>
          }
          items={items.map((item) => ({
            id: String(item.key),
            href: `/event/${eventId}/game/${item.key}`,
            textValue: item.label,
            content: (
              <div className="flex flex-col gap-1">
                <span>{item.label}</span>
                <span className="text-xs text-muted">{item.date}</span>
              </div>
            ),
          }))}
        />
      </div>
    </LoadingTemplate>
  );
};

export default History;
