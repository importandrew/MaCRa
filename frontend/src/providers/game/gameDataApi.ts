import { useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { GetRequest } from '../../components/Request';
import { GameData, GamePlayers, PlayerOrder } from '../../models/GameData';
import { GameResult, GameRoundLogs } from '../../models/GameResult';

export const useFetchGameData = (isToastEnabled = false) => {
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [data, setData] = useState<GameData | null>(null);
  const [playersList, setPlayersList] = useState<PlayerOrder[]>([]);

  const fetchGameData = async () => {
    setIsLoading(true);
    try {
      // Request game data
      const gameDataRequest = GetRequest('/api/now');
      const response = await fetch(gameDataRequest);
      const result: GameData | { error: string } = await response?.json();

      if ('error' in result) {
        isToastEnabled && toast.error(result.error);
        setIsError(true);
        return;
      }

      setData(result);

      // Request players list only once on game start or page refresh
      if (playersList.length) return;

      const gamePlayersRequest = GetRequest(`/api/game/${result.id}/players`);
      const playersResponse = await fetch(gamePlayersRequest);
      const playersResult: GamePlayers | { error: string } =
        await playersResponse?.json();

      if ('error' in playersResult) {
        isToastEnabled && toast.error(playersResult.error);
        setIsError(true);
        return;
      }

      setPlayersList(playersResult.players);
    } catch (error) {
      isToastEnabled &&
        toast.error('Failed to fetch game data. Please try again.');
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
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

    try {
      const response = await fetch(GetRequest('/api/last'));

      if (!response.ok) {
        isToastEnabled && toast.error(
          `${response.status} ${response.statusText} while processing URL /api/last`
        );
        setIsError(true);
        return;
      }

      const result = await response.json();
      setData(result || []);
    } catch (error) {
      isToastEnabled && toast.error(`Failed to fetch data`);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
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

    try {
      const response = await fetch(GetRequest(`/api/game/${gameId}/logs`));

      if (!response.ok) {
        toast.error(
          `${response.status} ${response.statusText} while processing URL /api/game/${gameId}/logs`
        );
        setIsError(true);
        return;
      }

      const result = await response.json();
      setData(result || []);
    } catch (error) {
      toast.error(`Failed to fetch data`);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isError, isLoading, fetchData };
};