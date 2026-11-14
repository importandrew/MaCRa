// JoinRequest mirrors models.JoinRequest: a pending request to join, as the
// owner sees it.
export type JoinRequest = {
  user_id: number;
  name: string;
  created_at: string;
};

// JoiningSettings mirrors models.JoiningSettings (PUT /event/:id/joining).
export type JoiningSettings = {
  accept_requests: boolean;
  // player_limit: null for no limit.
  player_limit: number | null;
};
