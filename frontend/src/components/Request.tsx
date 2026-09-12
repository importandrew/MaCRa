import { fetchLogoutApi } from '../providers/auth/authApi';
import { AUTH_USER_DATA } from '../providers/auth/authConstants';
import { showToast } from './Toast';

export type FetchResponse<T> = {
  data: T | null;
  status: number;
};

type FetchOptions = {
  isToastEnabled?: boolean;
  redirectOn401?: boolean;
};

async function globalFetch<T>(
  request: Request,
  params?: FetchOptions
): Promise<T | null> {
  const result = await fetchWithStatus<T>(request, params);
  return result.data;
}

async function fetchWithStatus<T>(
  request: Request,
  params?: FetchOptions
): Promise<FetchResponse<T>> {
  const { isToastEnabled = true, redirectOn401 = true } = params || {};

  try {
    const response = await fetch(request);

    if (response.status === 401 && redirectOn401) {
      await fetchLogoutApi();
      localStorage.removeItem(AUTH_USER_DATA);
      showToast({
        title: 'Unauthorized request!',
        description: 'Redirecting to login...',
        color: 'danger',
      });
      setTimeout(() => (window.location.href = '/login'), 3000);
      return { data: null, status: response.status };
    }

    if (!response.ok) {
      const data = await response
        .json()
        .catch(() => ({ error: 'Unknown error' }));
      if (isToastEnabled && response.status !== 403) {
        showToast({
          title: `Error: ${data.error}`,
          color: 'danger',
        });
      }
      return { data: null, status: response.status };
    }

    const data = JSON.parse((await response.text()) || '{}');
    return { data, status: response.status };
  } catch (error) {
    isToastEnabled &&
      showToast({
        title: `Error: ${error}`,
        color: 'danger',
      });
    return { data: null, status: 0 };
  }
}

async function authFetch<T>(
  request: Request,
  params?: Omit<FetchOptions, 'redirectOn401'>
): Promise<T | null> {
  return globalFetch<T>(request, { ...params, redirectOn401: true });
}

async function publicFetch<T>(
  request: Request,
  params?: Omit<FetchOptions, 'redirectOn401'>
): Promise<FetchResponse<T>> {
  return fetchWithStatus<T>(request, { ...params, redirectOn401: false });
}

async function globalBatchFetch<T>(
  requests: Request[],
  params?: FetchOptions
): Promise<T[] | null> {
  const { isToastEnabled = true } = params || {};
  try {
    const responses = await Promise.all(
      requests.map((request) => fetch(request))
    );

    const data = await Promise.all(
      responses.map((response) => response.json())
    );
    return data;
  } catch (error) {
    isToastEnabled &&
      showToast({
        title: `Error: ${error}`,
        color: 'danger',
      });
    return null;
  }
}

function PostRequest(url: string, body: unknown = {}): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

function GetRequest(url: string): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    credentials: 'include',
  });
}

function PutRequest(url: string, body: unknown = {}): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify(body),
  });
}

function PatchRequest(url: string, body: unknown = {}): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'PATCH',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

function DeleteRequest(url: string): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'DELETE',
    credentials: 'include',
  });
}

function eventPath(eventId: number, suffix: string): string {
  return `/api/event/${eventId}${suffix}`;
}

function eventGamePath(
  eventId: number,
  gameId: number,
  suffix: string
): string {
  return `/api/event/${eventId}/game/${gameId}${suffix}`;
}

export {
  globalFetch,
  authFetch,
  publicFetch,
  fetchWithStatus,
  globalBatchFetch,
  GetRequest,
  PostRequest,
  PutRequest,
  PatchRequest,
  DeleteRequest,
  eventPath,
  eventGamePath,
};
