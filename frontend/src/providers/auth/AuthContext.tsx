import { ReactNode, useCallback, useEffect, useRef, useState } from 'react';
import {
  authFetch,
  GetRequest,
  setSessionExpiredHandler,
} from '../../components/Request';
import { NavigateFunction } from 'react-router-dom';
import { UserData } from '../../models/UserData';
import { fetchLogoutApi } from './authApi';
import { AuthContext } from './authConstants';
import { clearStoredUser, readStoredUser, storeUser } from './authStorage';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const navigateRef = useRef<NavigateFunction | null>(null);

  // user is the one source of truth for who is signed in; the stored copy
  // only carries it across page loads.
  const [user, setUser] = useState<UserData | null>(() => readStoredUser());
  const isAuthenticated = user != null;

  // A 401 found by any request signs out at once (Request.ts has cleared the
  // stored user): the UI switches to the signed-out state and goes to /login
  // while the toast is still shown.
  useEffect(() => {
    setSessionExpiredHandler(() => {
      setUser(null);
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

  const updateUser = useCallback((data: UserData) => {
    storeUser(data);
    setUser({ id: data.id, name: data.name });
  }, []);

  // The stored user outlives the auth cookie (the browser drops it after its
  // Max-Age), so check the session once on start: a 401 signs out. The
  // answer also brings the current name (it may have changed on another
  // device).
  useEffect(() => {
    if (!readStoredUser()) return;
    authFetch<UserData>(GetRequest('/api/profile'), {
      isToastEnabled: false,
    }).then((profile) => {
      if (profile) updateUser(profile);
    });
  }, [updateUser]);

  const signIn = useCallback(
    (data: UserData) => {
      updateUser(data);
    },
    [updateUser]
  );

  const signOut = useCallback(async () => {
    const isSignOut = await fetchLogoutApi();

    if (!isSignOut) return;

    clearStoredUser();
    setUser(null);

    if (navigateRef.current) {
      navigateRef.current('/login');
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        signIn,
        signOut,
        updateUser,
        navigateRef,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
