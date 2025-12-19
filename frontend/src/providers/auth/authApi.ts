import { toast } from 'react-toastify';
import { PostRequest } from '../../components/Request';

export async function fetchLogoutApi(): Promise<boolean> {
  try {
    const logoutRequest = PostRequest('/auth/logout');
    const response = await fetch(logoutRequest);

    return response.ok;
  } catch (error) {
    toast('Network error during logout: ' + error);
    return false;
  }
}
