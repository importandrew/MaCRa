import type { EventRole } from './EventRole';
import type { GameRules } from './GameRules';
import type { EventSettings } from './EventSettings';
import type { TournamentSettings } from './TournamentSettings';

export type EventVisibility = 'private' | 'public';

export type ContentReach = 'public' | 'members' | 'staff';

export type Event = {
  id: number;
  type: 'tournament' | 'regular';
  name: string;
  description?: string | null;
  owner_user_id: number;
  visibility: EventVisibility;
  ratings_reach: ContentReach;
  games_reach: ContentReach;
  rules: GameRules;
  settings: EventSettings;
  // tournament: a tournament's own settings; null for a regular event.
  tournament: TournamentSettings | null;
  // started_at: when a tournament's first tour started; null before.
  started_at: string | null;
  // accept_requests: the owner takes requests to join (public events only).
  accept_requests: boolean;
  // player_limit: the most players the event takes by a request or the
  // invite link; null for no limit.
  player_limit: number | null;
  created_at: string;
  finished_at: string | null;
  is_participating?: boolean;
  // last_game_at: when its last game (not aborted) started; lists only.
  last_game_at?: string | null;
  is_owner?: boolean;
  roles?: EventRole[];
  // player_count comes only with GET /api/event/:id (the overview): non-members
  // don't get the member list, so the overview can't count it.
  player_count?: number;
  // rules_frozen also comes only with the overview: the event has games, so
  // its rules can't change.
  rules_frozen?: boolean;
};

// CONTENT_REACHES in the order the pickers show them (labels: reach.* texts).
export const CONTENT_REACHES: ContentReach[] = ['public', 'members', 'staff'];

// EventStatus mirrors models.EventStatus: a tournament is prepared, then in
// progress from its first tour, then finished; a regular event is in progress
// until it is finished.
export type EventStatus = 'preparation' | 'in_progress' | 'finished';

export function eventStatus(event: Event): EventStatus {
  if (event.finished_at) return 'finished';
  if (event.type === 'tournament' && !event.started_at) return 'preparation';
  return 'in_progress';
}
