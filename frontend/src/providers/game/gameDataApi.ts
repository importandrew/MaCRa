import { useCallback, useEffect, useRef, useState } from 'react';
import {
  DeleteRequest,
  GetRequest,
  globalBatchFetch,
  globalFetch,
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

export async function undoLastRound(gameId: number): Promise<boolean> {
  const request = DeleteRequest(`/api/game/${gameId}/undo-round`);
  const result = await globalFetch<null>(request);
  return result !== null;
}

export const useFetchPlayerGamesHistory = (playerId: number) => {
  const [data, setData] = useState<GamesHistoryData | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);

    const request = GetRequest(`/api/player/${playerId}/game-history`);
    const result = await globalFetch<GamesHistoryData>(request);

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

export const useFetchPlayerGameResults = (gameId: number) => {
  const [data, setData] = useState<GameResult | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);

    const request = GetRequest(`/api/game/${gameId}/result`);
    const result = await globalFetch<GameResult>(request);

    if (!result) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(result);
    setIsLoading(false);
  }, [gameId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, refetch: fetchData };
};

export const useFetchBatchGameResults = (gameIds: number[]) => {
  const [data, setData] = useState<Record<number, GameResult> | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (gameIds.length === 0) {
      setData({});
      return;
    }
    setIsLoading(true);
    setIsError(false);

    const requests = gameIds.map((id) => {
      return GetRequest(`/api/game/${id}/result`);
    });

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
  }, [gameIds]);

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
    setIsLoading(true);

    const gameDataRequest = GetRequest('/api/now');
    const data = await globalFetch<GameData>(gameDataRequest, {
      isToastEnabled: true,
    });

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);

    // Request players list only once on game start or page refresh
    if (playersList.length) {
      setIsLoading(false);
      return;
    }

    const gamePlayersRequest = GetRequest(`/api/game/${data.id}/players`);
    const playersResult = await globalFetch<GamePlayers>(gamePlayersRequest);

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

export const useFetchGameResultData = (gameId?: number) => {
  const [data, setData] = useState<GameResult | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);

    const lastGameRequest = GetRequest('/api/last');
    const gameResultsByIdRequest = GetRequest(`/api/game/${gameId}/result`);

    const data = await globalFetch<GameResult>(
      gameId ? gameResultsByIdRequest : lastGameRequest,
      {
        isToastEnabled: true,
      }
    );

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, fetchData };
};

export const useFetchGameLogsData = (gameId: number | null) => {
  const [data, setData] = useState<GameRoundLogs | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    if (gameId == null) {
      setData(null);
      return;
    }
    setIsLoading(true);
    setIsError(false);

    const gameLogsRequest = GetRequest(`/api/game/${gameId}/logs`);
    const data = await globalFetch<GameRoundLogs>(gameLogsRequest);

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);
    setIsLoading(false);
  }, [gameId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, fetchData };
};

export const useFetchPlayerStatus = () => {
  const [data, setData] = useState<PlayerStatus | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);

    const playerStatusRequest = GetRequest('/api/main');
    const data = await globalFetch<PlayerStatus>(playerStatusRequest);

    if (!data) {
      setData(null);
      setIsLoading(false);
      return;
    }

    setData(data);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isLoading, fetchData };
};

export const useFetchLiveGameData = (gameId: number) => {
  const [data, setData] = useState<GameView | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [playersList, setPlayersList] = useState<PlayerOrder[]>([]);
  const playersFetchedRef = useRef(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);

    const liveGameRequest = GetRequest(`/api/game/${gameId}/live`);
    const data = await globalFetch<GameView>(liveGameRequest, {
      isToastEnabled: false,
    });

    if (!data) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    setData(data);

    // Request players list only once on game start or page refresh
    if (playersFetchedRef.current) {
      setIsLoading(false);
      return;
    }

    const gamePlayersRequest = GetRequest(`/api/game/${data.id}/players`);
    const playersResult = await globalFetch<GamePlayers>(gamePlayersRequest);

    if (!playersResult) {
      setIsError(true);
    }

    setPlayersList(playersResult?.players || []);
    playersFetchedRef.current = true;
    setIsLoading(false);
  }, [gameId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, playersList, isLoading, isError, refetch: fetchData };
};
