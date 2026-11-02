import { useCallback, useRef, useState } from 'react';

type Settled = { key: unknown; isLoaded: boolean };

/**
 * useLoadingState tracks the `isLoading` flag of a data hook so that only the
 * first load shows a spinner. LoadingTemplate and similar wrappers replace the
 * whole page with a spinner while `isLoading` is true, so raising it on every
 * refetch made pages disappear and reappear ("flicker"). After a successful
 * load, a refetch keeps the loaded data on screen and swaps it in when it
 * arrives.
 *
 * `key` identifies what is loaded (an event id, a game id, …). `isLoading` is
 * true while no load for the current key has finished, from the very first
 * render (no empty frame before the fetch starts) and right after the key
 * changes. A retry after a failed load shows the spinner again. `enabled` is
 * false when the hook has nothing to fetch.
 */
export function useLoadingState(key: unknown = null, enabled = true) {
  const [settled, setSettled] = useState<Settled | null>(null);
  const settledRef = useRef<Settled | null>(null);

  const update = useCallback((value: Settled | null) => {
    settledRef.current = value;
    setSettled(value);
  }, []);

  const startLoading = useCallback(() => {
    const current = settledRef.current;
    if (current && Object.is(current.key, key) && !current.isLoaded) {
      update(null);
    }
  }, [key, update]);

  // finishLoading ends a load; isLoaded tells whether data for key arrived.
  const finishLoading = useCallback(
    (isLoaded: boolean) => update({ key, isLoaded }),
    [key, update]
  );

  const isLoading = enabled && !(settled && Object.is(settled.key, key));

  return { isLoading, startLoading, finishLoading };
}
