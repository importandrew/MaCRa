import type { EventRole } from './EventRole';
import type { GameRules } from './GameRules';

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
  created_at: string;
  finished_at: string | null;
  is_participating?: boolean;
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
