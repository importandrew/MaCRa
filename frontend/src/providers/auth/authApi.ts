import { showToast } from '../../components/Toast';
import { PostRequest } from '../../components/Request';

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
