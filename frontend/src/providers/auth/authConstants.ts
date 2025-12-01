import { useContext } from 'react';
import { AuthContext } from './AuthContext';

export const AUTH_USER_DATA = 'authUserData';

export const useAuth = () => useContext(AuthContext);
