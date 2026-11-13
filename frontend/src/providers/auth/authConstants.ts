import { createContext, useContext } from 'react';
import { AuthCreateContext } from './authTypes';

export const AuthContext = createContext<AuthCreateContext>({
  user: null,
  isAuthenticated: false,
  signIn: () => {},
  signOut: () => {},
  updateUser: () => {},
  navigateRef: { current: null },
});

export const useAuth = () => useContext(AuthContext);
