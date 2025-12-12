import { ReactNode, useCallback, useRef, useState } from 'react';
import { UserData } from '../../models/UserData';
import { fetchLogoutApi } from './authApi';
import { AUTH_USER_DATA, AuthContext } from './authConstants';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const navigateRef = useRef<((path: string) => void) | null>(null);

  const isAuthenticatedFromStore = AUTH_USER_DATA in localStorage;
  const [isAuthenticated, setIsAuthenticated] = useState(
    isAuthenticatedFromStore
  );

  const signIn = useCallback((data: UserData) => {
    localStorage.setItem(AUTH_USER_DATA, JSON.stringify(data));
    setIsAuthenticated(true);
  }, []);

  const signOut = useCallback(async () => {
    const isSignOut = await fetchLogoutApi();

    if (!isSignOut) return;

    localStorage.removeItem(AUTH_USER_DATA);
    setIsAuthenticated(false);

    if (navigateRef.current) {
      navigateRef.current('/login');
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, signIn, signOut, navigateRef }}
    >
      {children}
    </AuthContext.Provider>
  );
};
