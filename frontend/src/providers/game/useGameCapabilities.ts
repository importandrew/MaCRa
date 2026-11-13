import { useCallback } from 'react';
import {
  GetRequest,
  eventGamePath,
  publicFetch,
} from '../../components/Request';
import type { GameCapability, GameMe } from '../../models/EventMe';
import { useFetch } from '../useFetch';

// useGameCapabilities loads what the user may do with one game: undo, abort,
// finish, enter rounds. A menu that uses it mounts when it opens, so it asks
// again each time and follows the game (e.g. abort ends with the first
// round for seated players).
export function useGameCapabilities(
  eventId: number | null,
  gameId: number | null
) {
  const { data, refetch } = useFetch<GameMe>(
    eventId && gameId ? `${eventId}/${gameId}` : null,
    () =>
      publicFetch<GameMe>(GetRequest(eventGamePath(eventId!, gameId!, '/me')), {
        isToastEnabled: false,
      })
  );

  const hasCapability = useCallback(
    (cap: GameCapability) => data?.capabilities.includes(cap) ?? false,
    [data]
  );

  return { hasCapability, refetch };
}
