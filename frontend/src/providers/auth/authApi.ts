import { addToast } from '@heroui/react';
import { PostRequest } from '../../components/Request';

export async function fetchLogoutApi(): Promise<boolean> {
  try {
    const logoutRequest = PostRequest('/auth/logout');
    const response = await fetch(logoutRequest);

    return response.ok;
  } catch (error) {
    addToast({
      title: 'Network error during logout: ' + error,
      color: 'danger',
    });
    return false;
  }
}
