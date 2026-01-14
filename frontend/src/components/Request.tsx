import { fetchLogoutApi } from '../providers/auth/authApi';
import { AUTH_USER_DATA } from '../providers/auth/authConstants';
import { addToast } from '@heroui/react';

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
      addToast({
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
        addToast({
          title: `Error: ${data.error}`,
          color: 'danger',
        });
      return null;
    }

    const data = JSON.parse((await response.text()) || '{}');
    return data;
  } catch (error) {
    isToastEnabled &&
      addToast({
        title: 'Error occured: ' + error,
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

function DeleteRequest(url: string): Request {
  return new Request(import.meta.env.VITE_API_ADDR + url, {
    method: 'DELETE',
    credentials: 'include',
  });
}

export { globalFetch, GetRequest, PostRequest, DeleteRequest };
