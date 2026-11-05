import { MutableRefObject } from 'react';
import { NavigateFunction } from 'react-router-dom';
import { UserData } from '../../models/UserData';

export type AuthCreateContext = {
  isAuthenticated: boolean;
  signIn: (data: UserData) => void;
  signOut: () => void;
  navigateRef: MutableRefObject<NavigateFunction | null>;
};
