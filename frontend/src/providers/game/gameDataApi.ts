import { useCallback, useEffect, useRef, useState } from 'react';
import {
  DeleteRequest,
  GetRequest,
  authFetch,
  eventGamePath,
  eventPath,
  globalBatchFetch,
  publicFetch,
} from '../../components/Request';
import {
  GameData,
  GamePlayers,
  GamesHistoryData,
  PlayerOrder,
} from '../../models/GameData';
import {
  GameResult,
  GameRoundLogs,
  PlayerStatus,
} from '../../models/GameResult';
import { GameView } from '../../models/GameView';
import { getStoredEventId, useEvent } from '../event/eventConstants';

export async function undoLastRound(
  gameId: number,
  eventId?: number
): Promise<boolean> {
  const resolvedEventId = eventId ?? getStoredEventId();
  if (!resolvedEventId) return false;

  const request = DeleteRequest(
    eventGamePath(resolvedEventId, gameId, '/undo-round')
  );
  const result = await authFetch<null>(request);
  return result !== null;
}

export const useFetchPlayerGamesHistory = (playerId: number) => {
  const [data, setData] = useState<GamesHistoryData | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    const eventId = getStoredEventId();
    if (!eventId) {
      setIsError(true);
      setData(null);
      return;
    }

    setIsLoading(true);

    const request = GetRequest(
      eventPath(eventId, `/player/${playerId}/game-history`)
    );
    const result = await authFetch<GamesHistoryData>(request);

    if (!result) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result);
    setIsLoading(false);
  }, [playerId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
};

export const useFetchBatchGameResults = (
  gameIds: number[],
  eventId: number
) => {
  const [data, setData] = useState<Record<number, GameResult> | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (!gameIds.length) {
      setData({});
      return;
    }

    setIsLoading(true);
    setIsError(false);

    const requests = gameIds.map((id) =>
      GetRequest(eventGamePath(eventId, id, '/result'))
    );

    const results = await globalBatchFetch<GameResult>(requests);

    if (!results) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    const data: Record<number, GameResult> = results.reduce(
      (acc: Record<number, GameResult>, result: GameResult, index: number) => {
        const id = gameIds[index];
        if (result) {
          acc[id] = result;
        }

        return acc;
      },
      {} as Record<number, GameResult>
    );

    setData(data);
    setIsLoading(false);
  }, [gameIds, eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
};

export const useFetchGameData = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [data, setData] = useState<GameData | null>(null);
  const [playersList, setPlayersList] = useState<PlayerOrder[]>([]);

  const fetchGameData = async () => {
    const eventId = getStoredEventId();
    if (!eventId) {
      setIsError(true);
      return;
    }

    setIsLoading(true);

    const gameDataRequest = GetRequest(eventPath(eventId, '/now'));
    const data = await authFetch<GameData>(gameDataRequest);

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);

    if (playersList.length) {
      setIsLoading(false);
      return;
    }

    const gamePlayersRequest = GetRequest(
      eventGamePath(eventId, data.id, '/players')
    );
    const playersResult = await authFetch<GamePlayers>(gamePlayersRequest);

    if (!playersResult) {
      setIsError(true);
    }

    setPlayersList(playersResult?.players || []);
    setIsLoading(false);
  };

  useEffect(() => {
    fetchGameData();
  }, []);

  return { data, playersList, isLoading, isError, refetch: fetchGameData };
};

type GameResultFetchOptions = {
  gameId?: number;
  eventId?: number;
};

export const useFetchGameResultData = (options?: GameResultFetchOptions) => {
  const gameId = options?.gameId;
  const eventId = options?.eventId;
  const [data, setData] = useState<GameResult | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    if (gameId && eventId) {
      const request = GetRequest(eventGamePath(eventId, gameId, '/result'));
      const result = await publicFetch<GameResult>(request);

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
      return;
    }

    const storedEventId = getStoredEventId();
    if (!storedEventId) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    const lastGameRequest = GetRequest(eventPath(storedEventId, '/last'));
    const data = await authFetch<GameResult>(lastGameRequest);

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);
    setIsLoading(false);
  }, [gameId, eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, fetchData };
};

export const useFetchGameLogsData = (
  gameId: number | null,
  eventId?: number
) => {
  const [data, setData] = useState<GameRoundLogs | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (gameId == null) {
      setData(null);
      return;
    }

    const resolvedEventId = eventId ?? getStoredEventId();
    if (!resolvedEventId) {
      setIsError(true);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(eventGamePath(resolvedEventId, gameId, '/logs'));
    const result = eventId
      ? await publicFetch<GameRoundLogs>(request)
      : { data: await authFetch<GameRoundLogs>(request), status: 200 };

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
  }, [gameId, eventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, fetchData };
};

export const useFetchPlayerStatus = () => {
  const { selectedEventId } = useEvent();
  const [data, setData] = useState<PlayerStatus | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const fetchData = useCallback(async () => {
    if (!selectedEventId) {
      setIsError(false);
      setData(null);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);

    const playerStatusRequest = GetRequest(
      `/api/event/${selectedEventId}/main`
    );
    const data = await authFetch<PlayerStatus>(playerStatusRequest);

    if (!data) {
      setData(null);
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);
    setIsLoading(false);
  }, [selectedEventId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isLoading, isError, fetchData };
};

export type GameStatusResponse = {
  id: number;
  status: 'created' | 'live' | 'ended' | 'aborted';
};

export const useFetchGameStatus = (eventId: number, gameId: number) => {
  const [status, setStatus] = useState<GameStatusResponse['status'] | null>(
    null
  );
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (!eventId || !gameId) {
      setIsError(true);
      setStatus(null);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(eventGamePath(eventId, gameId, '/status'));
    const result = await publicFetch<GameStatusResponse>(request);

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

    setStatus(result.data.status);
    setIsLoading(false);
  }, [eventId, gameId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    status,
    isLoading,
    isError,
    isForbidden,
    refetch: fetchData,
  };
};

export const useFetchLiveGameData = (eventId: number, gameId: number) => {
  const [data, setData] = useState<GameView | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [playersList, setPlayersList] = useState<PlayerOrder[]>([]);
  const playersFetchedRef = useRef(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);
    setIsForbidden(false);

    const liveGameRequest = GetRequest(eventGamePath(eventId, gameId, '/live'));
    const liveResult = await publicFetch<GameView>(liveGameRequest);

    if (liveResult.status === 403) {
      setIsForbidden(true);
      setIsLoading(false);
      return;
    }

    if (!liveResult.data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(liveResult.data);

    if (playersFetchedRef.current) {
      setIsLoading(false);
      return;
    }

    const gamePlayersRequest = GetRequest(
      eventGamePath(eventId, gameId, '/players')
    );
    const playersResult = await publicFetch<GamePlayers>(gamePlayersRequest);

    if (playersResult.status === 403) {
      setIsForbidden(true);
      setIsLoading(false);
      return;
    }

    if (!playersResult.data) {
      setIsError(true);
    }

    setPlayersList(playersResult.data?.players || []);
    playersFetchedRef.current = true;
    setIsLoading(false);
  }, [eventId, gameId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    playersList,
    isLoading,
    isError,
    isForbidden,
    refetch: fetchData,
  };
};
