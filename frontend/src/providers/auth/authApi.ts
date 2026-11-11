import i18n from '../../i18n';
import { showToast } from '../../components/Toast';
import { GetRequest, PostRequest, PutRequest } from '../../components/Request';
import { authFetch } from '../../components/Request';
import { UserData } from '../../models/UserData';

export async function fetchLogoutApi(): Promise<boolean> {
  try {
    const logoutRequest = PostRequest('/auth/logout');
    const response = await fetch(logoutRequest);

    return response.ok;
  } catch (error) {
    showToast({
      title: i18n.t('errors.logoutNetwork', { error: String(error) }),
      variant: 'danger',
    });
    return false;
  }
}

export async function fetchProfileApi(): Promise<UserData | null> {
  const profileRequest = GetRequest('/api/profile');
  const data = await authFetch<{ id: number; name: string }>(profileRequest, {
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
  const data = await authFetch<{ id: number; name: string }>(updateRequest, {
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
