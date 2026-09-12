export type EventRole = 'player' | 'referee' | 'observer';

export type EventMember = {
  user_id: number;
  name: string;
  roles: EventRole[];
};
