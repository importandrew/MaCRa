import { useTranslation } from 'react-i18next';
import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetchEventGames } from '../../providers/event/eventApi';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import { EventGamesList } from './EventGamesList';

export default function LiveGames() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event, hasCapability, isForbidden } = useEventLayout();
  const canView = hasCapability('watch_live');
  const {
    games,
    hasMore,
    isLoading,
    isLoadingMore,
    isError,
    isForbidden: listForbidden,
    refetch,
    loadMore,
  } = useFetchEventGames(canView ? eventId : null, canView ? 'live' : null);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const forbidden = isForbidden || listForbidden || (!canView && !isLoading);
  const title = t('events.liveGamesTitle', {
    name: event?.name ?? t('events.event'),
  });

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={forbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate(`/event/${eventId}`)}
    >
      <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
        <h1 className="text-2xl font-semibold truncate min-w-0" title={title}>
          {title}
        </h1>
        <EventGamesList
          games={games}
          emptyContent={t('events.noLiveGames')}
          ariaLabel={t('events.liveGames')}
          mode="live"
          gameHref={(gameId) => `/event/${eventId}/game/${gameId}`}
          hasMore={hasMore}
          isLoadingMore={isLoadingMore}
          onLoadMore={loadMore}
        />
      </div>
    </LoadingTemplate>
  );
}
