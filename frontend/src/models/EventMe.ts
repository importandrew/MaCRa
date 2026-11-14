import { EventRole } from './EventRole';

// EventCapability is what the user may do in an event, whatever the game
// (GET /api/event/:id/me).
export type EventCapability =
  | 'manage_event'
  | 'manage_members'
  | 'create_game'
  | 'watch_live'
  | 'watch_results'
  | 'view_ratings';

export type EventMe = {
  is_owner: boolean;
  roles: EventRole[];
  capabilities: EventCapability[];
};

// GameCapability is what the user may do with one game
// (GET /api/event/:id/game/:gameId/me): the rules about being seated, the
// game's status and its rounds live on the backend.
export type GameCapability =
  | 'view_game'
  | 'enter_round'
  | 'enter_penalty'
  | 'referee_entry'
  | 'undo_round'
  | 'abort_game'
  | 'finish_game';

export type GameMe = {
  capabilities: GameCapability[];
};
