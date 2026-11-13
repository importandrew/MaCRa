import { ReactNode, useCallback, useContext, useEffect, useState } from 'react';
import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import { Event } from '../../models/Event';
import {
  FetchResponse,
  GetRequest,
  fetchWithStatus,
} from '../../components/Request';
import { showToast } from '../../components/Toast';
import { AuthContext } from '../auth/authConstants';
import {
  clearStoredEventId,
  EventContext,
  getStoredEventId,
  setStoredEventId,
} from './eventConstants';

async function fetchEventById(id: number): Promise<FetchResponse<Event>> {
  const request = GetRequest(`/api/event/${id}`);
  return fetchWithStatus<Event>(request, { isToastEnabled: false });
}

// isGone tells whether a failed load of the selected event means the user
// can't have it selected any more: no access (403) or no such event (404).
// Any other failure (network, 500) keeps the selection.
function isGone(status: number): boolean {
  return status === 403 || status === 404;
}

// toastNoAccess tells the user, once, why their selected event is gone.
function toastNoAccess(t: TFunction, name?: string) {
  showToast({
    title: name
      ? t('events.noAccessTo', { name })
      : t('events.noAccessToSelected'),
    variant: 'warning',
  });
}

export const EventProvider = ({ children }: { children: ReactNode }) => {
  const { t } = useTranslation();
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

  // dropInaccessibleEvent clears the selection once the user has lost access
  // to the selected event (a 403 from it), and says so once.
  const dropInaccessibleEvent = useCallback(() => {
    toastNoAccess(t, selectedEvent?.name);
    clearEvent();
  }, [t, selectedEvent, clearEvent]);

  const selectEvent = useCallback(
    async (id: number) => {
      setIsLoading(true);
      setStoredEventId(id);
      setSelectedEventId(id);

      const { data: event, status } = await fetchEventById(id);
      if (!event) {
        if (isGone(status)) {
          toastNoAccess(t);
        }
        clearEvent();
        setIsLoading(false);
        return;
      }

      setSelectedEvent(event);
      setIsLoading(false);
    },
    [t, clearEvent]
  );

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

    fetchEventById(storedId).then(({ data: event, status }) => {
      if (event) {
        setSelectedEvent(event);
      } else if (isGone(status)) {
        // The name is unknown: the event can't be loaded any more.
        toastNoAccess(t);
        clearEvent();
      }
      // Any other failure keeps the selection; the pages show their errors.
      setIsLoading(false);
    });
    // t is left out: a language switch mustn't reload the event.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, clearEvent]);

  return (
    <EventContext.Provider
      value={{
        selectedEventId,
        selectedEvent,
        isLoading,
        selectEvent,
        clearEvent,
        dropInaccessibleEvent,
      }}
    >
      {children}
    </EventContext.Provider>
  );
};
