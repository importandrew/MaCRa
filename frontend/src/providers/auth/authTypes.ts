import { MutableRefObject } from 'react';
import { UserData } from '../../models/UserData';

export type AuthCreateContext = {
  isAuthenticated: boolean;
  signIn: (data: UserData) => void;
  signOut: () => void;
  navigateRef: MutableRefObject<((path: string) => void) | null>;
};
