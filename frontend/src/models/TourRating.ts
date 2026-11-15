// TourRating mirrors models.TourRating (GET /api/event/:id/tours/rating): a
// tournament's rating by tour.
export type TourRating = {
  // tours: the finished tours, the columns of each row's tours.
  tours: number[];
  // in_progress: the live tour; 0 when none.
  in_progress: number;
  // rows: the active players by place, those who haven't played a finished
  // tour yet, then those who withdrew.
  rows: TourRatingRow[];
};

export type TourRatingRow = {
  player_id: number;
  name: string;
  // number: the participant number; null while not drawn.
  number: number | null;
  withdrawn: boolean;
  // rank: null for a player who withdrew or hasn't played yet; players
  // equal in TP and GP share it.
  rank: number | null;
  tp: number;
  gp: number;
  // tours: the result in each finished tour (as TourRating.tours); null
  // when the player didn't play it.
  tours: (TourResult | null)[];
};

export type TourResult = {
  tp: number;
  gp: number;
  table: number;
};
