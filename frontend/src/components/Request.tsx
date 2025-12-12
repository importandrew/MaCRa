import { toast } from 'react-toastify';
import { fetchLogoutApi } from '../providers/auth/authApi';
import { AUTH_USER_DATA } from '../providers/auth/authConstants';

async function globalFetch(request: Request, options = {}) {
  const response = await fetch(request, options);

  // *** GLOBAL 401 HANDLING LOGIC ***
  if (response.status === 401) {
    await fetchLogoutApi();
    localStorage.removeItem(AUTH_USER_DATA);
    toast.error('Unauthorised request! Redirecting to login...');
    setTimeout(() => (window.location.href = '/login'), 3000);
    throw new Error('Unauthorized access, redirecting to login.');
  }

  return response;
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

export { globalFetch, GetRequest, PostRequest };
