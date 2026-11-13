import { useCallback } from 'react';
import {
  DeleteRequest,
  GetRequest,
  FetchResponse,
  eventGamePath,
  eventPath,
  fetchWithStatus,
  globalBatchFetch,
  publicFetch,
} from '../../components/Request';
import { useFetch, type FetchKey } from '../useFetch';
import { GameData, GamesHistoryData } from '../../models/GameData';
import {
  GameResult,
  GameRoundLogs,
  PlayerStatus,
} from '../../models/GameResult';
import { GameView } from '../../models/GameView';
import type { NotFound } from '../../components/LoadingTemplate/NotFoundState';
import { useEvent } from '../event/eventConstants';

// undoLastRound undoes the last log entry of the game. logSize is the number
// of log entries the referee saw; the server answers 409 if the log has
// changed since, or if the final round can't be undone. The caller reports
// 409 itself.
export async function undoLastRound(
  eventId: number,
  gameId: number,
  logSize: number
): Promise<FetchResponse<null>> {
  const request = DeleteRequest(
    eventGamePath(eventId, gameId, `/undo-round?log_size=${logSize}`)
  );
  return fetchWithStatus<null>(request, { quietStatuses: [409] });
}

export const useFetchPlayerGamesHistory = (playerId: number) => {
  const { selectedEventId: eventId } = useEvent();
  const f = useFetch<GamesHistoryData>(
    eventId ? `${eventId}/${playerId}` : null,
    () =>
      fetchWithStatus<GamesHistoryData>(
        GetRequest(eventPath(eventId!, `/player/${playerId}/game-history`))
      )
  );
  return {
    data: f.data,
    isError: f.isFailed,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

export const useFetchBatchGameResults = (
  gameIds: number[],
  eventId: number | null
) => {
  const f = useFetch<Record<number, GameResult>>(
    eventId && gameIds.length ? `${eventId}/${gameIds.join(',')}` : null,
    async () => {
      const results = await globalBatchFetch<GameResult>(
        gameIds.map((id) => GetRequest(eventGamePath(eventId!, id, '/result')))
      );
      if (!results) return { data: null, status: 0 };

      // A game whose result can't be loaded is left out.
      const byId: Record<number, GameResult> = {};
      results.forEach((result, index) => {
        if (result) byId[gameIds[index]] = result;
      });
      return { data: byId, status: 200 };
    }
  );
  return {
    data: f.data,
    isError: f.isFailed,
    isLoading: f.isLoading,
    refetch: f.refetch,
  };
};

export const useFetchGameData = () => {
  const { selectedEventId: eventId } = useEvent();
  // Once a table is shown, a 404 means its game has ended or been aborted;
  // the caller tells the player (useShowGameGone), so there is no error
  // toast then. refetch returns the /now response.
  const f = useFetch<GameData>(eventId, (shown) =>
    fetchWithStatus<GameData>(GetRequest(eventPath(eventId!, '/now')), {
      quietStatuses: shown ? [404] : [],
    })
  );
  return {
    data: f.data,
    isLoading: f.isLoading,
    isError: f.isFailed,
    refetch: f.refetch,
  };
};

type GameResultFetchOptions = {
  gameId?: number;
  eventId?: number;
};

export const useFetchGameResultData = (options?: GameResultFetchOptions) => {
  const gameId = options?.gameId;
  const eventId = options?.eventId;
  const { selectedEventId } = useEvent();
  // The result of a given game (a public page), or the user's last game in
  // the selected event (/last).
  const isGivenGame = Boolean(gameId && eventId);
  let key: FetchKey = null;
  if (isGivenGame) key = `${eventId}/${gameId}`;
  else if (selectedEventId) key = `last/${selectedEventId}`;

  const f = useFetch<GameResult>(key, () =>
    isGivenGame
      ? publicFetch<GameResult>(
          GetRequest(eventGamePath(eventId!, gameId!, '/result'))
        )
      : fetchWithStatus<GameResult>(
          GetRequest(eventPath(selectedEventId!, '/last'))
        )
  );
  return {
    data: f.data,
    isError: f.isError || f.isNotFound,
    isForbidden: f.isForbidden,
    isLoading: f.isLoading,
    fetchData: f.refetch,
  };
};

// useFetchGameLogsData loads the round log of a game that has no table on
// screen (the result page). A table screen has the log in its table response.
export const useFetchGameLogsData = (
  gameId: number | null,
  eventId: number | null
) => {
  const f = useFetch<GameRoundLogs>(
    gameId != null && eventId ? `${eventId}/${gameId}` : null,
    () =>
      publicFetch<GameRoundLogs>(
        GetRequest(eventGamePath(eventId!, gameId!, '/logs'))
      )
  );
  // fetchData loads the log again and returns it (null on an error).
  const { refetch } = f;
  const fetchData = useCallback(
    async (): Promise<GameRoundLogs | null> => (await refetch()).data,
    [refetch]
  );
  return {
    data: f.data,
    isError: f.isError || f.isNotFound,
    isForbidden: f.isForbidden,
    isLoading: f.isLoading,
    fetchData,
  };
};

export const useFetchPlayerStatus = () => {
  const { selectedEventId } = useEvent();
  const f = useFetch<PlayerStatus>(selectedEventId, () =>
    fetchWithStatus<PlayerStatus>(
      GetRequest(eventPath(selectedEventId!, '/main'))
    )
  );
  return {
    data: f.data,
    isLoading: f.isLoading,
    isError: f.isFailed,
    fetchData: f.refetch,
  };
};

export type GameStatusResponse = {
  id: number;
  status: 'created' | 'live' | 'ended' | 'aborted';
};

export const useFetchGameStatus = (
  eventId: number | null,
  gameId: number | null
) => {
  const f = useFetch<GameStatusResponse>(
    eventId && gameId ? `${eventId}/${gameId}` : null,
    () =>
      publicFetch<GameStatusResponse>(
        GetRequest(eventGamePath(eventId!, gameId!, '/status')),
        { quietStatuses: [404] }
      )
  );

  // notFound: the URL's event or game doesn't exist (a bad id, or a 404 from
  // the server, which says which); no toast, the page says so.
  let notFound: NotFound | null = null;
  if (!eventId) notFound = 'event';
  else if (!gameId) notFound = 'game';
  else if (f.isNotFound) {
    notFound = f.error?.error === 'event_not_found' ? 'event' : 'game';
  }

  return {
    status: f.data?.status ?? null,
    isLoading: f.isLoading,
    isError: f.isError,
    isForbidden: f.isForbidden,
    notFound,
    refetch: f.refetch,
  };
};

export const useFetchLiveGameData = (
  eventId: number | null,
  gameId: number | null
) => {
  const f = useFetch<GameView>(
    eventId && gameId ? `${eventId}/${gameId}` : null,
    () =>
      publicFetch<GameView>(
        GetRequest(eventGamePath(eventId!, gameId!, '/live'))
      )
  );
  // refetch loads the table again and returns it (null on an error).
  const { refetch: load } = f;
  const refetch = useCallback(
    async (): Promise<GameView | null> => (await load()).data,
    [load]
  );
  return {
    data: f.data,
    isLoading: f.isLoading,
    isError: f.isError || f.isNotFound,
    isForbidden: f.isForbidden,
    refetch,
  };
};
