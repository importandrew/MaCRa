import { globalFetch, PostRequest } from '../../components/Request';

export async function fetchLogoutApi(): Promise<boolean> {
  try {
    const logoutRequest = PostRequest('/auth/logout');
    const response = await globalFetch(logoutRequest);

    return response.ok;
  } catch (error) {
    console.error('Network error during logout:', error);
    return false;
  }
}
