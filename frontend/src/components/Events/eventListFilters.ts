import type { Event } from '../../models/Event';

// matchesQuery: the event's name contains the query, ignoring case and the
// spaces around it; an empty query matches every event.
export function matchesQuery(event: Event, query: string): boolean {
  const q = query.trim().toLocaleLowerCase();
  return !q || event.name.toLocaleLowerCase().includes(q);
}

// splitFinished splits events into active and finished ones, keeping their
// order (the server's: active by their last game, finished by when they
// finished).
export function splitFinished(events: Event[]): {
  active: Event[];
  finished: Event[];
} {
  return {
    active: events.filter((event) => !event.finished_at),
    finished: events.filter((event) => event.finished_at),
  };
}
