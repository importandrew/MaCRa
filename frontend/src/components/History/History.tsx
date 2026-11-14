import { useTranslation } from 'react-i18next';
import { useCallback, useMemo, useState } from 'react';
import { HistoryProps } from '../../models/GameData';
import { useAuth } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import {
  useFetchBatchGameResults,
  useFetchPlayerGamesHistory,
} from '../../providers/game/gameDataApi';
import { getDateTimeFormatted } from '../../utils/date-tools';
import NavList from '../NavList';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import SegmentedPicker from '../SegmentedPicker';
import PlayerStatsView from '../Stats/PlayerStatsView';
import { useFetchPlayerStats } from '../../providers/event/statsApi';

const History = (props: HistoryProps) => {
  const { user } = useAuth();
  const { t } = useTranslation();
  // History is behind RequireAuth, so a user is signed in.
  const { playerId = user?.id ?? 0 } = props;
  const { data, isLoading, isError, refetch } =
    useFetchPlayerGamesHistory(playerId);

  const { selectedEventId: eventId } = useEvent();
  const [tab, setTab] = useState<'games' | 'stats'>('games');
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
    const playerResult = result.player_results.find((pr) => pr.id === playerId);
    return playerResult;
  };

  const items = data
    ? data.game_ids.map((id) => {
        const playerResult = getPlayerResult(id);
        const result = resultsMap?.[id];
        const label = playerResult
          ? t('history.result', { tp: playerResult.tp, gp: playerResult.gp })
          : resultsLoading
            ? t('common.loading')
            : resultsError
              ? t('history.loadFailed')
              : t('history.noResult');

        return {
          key: id,
          label,
          date: getDateTimeFormatted(result?.finished_at || ''),
        };
      })
    : [];

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] justify-start items-center gap-4 md:my-3 px-3 md:px-0">
      <SegmentedPicker
        ariaLabel={t('history.sections')}
        options={[
          { value: 'games', label: t('history.games') },
          { value: 'stats', label: t('history.stats') },
        ]}
        selectedKey={tab}
        onSelectionChange={setTab}
      />
      {tab === 'stats' ? (
        <HistoryStats eventId={eventId} playerId={playerId} />
      ) : (
        <LoadingTemplate
          isLoading={isLoadingState}
          isError={isErrorState}
          onRefresh={onRefresh}
        >
          <NavList
            ariaLabel={t('landing.playedGames')}
            emptyContent={
              <p className="text-center text-3xl">{t('rating.noGames')}</p>
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
        </LoadingTemplate>
      )}
    </div>
  );
};

// HistoryStats is the history's "Statistics" tab: the player's statistics in
// the selected event.
function HistoryStats({
  eventId,
  playerId,
}: {
  eventId: number | null;
  playerId: number;
}) {
  const { stats, isLoading, isError, refetch } = useFetchPlayerStats(
    eventId,
    playerId
  );
  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      onRefresh={refetch}
    >
      {stats && <PlayerStatsView stats={stats} />}
    </LoadingTemplate>
  );
}

export default History;
