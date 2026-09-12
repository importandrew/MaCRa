import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetchEventGames } from '../../providers/event/eventApi';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import {
  EventGamesList,
  useEventGamesNavigate,
} from './EventGamesList';

export default function LiveGames() {
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
  const onSelect = useEventGamesNavigate(eventId);

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const forbidden =
    isForbidden || listForbidden || (!canView && !isLoading);

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={forbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate(`/event/${eventId}`)}
    >
      <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
        <h1
          className="text-2xl font-semibold truncate min-w-0"
          title={`Live games — ${event?.name ?? 'Event'}`}
        >
          Live games — {event?.name ?? 'Event'}
        </h1>
        <EventGamesList
          games={games}
          emptyContent="No live games"
          ariaLabel="Live event games"
          mode="live"
          onSelect={onSelect}
          hasMore={hasMore}
          isLoadingMore={isLoadingMore}
          onLoadMore={loadMore}
        />
      </div>
    </LoadingTemplate>
  );
}
