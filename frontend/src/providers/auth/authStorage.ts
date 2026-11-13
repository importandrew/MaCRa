import type { UserData } from '../../models/UserData';

// AUTH_USER_DATA is the localStorage key of the signed-in user (id and name
// only). AuthContext owns it; Request.ts, outside React, checks whether it
// exists. Components use useAuth().user, never the storage.
export const AUTH_USER_DATA = 'authUserData';

// readStoredUser returns the stored user, or null if there is none or it
// can't be read. User data saved by an older version together with the
// session token is rewritten without it: the token belongs only in the
// HttpOnly cookie.
export function readStoredUser(): UserData | null {
  try {
    const stored = JSON.parse(localStorage.getItem(AUTH_USER_DATA) || 'null');
    if (!stored || typeof stored.id !== 'number') return null;
    const user: UserData = { id: stored.id, name: String(stored.name ?? '') };
    if ('token' in stored) storeUser(user);
    return user;
  } catch {
    return null;
  }
}

// storeUser keeps only the fields the UI needs.
export function storeUser({ id, name }: UserData) {
  localStorage.setItem(AUTH_USER_DATA, JSON.stringify({ id, name }));
}

export function clearStoredUser() {
  localStorage.removeItem(AUTH_USER_DATA);
}

export function hasStoredUser(): boolean {
  return AUTH_USER_DATA in localStorage;
}
