// Tour models mirror models/tour.go (GET /api/event/:id/tours).
// draft: generated, seen by the staff only; published: the players see
// their tables.
export type TourStatus =
  'draft' | 'published' | 'live' | 'finished' | 'aborted';

// fixed: the tour's tables from the seating prepared for all tours (plan).
export type SeatingMethod = 'random' | 'swiss' | 'interval' | 'fixed';

// SEATING_METHODS in the order the desk offers them (labels:
// tours.method.<method>); fixed only with a plan.
export const SEATING_METHODS: SeatingMethod[] = [
  'random',
  'swiss',
  'interval',
  'fixed',
];

export type Tour = {
  id: number;
  number: number;
  status: TourStatus;
  seating: SeatingMethod;
  // seating_step: an interval seating's step; null otherwise.
  seating_step: number | null;
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
  // repeats and repeat_pairs: who at the table met before, how often.
  repeats: number;
  repeat_pairs: { a: number; b: number; met: number }[];
  // places: the range of places in the standings; null before any tour
  // is finished.
  places: { from: number; to: number } | null;
};

export type TourState = {
  // current: the last tour not aborted; null before the first.
  current: Tour | null;
  tables: TourTable[];
  // next: the number of the tour to publish next; 0 while one is published
  // or live, and when all are played.
  next: number;
  // repeats: repeat meetings at the current tour's tables; max_met: the
  // most times a pair seated together met before; players_with_repeats:
  // players who meet someone again.
  repeats: number;
  max_met: number;
  players_with_repeats: number;
  // active: participants who can be seated; missing: players lacking for
  // whole tables.
  active: number;
  missing: number;
  // plan: the seating prepared for all tours; null when there is none.
  plan: SeatingPlanInfo | null;
};

// SeatingPlanInfo mirrors models.SeatingPlanInfo: where the plan came from,
// its tours, the most tables in a tour, the participant numbers it seats,
// and its repeat meetings among its own tours.
export type SeatingPlanInfo = {
  source: 'generated' | 'imported';
  created_at: string;
  tours: number[];
  tables: number;
  players: number;
  repeats: number;
  max_met: number;
};
