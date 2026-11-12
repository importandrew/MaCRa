import { useCallback, useEffect, useState } from 'react';
import {
  DeleteRequest,
  GetRequest,
  FetchResponse,
  authFetch,
  eventGamePath,
  eventPath,
  fetchWithStatus,
  globalBatchFetch,
  publicFetch,
} from '../../components/Request';
import { useLoadingState } from '../useLoadingState';
import { GameData, GamesHistoryData } from '../../models/GameData';
import {
  GameResult,
  GameRoundLogs,
  PlayerStatus,
} from '../../models/GameResult';
import { GameView } from '../../models/GameView';
import { getStoredEventId, useEvent } from '../event/eventConstants';

// undoLastRound undoes the last log entry of the game. logSize is the number
// of log entries the referee saw; the server answers 409 if the log has
// changed since, or if the final round can't be undone. The caller reports
// 409 itself.
export async function undoLastRound(
  gameId: number,
  logSize: number,
  eventId?: number
): Promise<FetchResponse<null>> {
  const resolvedEventId = eventId ?? getStoredEventId();
  if (!resolvedEventId) return { data: null, status: 0 };

  const request = DeleteRequest(
    eventGamePath(resolvedEventId, gameId, `/undo-round?log_size=${logSize}`)
  );
  return fetchWithStatus<null>(request, { quietStatuses: [409] });
}

export const useFetchPlayerGamesHistory = (playerId: number) => {
  const [data, setData] = useState<GamesHistoryData | null>(null);
  const [isError, setIsError] = useState(false);
  const { isLoading, startLoading, finishLoading } = useLoadingState(playerId);

  const fetchData = useCallback(async () => {
    const eventId = getStoredEventId();
    if (!eventId) {
      setIsError(true);
      setData(null);
      finishLoading(false);
      return;
    }

    startLoading();
    setIsError(false);

    const request = GetRequest(
      eventPath(eventId, `/player/${playerId}/game-history`)
    );
    const result = await authFetch<GamesHistoryData>(request);

    if (!result) {
      setIsError(true);
      finishLoading(false);
      return;
    }

    setData(result);
    finishLoading(true);
  }, [playerId, startLoading, finishLoading]);

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
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    `${eventId}/${gameIds.join(',')}`,
    gameIds.length > 0
  );

  const fetchData = useCallback(async () => {
    if (!gameIds.length) {
      setData({});
      return;
    }

    startLoading();
    setIsError(false);

    const requests = gameIds.map((id) =>
      GetRequest(eventGamePath(eventId, id, '/result'))
    );

    const results = await globalBatchFetch<GameResult>(requests);

    if (!results) {
      setIsError(true);
      finishLoading(false);
      return;
    }

    const data: Record<number, GameResult> = results.reduce(
      (
        acc: Record<number, GameResult>,
        result: GameResult | null,
        index: number
      ) => {
        const id = gameIds[index];
        if (result) {
          acc[id] = result;
        }

        return acc;
      },
      {} as Record<number, GameResult>
    );

    setData(data);
    finishLoading(true);
  }, [gameIds, eventId, startLoading, finishLoading]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
};

