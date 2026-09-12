import { createContext, useContext } from 'react';
import { Event } from '../../models/Event';

export const SELECTED_EVENT_ID = 'selectedEventId';

export function getStoredEventId(): number | null {
  const raw = localStorage.getItem(SELECTED_EVENT_ID);
  if (!raw) return null;
  const id = parseInt(raw, 10);
  return Number.isNaN(id) || id <= 0 ? null : id;
}

export function setStoredEventId(id: number): void {
  localStorage.setItem(SELECTED_EVENT_ID, String(id));
}

export function clearStoredEventId(): void {
  localStorage.removeItem(SELECTED_EVENT_ID);
}

export type EventContextValue = {
  selectedEventId: number | null;
  selectedEvent: Event | null;
  isLoading: boolean;
  selectEvent: (id: number) => Promise<void>;
  clearEvent: () => void;
};

export const EventContext = createContext<EventContextValue>({
  selectedEventId: null,
  selectedEvent: null,
  isLoading: false,
  selectEvent: async () => {},
  clearEvent: () => {},
});

export const useEvent = () => useContext(EventContext);
