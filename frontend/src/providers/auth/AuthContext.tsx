import { ReactNode, useCallback, useEffect, useRef, useState } from 'react';
import {
  authFetch,
  GetRequest,
  setSessionExpiredHandler,
} from '../../components/Request';
import { NavigateFunction } from 'react-router-dom';
import { UserData } from '../../models/UserData';
import { fetchLogoutApi } from './authApi';
import { AUTH_USER_DATA, AuthContext } from './authConstants';

// storeUserData keeps only the fields the UI needs. Older versions also saved
// the session token from the login response here, where any script could read
// it; the token belongs only in the HttpOnly cookie.
const storeUserData = ({ id, name }: UserData) => {
  localStorage.setItem(AUTH_USER_DATA, JSON.stringify({ id, name }));
};

// dropStoredToken rewrites user data saved by an older version without its
// token.
const dropStoredToken = () => {
  try {
    const stored = JSON.parse(localStorage.getItem(AUTH_USER_DATA) || 'null');
    if (stored && typeof stored === 'object' && 'token' in stored) {
      storeUserData(stored as UserData);
    }
  } catch {
    // Unreadable data is left as is: it holds no usable token either.
  }
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const navigateRef = useRef<NavigateFunction | null>(null);

  const isAuthenticatedFromStore = AUTH_USER_DATA in localStorage;
  const [isAuthenticated, setIsAuthenticated] = useState(
    isAuthenticatedFromStore
  );

  // A 401 found by any request signs out at once: the UI switches to the
  // signed-out state and goes to /login while the toast is still shown.
  useEffect(() => {
    setSessionExpiredHandler(() => {
      setIsAuthenticated(false);
      if (navigateRef.current) {
        // Login returns to `from` after a sign-in, as with RequireAuth.
        const { pathname, search } = window.location;
        navigateRef.current('/login', {
          state:
            pathname === '/login' ? undefined : { from: { pathname, search } },
        });
      } else {
        window.location.href = '/login';
      }
    });
    return () => setSessionExpiredHandler(null);
  }, []);

  // The stored user data outlives the auth cookie (the browser drops it after
  // its Max-Age), so check the session once on start: a 401 logs out.
  useEffect(() => {
    if (isAuthenticatedFromStore) {
      dropStoredToken();
      authFetch(GetRequest('/api/profile'), { isToastEnabled: false });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const signIn = useCallback((data: UserData) => {
    storeUserData(data);
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
