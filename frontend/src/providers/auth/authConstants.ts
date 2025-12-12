import { createContext, useContext } from 'react';
import { AuthCreateContext } from './authTypes';

export const AUTH_USER_DATA = 'authUserData';

export const AuthContext = createContext<AuthCreateContext>({
  isAuthenticated: false,
  signIn: () => {},
  signOut: () => {},
  navigateRef: { current: null },
});

export const useAuth = () => useContext(AuthContext);
