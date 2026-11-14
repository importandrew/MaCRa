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

// EventsTab is the tab of /events: the user's events or the public ones.
export type EventsTab = 'mine' | 'public';

// EVENTS_TAB is the localStorage key of the /events tab the user picked,
// a convenience on this device only.
const EVENTS_TAB = 'eventsTab';

// readEventsTab returns the tab picked last time, or null (none, or storage
// unavailable).
export function readEventsTab(): EventsTab | null {
  try {
    const stored = localStorage.getItem(EVENTS_TAB);
    return stored === 'mine' || stored === 'public' ? stored : null;
  } catch {
    return null;
  }
}

export function storeEventsTab(tab: EventsTab): void {
  try {
    localStorage.setItem(EVENTS_TAB, tab);
  } catch {
    // Remembering the tab is only a convenience.
  }
}
