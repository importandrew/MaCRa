// Tour models mirror models/tour.go (GET /api/event/:id/tours).
export type TourStatus = 'published' | 'live' | 'finished' | 'aborted';

export type SeatingMethod = 'random';

export type Tour = {
  id: number;
  number: number;
  status: TourStatus;
  seating: SeatingMethod;
  published_at: string;
  started_at: string | null;
  finished_at: string | null;
};

export type TourSeat = {
  user_id: number;
  name: string;
  number: number;
  wind: 'east' | 'south' | 'west' | 'north';
};

export type TourTable = {
  table: number;
  game_id: number;
  // seats: east first.
  seats: TourSeat[];
};

export type TourState = {
  // current: the last tour not aborted; null before the first.
  current: Tour | null;
  tables: TourTable[];
  // next: the number of the tour to publish next; 0 while one is published
  // or live, and when all are played.
  next: number;
  // active: participants who can be seated; missing: players lacking for
  // whole tables.
  active: number;
  missing: number;
};
