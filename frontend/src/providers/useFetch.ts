import { useCallback, useEffect, useRef, useState } from 'react';
import type { FetchResponse } from '../components/Request';
import type { ApiError } from '../models/ApiError';
import { useLoadingState } from './useLoadingState';

// FetchKey says what a data hook loads (an event id, "event/game", …): a new
// key loads again. null means there is nothing to load.
export type FetchKey = string | number | null;

type Loaded<T> = { key: FetchKey; value: T };
type Failure = { key: FetchKey; status: number; error?: ApiError };

// useFetch is the one data-loading hook the API hooks are built on. load
// makes the request (it gets the data shown now, e.g. to keep a 404 quiet
// once something is on screen); the hook loads on mount and whenever key
// changes, and refetch loads again.
//
// - isLoading is true only until the first load for the key has finished
//   (useLoadingState): a refetch keeps the loaded data on screen. While a new
//   key loads, the previous data stays (pages show their spinner then).
// - A failed load sets isForbidden (403), isNotFound (404) or isError (any
//   other failure, a network error included); error is the server's body.
//   A failure belongs to its key: a new key starts without one.
// - A response that arrives after a newer request was started (the key
//   changed, or refetch was called again) is ignored.
// - mutate replaces the data with a value the caller already has (e.g. the
//   table a round entry answered with), without a request.
export function useFetch<T>(
  key: FetchKey,
  load: (shown: T | null) => Promise<FetchResponse<T>>
) {
  const { isLoading, startLoading, finishLoading } = useLoadingState(
    key,
    key != null
  );
  const [loaded, setLoaded] = useState<Loaded<T> | null>(null);
  const [failure, setFailure] = useState<Failure | null>(null);
  const data = key == null ? null : (loaded?.value ?? null);

  // The latest load and data, for the loads below to use without depending
  // on them: load is a new function on every render.
  const loadRef = useRef(load);
  const dataRef = useRef(data);
  useEffect(() => {
    loadRef.current = load;
    dataRef.current = data;
  });
  const requestRef = useRef(0);

  // request loads for key; state changes only after the response, so the
  // effect below may call it.
  const request = useCallback(async (): Promise<FetchResponse<T>> => {
    const id = ++requestRef.current;
    if (key == null) return { data: null, status: 0 };

    const result = await loadRef.current(dataRef.current);
    if (id !== requestRef.current) return result;

    if (result.data == null) {
      setFailure({ key, status: result.status, error: result.error });
      finishLoading(false);
    } else {
      setLoaded({ key, value: result.data });
      setFailure(null);
      finishLoading(true);
    }
    return result;
  }, [key, finishLoading]);

  useEffect(() => {
    void request();
  }, [request]);

  // refetch loads again; after a failure it shows the spinner again.
  const refetch = useCallback(() => {
    startLoading();
    setFailure(null);
    return request();
  }, [startLoading, request]);

  const mutate = useCallback(
    (value: T) => {
      if (key != null) setLoaded({ key, value });
    },
    [key]
  );

  const status = failure?.key === key ? failure.status : undefined;
  return {
    data,
    isLoading,
    isForbidden: status === 403,
    isNotFound: status === 404,
    isError: status != null && status !== 403 && status !== 404,
    // isFailed: any failure, for hooks that don't tell 403 or 404 apart.
    isFailed: status != null,
    error: failure?.key === key ? failure.error : undefined,
    refetch,
    mutate,
  };
}
