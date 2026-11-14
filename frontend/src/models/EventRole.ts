// substitute: in a tournament only, never together with player.
export type EventRole = 'player' | 'referee' | 'observer' | 'substitute';

export type EventMember = {
  user_id: number;
  name: string;
  roles: EventRole[];
};
