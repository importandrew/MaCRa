import { useCallback, useEffect, useRef, useState } from 'react';
import type { FetchResponse } from '../components/Request';
import type { ApiError } from '../models/ApiError';
import { useLoadingState } from './useLoadingState';

// FetchKey says what a data hook loads (an event id, "event/game", …): a new
// key loads again. null means there is nothing to load.
export type FetchKey = string | number | null;

type Failure = { status: number; error?: ApiError };

// useFetch is the one data-loading hook the API hooks are built on. load
// makes the request (it gets the data shown now, e.g. to keep a 404 quiet
// once something is on screen); the hook loads on mount and whenever key
// changes, and refetch loads again.
//
// - isLoading is true only until the first load for the key has finished
//   (useLoadingState): a refetch keeps the loaded data on screen.
// - A failed load sets isForbidden (403), isNotFound (404) or isError (any
//   other failure, a network error included); error is the server's body.
// - A response that arrives after a newer request was started (the key
//   changed, or refetch was called again) is ignored.
export function useFetch<T>(
  key: FetchKey,
  load: (shown: T | null) => Promise<FetchResponse<T>>
) {
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    key,
    key != null
  );
  const [data, setData] = useState<T | null>(null);
  const [failure, setFailure] = useState<Failure | null>(null);

  // The latest load and data, for refetch to use without depending on them:
  // load is a new function on every render.
  const loadRef = useRef(load);
  const dataRef = useRef(data);
  useEffect(() => {
    loadRef.current = load;
    dataRef.current = data;
  });
  const requestRef = useRef(0);

  const refetch = useCallback(async (): Promise<FetchResponse<T>> => {
    const request = ++requestRef.current;
    if (key == null) {
      setData(null);
      setFailure(null);
      return { data: null, status: 0 };
    }

    startLoading();
    setFailure(null);
    const result = await loadRef.current(dataRef.current);
    if (request !== requestRef.current) {
      return result;
    }

    if (result.data == null) {
      setFailure({ status: result.status, error: result.error });
      finishLoading(false);
      return result;
    }

    setData(result.data);
    finishLoading(true);
    return result;
  }, [key, startLoading, finishLoading]);

  useEffect(() => {
    void refetch();
  }, [refetch]);

  const status = failure?.status;
  return {
    data,
    isLoading,
    isForbidden: status === 403,
    isNotFound: status === 404,
    isError: failure != null && status !== 403 && status !== 404,
    // isFailed: any failure, for hooks that don't tell 403 or 404 apart.
    isFailed: failure != null,
    error: failure?.error,
    refetch,
  };
}
