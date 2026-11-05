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

// Several requests can fail with 401 at once; log out and redirect only once.
// The flag stays set until the next sign-in (resetSessionExpired).
let isSessionExpiredHandled = false;
let sessionExpiredHandler: (() => void) | null = null;

// setSessionExpiredHandler lets the auth provider switch to the signed-out
// state and navigate to /login without a page reload, so the toast stays.
function setSessionExpiredHandler(handler: (() => void) | null) {
  sessionExpiredHandler = handler;
}

// resetSessionExpired re-arms the session-expired handling after a sign-in.
function resetSessionExpired() {
  isSessionExpiredHandled = false;
}

// handleSessionExpired signs out and sends the user to /login. wasSignedIn
// picks the toast: an expired session, or a visitor who never signed in and
// hit an endpoint that needs an account.
function handleSessionExpired(wasSignedIn: boolean) {
  if (isSessionExpiredHandled) return;
  isSessionExpiredHandled = true;

  localStorage.removeItem(AUTH_USER_DATA);
  showToast(
    wasSignedIn
      ? {
          title: 'Session expired',
          description: 'Please log in again.',
          variant: 'danger',
        }
      : {
          title: 'Login required',
          description: 'Please log in to continue.',
          variant: 'warning',
        }
  );

  if (sessionExpiredHandler) {
    sessionExpiredHandler();
  } else {
    window.location.href = '/login';
  }

  // The UI doesn't wait for this; it only makes sure the cookie is gone.
  void fetchLogoutApi();
}

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
    // A 401 from a public endpoint means the backend cleared a stale auth
    // cookie; an anonymous visitor then just repeats the request without it.
    const retryRequest = request.clone();
    let response = await fetch(request);

    if (response.status === 401) {
      const isSignedIn = AUTH_USER_DATA in localStorage;

      if (redirectOn401 || isSignedIn || isSessionExpiredHandled) {
        handleSessionExpired(isSignedIn);
        return { data: null, status: response.status };
      }

      response = await fetch(retryRequest);
    }

    if (!response.ok) {
      const data = await response
        .json()
        .catch(() => ({ error: 'Unknown error' }));
      if (isToastEnabled && response.status !== 403) {
        showToast({
          title: `Error: ${data.error}`,
          variant: 'danger',
        });
      }
      return { data: null, status: response.status };
    }

    const data = JSON.parse((await response.text()) || '{}');
    return { data, status: response.status };
  } catch (error) {
    if (isToastEnabled) {
      showToast({
        title: `Error: ${error}`,
        variant: 'danger',
      });
    }
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

// globalBatchFetch runs the requests in parallel like publicFetch, so a 401
// gets the same session-expired handling. A failed request yields null in its
// slot; a network failure fails the whole batch with a single toast.
async function globalBatchFetch<T>(
  requests: Request[],
  params?: Omit<FetchOptions, 'redirectOn401'>
): Promise<(T | null)[] | null> {
  const { isToastEnabled = true } = params || {};

  const results = await Promise.all(
    requests.map((request) =>
      fetchWithStatus<T>(request, {
        isToastEnabled: false,
        redirectOn401: false,
      })
    )
  );

  if (results.some((result) => result.status === 0)) {
    if (isToastEnabled) {
      showToast({ title: 'Error: network request failed', variant: 'danger' });
    }
    return null;
  }

  return results.map((result) => result.data);
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
  setSessionExpiredHandler,
  resetSessionExpired,
  GetRequest,
  PostRequest,
  PutRequest,
  PatchRequest,
  DeleteRequest,
  eventPath,
  eventGamePath,
};
