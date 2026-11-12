import { useCallback } from 'react';
import { EventCapability, EventMe } from '../../models/EventMe';
import { useFetchEventMe } from './eventApi';

export function hasCapability(
  me: EventMe | null,
  cap: EventCapability
): boolean {
  return me?.capabilities.includes(cap) ?? false;
}

export function useEventCapabilities(eventId: number | null) {
  const { data, isError, isForbidden, isNotFound, isLoading, refetch } =
    useFetchEventMe(eventId);

  const check = useCallback(
    (cap: EventCapability) => hasCapability(data, cap),
    [data]
  );

  return {
    me: data,
    isError,
    isForbidden,
    isNotFound,
    isLoading,
    refetch,
    hasCapability: check,
  };
}
