import { createContext, useContext } from 'react';
import { Event } from '../../models/Event';
import { EventMe } from '../../models/EventMe';

export type EventLayoutContextValue = {
  eventId: number;
  event: Event | null;
  me: EventMe | null;
  isLoading: boolean;
  isError: boolean;
  isForbidden: boolean;
  refetch: () => void;
  hasCapability: (
    cap: import('../../models/EventMe').EventCapability
  ) => boolean;
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
