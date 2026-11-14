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

// JoinRequestState is what the overview offers a non-member about joining by
// a request: ask (a visitor signs in first), the pending request (it can be
// withdrawn), a declined one, or no places left.
export type JoinRequestState = 'available' | 'pending' | 'declined' | 'full';

export type EventMe = {
  is_owner: boolean;
  roles: EventRole[];
  capabilities: EventCapability[];
  // join_request is missing when the overview offers nothing about it.
  join_request?: JoinRequestState;
  // pending_requests: requests waiting for an answer, for whoever manages the
  // members (missing when there are none).
  pending_requests?: number;
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
