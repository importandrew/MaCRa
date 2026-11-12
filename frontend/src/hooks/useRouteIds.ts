import { useParams } from 'react-router-dom';

// parseId returns a route parameter as a positive integer, or null when it is
// missing or isn't one.
function parseId(raw: string | undefined): number | null {
  if (!raw || !/^\d+$/.test(raw)) return null;
  const id = Number(raw);
  return id > 0 ? id : null;
}

// useRouteIds returns the event and game of the URL (/event/:eventId/...,
// /event/:eventId/game/:gameId), null where the route has none or a bad one.
// Pages about the selected event (/now, /last, /new, /history) use useEvent()
// instead.
export function useRouteIds(): {
  eventId: number | null;
  gameId: number | null;
} {
  const { eventId, gameId } = useParams<{
    eventId?: string;
    gameId?: string;
  }>();
  return { eventId: parseId(eventId), gameId: parseId(gameId) };
}
