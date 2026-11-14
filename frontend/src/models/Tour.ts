// Tour models mirror models/tour.go (GET /api/event/:id/tours).
// draft: generated, seen by the staff only; published: the players see
// their tables.
export type TourStatus =
  'draft' | 'published' | 'live' | 'finished' | 'aborted';

export type SeatingMethod = 'random';

export type Tour = {
  id: number;
  number: number;
  status: TourStatus;
  seating: SeatingMethod;
  // published_at: null while the seating is a draft.
  published_at: string | null;
  // minutes: the tour's length, set when it starts.
  minutes: number | null;
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
  // status: the table's game: created (tour not started), live, ended,
  // aborted.
  status: 'created' | 'live' | 'ended' | 'aborted';
  // extra_minutes: time a referee added to the table.
  extra_minutes: number;
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
