import { fetchLogoutApi } from '../providers/auth/authApi';
import { AUTH_USER_DATA } from '../providers/auth/authConstants';
import { showToast } from './Toast';

async function globalFetch<T>(
  request: Request,
  params?: { isToastEnabled: boolean }
): Promise<T | null> {
  const { isToastEnabled = true } = params || {};
  try {
    const response = await fetch(request);

    // Global 401 Unauthorized handler
    if (response.status === 401) {
      await fetchLogoutApi();
      localStorage.removeItem(AUTH_USER_DATA);
      showToast({
        title: 'Unauthorized request!',
        description: 'Redirecting to login...',
        color: 'danger',
      });
      setTimeout(() => (window.location.href = '/login'), 3000);
      return null;
    }

    // Global default error handler
    if (!response.ok) {
      const data = await response.json();
      isToastEnabled &&
        showToast({
          title: `Error: ${data.error}`,
          color: 'danger',
        });
      return null;
    }

    const data = JSON.parse((await response.text()) || '{}');
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

async function globalBatchFetch<T>(
  requests: Request[],
  params?: { isToastEnabled: boolean }
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

function DeleteRequest(url: string): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'DELETE',
    credentials: 'include',
  });
}

export {
  globalFetch,
  globalBatchFetch,
  GetRequest,
  PostRequest,
  PutRequest,
  DeleteRequest,
};
