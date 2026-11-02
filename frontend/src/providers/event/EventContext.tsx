import { ReactNode, useCallback, useContext, useEffect, useState } from 'react';
import { Event } from '../../models/Event';
import { GetRequest, authFetch } from '../../components/Request';
import { AuthContext } from '../auth/authConstants';
import {
  clearStoredEventId,
  EventContext,
  getStoredEventId,
  setStoredEventId,
} from './eventConstants';

async function fetchEventById(id: number): Promise<Event | null> {
  const request = GetRequest(`/api/event/${id}`);
  return authFetch<Event>(request, { isToastEnabled: false });
}

export const EventProvider = ({ children }: { children: ReactNode }) => {
  const { isAuthenticated } = useContext(AuthContext);
  const [selectedEventId, setSelectedEventId] = useState<number | null>(() =>
    getStoredEventId()
  );
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
  // Starts as true when the effect below will load the stored event, so the
  // first render doesn't show the "no event" state for a moment.
  const [isLoading, setIsLoading] = useState(
    () => isAuthenticated && getStoredEventId() != null
  );

  const clearEvent = useCallback(() => {
    clearStoredEventId();
    setSelectedEventId(null);
    setSelectedEvent(null);
  }, []);

  const selectEvent = useCallback(async (id: number) => {
    setIsLoading(true);
    setStoredEventId(id);
    setSelectedEventId(id);

    const event = await fetchEventById(id);
    if (!event) {
      clearStoredEventId();
      setSelectedEventId(null);
      setSelectedEvent(null);
      setIsLoading(false);
      return;
    }

    setSelectedEvent(event);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    if (!isAuthenticated) {
      clearEvent();
      return;
    }

    const storedId = getStoredEventId();
    if (!storedId) {
      setSelectedEventId(null);
      setSelectedEvent(null);
      return;
    }

    setSelectedEventId(storedId);
    setIsLoading(true);

    fetchEventById(storedId).then((event) => {
      if (!event) {
        clearStoredEventId();
        setSelectedEventId(null);
        setSelectedEvent(null);
      } else {
        setSelectedEvent(event);
      }
      setIsLoading(false);
    });
  }, [isAuthenticated, clearEvent]);

  return (
    <EventContext.Provider
      value={{
        selectedEventId,
        selectedEvent,
        isLoading,
        selectEvent,
        clearEvent,
      }}
    >
      {children}
    </EventContext.Provider>
  );
};
