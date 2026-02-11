import { showToast } from '../../components/Toast';
import { GetRequest, PostRequest, PutRequest } from '../../components/Request';
import { globalFetch } from '../../components/Request';
import { UserData } from '../../models/UserData';

export async function fetchLogoutApi(): Promise<boolean> {
  try {
    const logoutRequest = PostRequest('/auth/logout');
    const response = await fetch(logoutRequest);

    return response.ok;
  } catch (error) {
    showToast({
      title: `Network error during logout: ${error}`,
      color: 'danger',
    });
    return false;
  }
}

export async function fetchProfileApi(): Promise<UserData | null> {
  const profileRequest = GetRequest('/api/profile');
  const data = await globalFetch<{ id: number; name: string }>(profileRequest, {
    isToastEnabled: true,
  });

  if (!data) {
    return null;
  }

  return {
    id: data.id,
    name: data.name,
  };
}

export async function updateProfileApi(name: string): Promise<UserData | null> {
  const updateRequest = PutRequest('/api/profile', { name });
  const data = await globalFetch<{ id: number; name: string }>(updateRequest, {
    isToastEnabled: true,
  });

  if (!data) {
    return null;
  }

  return {
    id: data.id,
    name: data.name,
  };
}
