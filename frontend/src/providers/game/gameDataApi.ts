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
  const [hasLoadedPlayers, setHasLoadedPlayers] = useState(false);

  const fetchGameData = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);
    try {
      // Request game data
      const gameDataRequest = GetRequest('/api/now');
      const response = await fetch(gameDataRequest);

      if (!response.ok) {
        isToastEnabled && toast.error(
          `${response.status} ${response.statusText} while processing URL /api/now`
        );
        setIsError(true);
        setIsLoading(false);
        return;
      }

      const result: GameData | { error: string } = await response?.json();

      if ('error' in result) {
        isToastEnabled && toast.error(result.error);
        setIsError(true);
        setIsLoading(false);
        return;
      }

      setData(result);
      setIsError(false);

      // Request players list only once on game start or page refresh
      // This ensures PlayerOrder remains consistent through game progression
      if (hasLoadedPlayers) {
        setIsLoading(false);
        return;
      }

      const gamePlayersRequest = GetRequest(`/api/game/${result.id}/players`);
      const playersResponse = await fetch(gamePlayersRequest);

      if (!playersResponse.ok) {
        isToastEnabled && toast.error(
          `${playersResponse.status} ${playersResponse.statusText} while processing URL /api/game/${result.id}/players`
        );
        setIsError(true);
        setIsLoading(false);
        return;
      }

      const playersResult: GamePlayers | { error: string } =
        await playersResponse?.json();

      if ('error' in playersResult) {
        isToastEnabled && toast.error(playersResult.error);
        setIsError(true);
        setIsLoading(false);
        return;
      }

      setPlayersList(playersResult.players);
      setHasLoadedPlayers(true);
      setIsError(false);
      setIsLoading(false);
    } catch (error) {
      isToastEnabled &&
        toast.error('Failed to fetch game data. Please try again.');
      setIsError(true);
      setIsLoading(false);
    }
  }, [isToastEnabled, hasLoadedPlayers]);

  useEffect(() => {
    fetchGameData();
  }, [fetchGameData]);

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