import { MutableRefObject } from 'react';
import { NavigateFunction } from 'react-router-dom';
import { UserData } from '../../models/UserData';

export type AuthCreateContext = {
  // user is the signed-in user (id and name), null for a visitor.
  user: UserData | null;
  isAuthenticated: boolean;
  signIn: (data: UserData) => void;
  signOut: () => void;
  // updateUser shows changed user data (a rename) everywhere at once.
  updateUser: (data: UserData) => void;
  navigateRef: MutableRefObject<NavigateFunction | null>;
};
