export type EventRatingEntry = {
  rank: number;
  player_id: number;
  name: string;
  games: number;
  tp: number;
  gp: number;
};

export type EventRatingsData = {
  ratings: EventRatingEntry[];
};
