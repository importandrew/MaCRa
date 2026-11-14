import type { Event } from './Event';

// EventInvite mirrors models.EventInvite: what an invite link shows before
// joining.
export type EventInvite = {
  event_id: number;
  name: string;
  type: Event['type'];
  // is_member: the user is already a member (false for a visitor).
  is_member: boolean;
};
