import type { EventRole } from './EventRole';

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
  created_at: string;
  finished_at: string | null;
  is_participating?: boolean;
  is_owner?: boolean;
  roles?: EventRole[];
};

export const CONTENT_REACH_OPTIONS: { value: ContentReach; label: string }[] = [
  { value: 'public', label: 'Everyone' },
  { value: 'members', label: 'Participants only' },
  { value: 'staff', label: 'Staff only' },
];

export function reachLabel(reach: ContentReach): string {
  return CONTENT_REACH_OPTIONS.find((o) => o.value === reach)?.label ?? reach;
}
