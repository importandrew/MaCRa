export type EventRatingEntry = {
  // rank: none for a player below the event's minimum games.
  rank?: number;
  player_id: number;
  name: string;
  games: number;
  tp: number;
  gp: number;
};

// EventRatingsData: the ranked players, and below them, without a place,
// the players with fewer than min_games ended games, in the same order.
export type EventRatingsData = {
  ratings: EventRatingEntry[];
  below_minimum: EventRatingEntry[];
  min_games: number;
};
