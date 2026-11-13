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
    isAuthenticated ? getStoredEventId() : null
  );
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
  // Starts as true when the effect below will load the stored event, so the
  // first render doesn't show the "no event" state for a moment.
  const [isLoading, setIsLoading] = useState(
    () => isAuthenticated && getStoredEventId() != null
  );

  // Signing in or out starts the selection afresh: from the stored event
  // after signing in, none after signing out. Adjusted while rendering (the
  // React pattern for state that follows a changed input), not in an effect.
  const [authSeen, setAuthSeen] = useState(isAuthenticated);
  if (authSeen !== isAuthenticated) {
    setAuthSeen(isAuthenticated);
    const storedId = isAuthenticated ? getStoredEventId() : null;
    setSelectedEventId(storedId);
    setSelectedEvent(null);
    setIsLoading(storedId != null);
  }

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

  // After signing out, the stored selection goes too.
  useEffect(() => {
    if (!isAuthenticated) clearStoredEventId();
  }, [isAuthenticated]);

  // Load the stored event on start and after signing in. State changes only
  // when the answer comes.
  useEffect(() => {
    const storedId = isAuthenticated ? getStoredEventId() : null;
    if (!storedId) return;

    let isCurrent = true;
    fetchEventById(storedId).then(({ data: event, status }) => {
      if (!isCurrent) return;
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
    return () => {
      isCurrent = false;
    };
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
