// Tour models mirror models/tour.go (GET /api/event/:id/tours).
// draft: generated, seen by the staff only; published: the players see
// their tables.
export type TourStatus =
  'draft' | 'published' | 'live' | 'finished' | 'aborted';

// fixed: the tour's tables imported from a file.
export type SeatingMethod = 'random' | 'swiss' | 'interval' | 'fixed';

// SEATING_METHODS in the order the desk offers them (labels:
// tours.method.<method>); fixed only for a tour with tables from a file.
export const SEATING_METHODS: SeatingMethod[] = [
  'random',
  'swiss',
  'interval',
  'fixed',
];

// Seating is how a tour is to be seated (generateTour, setTourPlan): the
// method, an interval step (undefined: the number of tables), and for Swiss
// and interval whether to avoid repeat meetings.
export type Seating = {
  seating: SeatingMethod;
  step?: number;
  avoid_repeats: boolean;
};

export type Tour = {
  id: number;
  number: number;
  status: TourStatus;
  seating: SeatingMethod;
  // seating_step: an interval seating's step; null otherwise.
  seating_step: number | null;
  // avoid_repeats: a Swiss or interval seating moved players apart to avoid
  // repeat meetings.
  avoid_repeats: boolean;
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
  // sitter: who sits at the seat instead of its player (a substitute); null
  // while the player does.
  sitter: { id: number; name: string } | null;
  // forfeit: the player gets 0 TP and 0 GP for the game.
  forfeit: boolean;
};

export type TourTable = {
  table: number;
  game_id: number;
  // status: the table's game: created (tour not started), live, ended,
  // aborted.
  status: 'created' | 'live' | 'ended' | 'aborted';
  // extra_minutes: time a referee added to the table.
  extra_minutes: number;
  // log_size: the game's log entries, sent back with a change at the table.
  log_size: number;
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
  // plan: each tour's planned seating, 1 to the number of tours; file: the
  // tables imported from a file, null without them.
  plan: TourPlan[];
  file: SeatingFileInfo | null;
  // substitutes: the tournament's substitutes, with the table each sits at
  // now (0: none).
  substitutes: Substitute[];
};

export type Substitute = { id: number; name: string; table: number };

// TourPlan mirrors models.TourPlan: the seating a tour is to have.
export type TourPlan = {
  number: number;
  seating: SeatingMethod;
  seating_step: number | null;
  avoid_repeats: boolean;
  // file_tables: the tour's tables from a file (0 without them).
  file_tables: number;
  // seated: the tour's seating is published or played; its plan is final.
  seated: boolean;
};

// SeatingFileInfo mirrors models.SeatingFileInfo: the tables from a file,
// their tours, the most tables in a tour, the participant numbers they seat,
// and their repeat meetings among their own tours.
export type SeatingFileInfo = {
  tours: number[];
  tables: number;
  players: number;
  repeats: number;
  max_met: number;
};
