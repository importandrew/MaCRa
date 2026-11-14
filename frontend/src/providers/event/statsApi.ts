import { GetRequest, eventPath, publicFetch } from '../../components/Request';
import { useFetch } from '../useFetch';
import type { PlayerStats } from '../../models/PlayerStats';

// A player's statistics in an event: one's own always, others' as the
// event's rating. Public: a visitor sees them in a public event whose rating
// is public.
export const useFetchPlayerStats = (
  eventId: number | null,
  playerId: number | null
) => {
  const key = eventId != null && playerId ? `${eventId}:${playerId}` : null;
  const f = useFetch<PlayerStats>(key, () =>
    publicFetch<PlayerStats>(
      GetRequest(eventPath(eventId!, `/player/${playerId}/stats`)),
      { quietStatuses: [404] }
    )
  );
  return {
    stats: f.data,
    isLoading: f.isLoading,
    isError: f.isError,
    isForbidden: f.isForbidden,
    isNotFound: f.isNotFound,
    refetch: f.refetch,
  };
};
