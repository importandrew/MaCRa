import { useCallback, useState } from 'react';
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
  const key = eventId != null && status != null ? `${eventId}/${status}` : null;
  const pagePath = (offset: number) =>
    eventPath(
      eventId!,
      `/games?status=${status}&limit=${DEFAULT_GAMES_PAGE_SIZE}&offset=${offset}`
    );

  // The first page comes from useFetch; "Load more" appends the next pages,
  // kept with the key they belong to.
  const first = useFetch<EventGamesListData>(key, () =>
    publicFetch<EventGamesListData>(GetRequest(pagePath(0)))
  );
  const [more, setMore] = useState<{
    key: string | null;
    games: EventGamesListData['games'];
    hasMore: boolean;
  } | null>(null);
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const moreForKey = more?.key === key ? more : null;
  const games = [...(first.data?.games ?? []), ...(moreForKey?.games ?? [])];
  const hasMore = moreForKey
    ? moreForKey.hasMore
    : (first.data?.has_more ?? false);

  const { refetch: refetchFirst } = first;
  const refetch = useCallback(() => {
    setMore(null);
    return refetchFirst();
  }, [refetchFirst]);

  const loadMore = async () => {
    if (!hasMore || isLoadingMore || first.isLoading || key == null) return;
    setIsLoadingMore(true);
    const result = await publicFetch<EventGamesListData>(
      GetRequest(pagePath(games.length))
    );
    setIsLoadingMore(false);
    if (!result.data) return;
    setMore({
      key,
      games: [...(moreForKey?.games ?? []), ...result.data.games],
      hasMore: result.data.has_more,
    });
  };

  return {
    games,
    hasMore,
    isError: first.isError || first.isNotFound,
    isForbidden: first.isForbidden,
    isLoading: first.isLoading,
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
