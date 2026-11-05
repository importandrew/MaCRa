import { ReactNode, useCallback, useEffect, useRef, useState } from 'react';
import { authFetch, GetRequest } from '../../components/Request';
import { UserData } from '../../models/UserData';
import { fetchLogoutApi } from './authApi';
import { AUTH_USER_DATA, AuthContext } from './authConstants';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const navigateRef = useRef<((path: string) => void) | null>(null);

  const isAuthenticatedFromStore = AUTH_USER_DATA in localStorage;
  const [isAuthenticated, setIsAuthenticated] = useState(
    isAuthenticatedFromStore
  );

  // The stored user data outlives the auth cookie (the browser drops it after
  // its Max-Age), so check the session once on start: a 401 logs out.
  useEffect(() => {
    if (isAuthenticatedFromStore) {
      authFetch(GetRequest('/api/profile'), { isToastEnabled: false });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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
