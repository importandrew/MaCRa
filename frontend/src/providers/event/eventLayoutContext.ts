import { createContext, useContext } from 'react';
import { Event } from '../../models/Event';
import { EventCapability, EventMe } from '../../models/EventMe';

export type EventLayoutContextValue = {
  eventId: number;
  event: Event | null;
  me: EventMe | null;
  isLoading: boolean;
  isError: boolean;
  isForbidden: boolean;
  // refetch reloads the event and the capabilities; it resolves when both
  // have loaded (a form keeps its edits on screen until then).
  refetch: () => Promise<void>;
  hasCapability: (cap: EventCapability) => boolean;
};

export const EventLayoutContext = createContext<EventLayoutContextValue | null>(
  null
);

export function useEventLayout() {
  const ctx = useContext(EventLayoutContext);
  if (!ctx) {
    throw new Error('useEventLayout must be used within EventLayout');
  }
  return ctx;
}

// useEventPage is useEventLayout for a page of the event that needs one of
// caps: canView says whether the user has one. EventLayout shows its pages
// only once the event and the capabilities have loaded and the event is
// accessible, so a page is forbidden exactly when canView is false (or when
// its own request answers 403).
export function useEventPage(...caps: EventCapability[]) {
  const layout = useEventLayout();
  const canView = caps.some((cap) => layout.hasCapability(cap));
  return { ...layout, canView };
}
