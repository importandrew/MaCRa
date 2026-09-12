import { EventRole } from './EventRole';

export type EventCapability =
  | 'manage_event'
  | 'manage_members'
  | 'create_game'
  | 'abort_game'
  | 'finish_game'
  | 'undo_round'
  | 'watch_live'
  | 'watch_results'
  | 'view_ratings';

export type EventMe = {
  is_owner: boolean;
  roles: EventRole[];
  capabilities: EventCapability[];
};
