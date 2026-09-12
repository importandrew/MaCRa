import { useCallback, useEffect, useState } from 'react';
import {
  GetRequest,
  authFetch,
  PatchRequest,
  PostRequest,
  DeleteRequest,
  publicFetch,
  eventPath,
} from '../../components/Request';
import { showToast } from '../../components/Toast';
import { Event, EventVisibility, ContentReach } from '../../models/Event';
import { EventMe } from '../../models/EventMe';
import { EventMember, EventRole } from '../../models/EventRole';
import { GamesHistoryData } from '../../models/GameData';
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
  const [data, setData] = useState<Event | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (id == null) {
      setData(null);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(`/api/event/${id}`);
    const result = await publicFetch<Event>(request);

    if (result.status === 403) {
      setIsForbidden(true);
      setIsLoading(false);
      return;
    }

    if (!result.data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result.data);
    setIsLoading(false);
  }, [id]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, refetch: fetchData };
};

export const useFetchEventMe = (eventId: number | null) => {
  const [data, setData] = useState<EventMe | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(() => eventId != null);

  const fetchData = useCallback(async () => {
    if (eventId == null) {
      setData(null);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(`/api/event/${eventId}/me`);
    const result = await publicFetch<EventMe>(request);

    if (result.status === 403) {
      setIsForbidden(true);
      setIsLoading(false);
      return;
    }

    if (!result.data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result.data);
    setIsLoading(false);
  }, [eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, refetch: fetchData };
};

export const useFetchEventMembers = (eventId: number | null) => {
  const [data, setData] = useState<EventMember[] | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (eventId == null) {
      setData(null);
      return;
    }

    setIsLoading(true);
    setIsError(false);

    const request = GetRequest(`/api/event/${eventId}/members`);
    const result = await authFetch<EventMember[]>(request);

    if (!result) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result);
    setIsLoading(false);
  }, [eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
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
  const response = await fetch(request);

  if (!response.ok) {
    const data = await response
      .json()
      .catch(() => ({ error: 'Unknown error' }));
    showToast({
      title: `Error: ${data.error}`,
      color: 'danger',
    });
    return false;
  }

  return true;
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
  const [data, setData] = useState<Event[] | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);

    const request = GetRequest('/api/event');
    const result = await publicFetch<Event[]>(request);

    if (result.status === 403 || !result.data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result.data);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
};

export const useFetchEventGameHistory = (eventId: number | null) => {
  const [data, setData] = useState<GamesHistoryData | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (eventId == null) {
      setData(null);
      return;
    }

    setIsLoading(true);
    setIsError(false);

    const request = GetRequest(eventPath(eventId, '/game-history'));
    const result = await authFetch<GamesHistoryData>(request);

    if (!result) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result);
    setIsLoading(false);
  }, [eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
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
  const [isLoading, setIsLoading] = useState(false);
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
        setIsLoading(true);
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
        setIsLoading(false);
        setIsLoadingMore(false);
        return;
      }

      if (!result.data) {
        setIsError(true);
        setIsLoading(false);
        setIsLoadingMore(false);
        return;
      }

      setGames((prev) =>
        append ? [...prev, ...result.data!.games] : result.data!.games
      );
      setHasMore(result.data.has_more);
      setIsLoading(false);
      setIsLoadingMore(false);
    },
    [eventId, status]
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
  const [data, setData] = useState<EventRatingsData | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (eventId == null) {
      setData(null);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(eventPath(eventId, '/ratings'));
    const result = await publicFetch<EventRatingsData>(request);

    if (result.status === 403) {
      setIsForbidden(true);
      setIsLoading(false);
      return;
    }

    if (!result.data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result.data);
    setIsLoading(false);
  }, [eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, refetch: fetchData };
};