export const useFetchGameData = () => {
  const { isLoading, startLoading, finishLoading } = useLoadingState();
  const [isError, setIsError] = useState(false);
  const [data, setData] = useState<GameData | null>(null);

  // fetchGameData loads the table and returns the /now response (data is
  // null on an error). Once a table is shown, a 404 means its game has ended
  // or been aborted; the caller tells the player (useShowGameGone), so there
  // is no error toast then.
  const fetchGameData = async (): Promise<FetchResponse<GameData>> => {
    const shownGame = data;
    const eventId = getStoredEventId();
    if (!eventId) {
      setIsError(true);
      finishLoading(false);
      return { data: null, status: 0 };
    }

    startLoading();
    setIsError(false);

    const gameDataRequest = GetRequest(eventPath(eventId, '/now'));
    const result = await fetchWithStatus<GameData>(gameDataRequest, {
      quietStatuses: shownGame ? [404] : [],
    });
    const { data: fresh } = result;

    if (!fresh) {
      setIsError(true);
      finishLoading(false);
      return result;
    }

    setData(fresh);
    finishLoading(true);
    return result;
  };

  useEffect(() => {
    fetchGameData();
    // Fetch once on mount. fetchGameData isn't memoized and reads data, so
    // listing it would refetch after every load.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { data, isLoading, isError, refetch: fetchGameData };
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
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    `${eventId}/${gameId}`
  );

  const fetchData = useCallback(async () => {
    startLoading();
    setIsError(false);
    setIsForbidden(false);

    if (gameId && eventId) {
      const request = GetRequest(eventGamePath(eventId, gameId, '/result'));
      const result = await publicFetch<GameResult>(request);

      if (result.status === 403) {
        setIsForbidden(true);
        finishLoading(false);
        return;
      }

      if (!result.data) {
        setIsError(true);
        finishLoading(false);
        return;
      }

      setData(result.data);
      finishLoading(true);
      return;
    }

    const storedEventId = getStoredEventId();
    if (!storedEventId) {
      setIsError(true);
      finishLoading(false);
      return;
    }

    const lastGameRequest = GetRequest(eventPath(storedEventId, '/last'));
    const data = await authFetch<GameResult>(lastGameRequest);

    if (!data) {
      setIsError(true);
      finishLoading(false);
      return;
    }

    setData(data);
    finishLoading(true);
  }, [gameId, eventId, startLoading, finishLoading]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, fetchData };
};

// useFetchGameLogsData loads the round log of a game that has no table on
// screen (the result page). A table screen has the log in its table response.
export const useFetchGameLogsData = (
  gameId: number | null,
  eventId?: number
) => {
  const [data, setData] = useState<GameRoundLogs | null>(null);
  const [isError, setIsError] = useState(false);
  const [isForbidden, setIsForbidden] = useState(false);
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    gameId,
    gameId != null
  );

  // fetchData loads the log and returns it (null on an error).
  const fetchData = useCallback(async (): Promise<GameRoundLogs | null> => {
    if (gameId == null) {
      setData(null);
      return null;
    }

    const resolvedEventId = eventId ?? getStoredEventId();
    if (!resolvedEventId) {
      setIsError(true);
      finishLoading(false);
      return null;
    }

    startLoading();
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(eventGamePath(resolvedEventId, gameId, '/logs'));
    const result = eventId
      ? await publicFetch<GameRoundLogs>(request)
      : { data: await authFetch<GameRoundLogs>(request), status: 200 };

    if (result.status === 403) {
      setIsForbidden(true);
      finishLoading(false);
      return null;
    }

    if (!result.data) {
      setIsError(true);
      finishLoading(false);
      return null;
    }

    setData(result.data);
    finishLoading(true);
    return result.data;
  }, [gameId, eventId, startLoading, finishLoading]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isForbidden, isLoading, fetchData };
};

export const useFetchPlayerStatus = () => {
  const { selectedEventId } = useEvent();
  const [data, setData] = useState<PlayerStatus | null>(null);
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    selectedEventId,
    !!selectedEventId
  );
  const [isError, setIsError] = useState(false);

  const fetchData = useCallback(async () => {
    if (!selectedEventId) {
      setIsError(false);
      setData(null);
      finishLoading(false);
      return;
    }

    startLoading();

    const playerStatusRequest = GetRequest(
      `/api/event/${selectedEventId}/main`
    );
    const data = await authFetch<PlayerStatus>(playerStatusRequest);

    if (!data) {
      setData(null);
      setIsError(true);
      finishLoading(false);
      return;
    }

    setData(data);
    finishLoading(true);
  }, [selectedEventId, startLoading, finishLoading]);

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
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    `${eventId}/${gameId}`
  );

  const fetchData = useCallback(async () => {
    if (!eventId || !gameId) {
      setIsError(true);
      setStatus(null);
      finishLoading(false);
      return;
    }

    startLoading();
    setIsError(false);
    setIsForbidden(false);

    const request = GetRequest(eventGamePath(eventId, gameId, '/status'));
    const result = await publicFetch<GameStatusResponse>(request);

    if (result.status === 403) {
      setIsForbidden(true);
      finishLoading(false);
      return;
    }

    if (!result.data) {
      setIsError(true);
      finishLoading(false);
      return;
    }

    setStatus(result.data.status);
    finishLoading(true);
  }, [eventId, gameId, startLoading, finishLoading]);

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
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    `${eventId}/${gameId}`
  );

  // fetchData loads the table and returns it (null on an error).
  const fetchData = useCallback(async (): Promise<GameView | null> => {
    startLoading();
    setIsError(false);
    setIsForbidden(false);

    const liveGameRequest = GetRequest(eventGamePath(eventId, gameId, '/live'));
    const liveResult = await publicFetch<GameView>(liveGameRequest);

    if (liveResult.status === 403) {
      setIsForbidden(true);
      finishLoading(false);
      return null;
    }

    if (!liveResult.data) {
      setIsError(true);
      finishLoading(false);
      return null;
    }

    setData(liveResult.data);
    finishLoading(true);
    return liveResult.data;
  }, [eventId, gameId, startLoading, finishLoading]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    isLoading,
    isError,
    isForbidden,
    refetch: fetchData,
  };
};
