import { useCallback, useEffect, useState } from 'react';
import { GetRequest, globalFetch } from '../../components/Request';
import { GameData, GamePlayers, PlayerOrder } from '../../models/GameData';
import { GameResult, GameRoundLogs } from '../../models/GameResult';

export const useFetchGameData = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [data, setData] = useState<GameData | null>(null);
  const [playersList, setPlayersList] = useState<PlayerOrder[]>([]);

  const fetchGameData = async () => {
    setIsLoading(true);

    const gameDataRequest = GetRequest('/api/now');
    const data = await globalFetch<GameData>(gameDataRequest);

    if (!data || 'error' in data) {
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

    if (!playersResult || 'error' in playersResult) {
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

export const useFetchGameResultData = (isToastEnabled = false) => {
  const [data, setData] = useState<GameResult | null>(null);
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);

    const lastGameRequest = GetRequest('/api/last');
    const data = await globalFetch<GameResult>(lastGameRequest, {
      isToastEnabled,
    });

    if (!data || 'error' in data) {
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
    setIsLoading(true);
    setIsError(false);

    const gameLogsRequest = GetRequest(`/api/game/${gameId}/logs`);
    const data = await globalFetch<GameRoundLogs>(gameLogsRequest);

    if (!data || 'error' in data) {
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
