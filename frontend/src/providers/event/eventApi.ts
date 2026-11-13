import { useCallback, useEffect, useState } from 'react';
import {
  GetRequest,
  authFetch,
  PatchRequest,
  PostRequest,
  DeleteRequest,
  publicFetch,
  eventPath,
  fetchWithStatus,
} from '../../components/Request';
import { useLoadingState } from '../useLoadingState';
import { useFetch } from '../useFetch';
import { Event, EventVisibility, ContentReach } from '../../models/Event';
import { EventMe } from '../../models/EventMe';
import { EventMember, EventRole } from '../../models/EventRole';
import { EventRatingsData } from '../../models/EventRating';
import {
  EventGameListStatus,
  EventGamesListData,
} from '../../models/EventGames';

export type CreateEventPayload = {
  type: 'regular' | 'tournament';
  name: string;
  description?: string;
  visibility?: EventVisibility;
  ratings_reach?: ContentReach;
  games_reach?: ContentReach;
};

export type UpdateEventPayload = {
  name?: string;
  description?: string;
  visibility?: EventVisibility;
  ratings_reach?: ContentReach;
  games_reach?: ContentReach;
};

export const useFetchEvent = (id: number | null) => {
  // A 404 is quiet: the page says "Event not found" (isNotFound).
  const f = useFetch<Event>(id, () =>
    publicFetch<Event>(GetRequest(`/api/event/${id}`), {
      quietStatuses: [404],
    })
  );
  return {
    data: f.data,
    isError: f.isError,
    isForbidden: f.isForbidden,
    isNotFound: f.isNotFound,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

export const useFetchEventMe = (eventId: number | null) => {
  // A 404 is quiet: the page says "Event not found" (isNotFound).
  const f = useFetch<EventMe>(eventId, () =>
    publicFetch<EventMe>(GetRequest(`/api/event/${eventId}/me`), {
      quietStatuses: [404],
    })
  );
  return {
    data: f.data,
    isError: f.isError,
    isForbidden: f.isForbidden,
    isNotFound: f.isNotFound,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

export const useFetchEventMembers = (eventId: number | null) => {
  const f = useFetch<EventMember[]>(eventId, () =>
    fetchWithStatus<EventMember[]>(GetRequest(`/api/event/${eventId}/members`))
  );
  return {
    data: f.data,
    isError: f.isFailed,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

export async function createEvent(
  payload: CreateEventPayload
): Promise<number | null> {
  const request = PostRequest('/api/event', payload);
  const result = await authFetch<{ id: number }>(request);
  return result?.id ?? null;
}

export async function updateEvent(
  id: number,
  payload: UpdateEventPayload
): Promise<boolean> {
  const request = PatchRequest(`/api/event/${id}`, payload);
  const result = await authFetch<null>(request);
  return result !== null;
}

export async function addEventMember(
  eventId: number,
  userId: number,
  roles: EventRole[]
): Promise<boolean> {
  const request = PostRequest(`/api/event/${eventId}/members`, {
    user_id: userId,
    roles,
  });
  const result = await authFetch<null>(request);
  return result !== null;
}

export async function removeEventMember(
  eventId: number,
  userId: number
): Promise<boolean> {
  const request = DeleteRequest(`/api/event/${eventId}/members?ids=${userId}`);
  const result = await authFetch<null>(request);
  return result !== null;
}

export async function finishEvent(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, '/finish'))
  );
  return result !== null;
}

export const useFetchVisibleEvents = () => {
  const f = useFetch<Event[]>('events', () =>
    publicFetch<Event[]>(GetRequest('/api/event'))
  );
  return {
    data: f.data,
    isError: f.isFailed,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

const DEFAULT_GAMES_PAGE_SIZE = 50;

export const useFetchEventGames = (
  eventId: number | null,
  status: EventGameListStatus | null
) => {
  const [games, setGames] = useState<EventGamesListData['games']>([]);
  const [hasMore, setHasMore] = useState(false);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    `${eventId}/${status}`,
    eventId != null && status != null
  );
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const fetchPage = useCallback(
    async (offset: number, append: boolean) => {
      if (eventId == null || status == null) {
        setGames([]);
        setHasMore(false);
        return;
      }

      if (append) {
        setIsLoadingMore(true);
      } else {
        startLoading();
      }
      setIsError(false);
      setIsForbidden(false);

      const request = GetRequest(
        eventPath(
          eventId,
          `/games?status=${status}&limit=${DEFAULT_GAMES_PAGE_SIZE}&offset=${offset}`
        )
      );
      const result = await publicFetch<EventGamesListData>(request);

      if (result.status === 403) {
        setIsForbidden(true);
        finishLoading(false);
        setIsLoadingMore(false);
        return;
      }

      if (!result.data) {
        setIsError(true);
        finishLoading(false);
        setIsLoadingMore(false);
        return;
      }

      setGames((prev) =>
        append ? [...prev, ...result.data!.games] : result.data!.games
      );
      setHasMore(result.data.has_more);
      finishLoading(true);
      setIsLoadingMore(false);
    },
    [eventId, status, startLoading, finishLoading]
  );

  const refetch = useCallback(() => {
    return fetchPage(0, false);
  }, [fetchPage]);

  const loadMore = useCallback(() => {
    if (!hasMore || isLoadingMore || isLoading) return;
    return fetchPage(games.length, true);
  }, [fetchPage, games.length, hasMore, isLoading, isLoadingMore]);

  useEffect(() => {
    fetchPage(0, false);
  }, [fetchPage]);

  return {
    games,
    hasMore,
    isError,
    isForbidden,
    isLoading,
    isLoadingMore,
    refetch,
    loadMore,
  };
};

export const useFetchEventRatings = (eventId: number | null) => {
  const f = useFetch<EventRatingsData>(eventId, () =>
    publicFetch<EventRatingsData>(GetRequest(eventPath(eventId!, '/ratings')))
  );
  return {
    data: f.data,
    isError: f.isError || f.isNotFound,
    isForbidden: f.isForbidden,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};
