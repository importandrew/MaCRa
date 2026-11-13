import { useCallback, useEffect, useState } from 'react';
import {
  GetRequest,
  eventGamePath,
  publicFetch,
} from '../../components/Request';
import type { GameCapability, GameMe } from '../../models/EventMe';

// useGameCapabilities loads what the user may do with one game: undo, abort,
// finish, enter rounds. A menu that uses it mounts when it opens, so it asks
// again each time and follows the game (e.g. abort ends with the first
// round for seated players).
export function useGameCapabilities(
  eventId: number | null,
  gameId: number | null
) {
  const [data, setData] = useState<GameMe | null>(null);

  const fetchData = useCallback(async () => {
    if (!eventId || !gameId) {
      setData(null);
      return;
    }
    const result = await publicFetch<GameMe>(
      GetRequest(eventGamePath(eventId, gameId, '/me')),
      { isToastEnabled: false }
    );
    setData(result.data);
  }, [eventId, gameId]);

  useEffect(() => {
    void fetchData();
  }, [fetchData]);

  const hasCapability = useCallback(
    (cap: GameCapability) => data?.capabilities.includes(cap) ?? false,
    [data]
  );

  return { hasCapability, refetch: fetchData };
}
