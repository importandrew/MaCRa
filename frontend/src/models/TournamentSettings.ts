// TournamentSettings mirrors models.TournamentSettings: a tournament's own
// settings (PUT /api/event/:id/tournament).
export type TimeUpRule = 'finish_hand' | 'extra_hands' | 'stop_now';

export type TournamentSettings = {
  tours: number;
  tour_minutes: number;
  // time_up: what happens at a table when the tour's time is up.
  time_up: TimeUpRule;
  // extra_hands: hands played after the one in progress, with 'extra_hands';
  // 0 otherwise.
  extra_hands: number;
  // ties: what decides between players equal in TP and GP: 'share' the
  // place, or 'break' it by tie_breakers in order, then the draw number.
  ties: 'share' | 'break';
  tie_breakers: TieBreaker[] | null;
};

export type TieBreaker = 'head_to_head' | 'first_places';

// TiesChoice is how the settings offer the tie rule: sharing, or a
// tie-break with its criteria in order (the draw number always last).
// Labels: tournament.ties.<choice>.
export type TiesChoice =
  'share' | 'h2h_first' | 'first_h2h' | 'h2h' | 'first' | 'number';

export const TIES_CHOICES: TiesChoice[] = [
  'share',
  'h2h_first',
  'first_h2h',
  'h2h',
  'first',
  'number',
];

const TIE_BREAKERS: Record<Exclude<TiesChoice, 'share'>, TieBreaker[]> = {
  h2h_first: ['head_to_head', 'first_places'],
  first_h2h: ['first_places', 'head_to_head'],
  h2h: ['head_to_head'],
  first: ['first_places'],
  number: [],
};

// tiesChoice reads the settings' tie rule as a choice.
export function tiesChoice(s: TournamentSettings): TiesChoice {
  if (s.ties !== 'break') return 'share';
  const key = (s.tie_breakers ?? []).join(',');
  const found = (
    Object.keys(TIE_BREAKERS) as Exclude<TiesChoice, 'share'>[]
  ).find((choice) => TIE_BREAKERS[choice].join(',') === key);
  return found ?? 'number';
}

// withTies sets the settings' tie rule from a choice.
export function withTies(
  s: TournamentSettings,
  choice: TiesChoice
): TournamentSettings {
  return choice === 'share'
    ? { ...s, ties: 'share', tie_breakers: [] }
    : { ...s, ties: 'break', tie_breakers: TIE_BREAKERS[choice] };
}

// TIME_UP_RULES in the order the settings list them (labels:
// tournament.timeUp.<rule>).
export const TIME_UP_RULES: TimeUpRule[] = [
  'finish_hand',
  'extra_hands',
  'stop_now',
];

// Limits, mirroring the backend (models/tournament.go).
export const MIN_TOURS = 1;
export const MAX_TOURS = 15;
export const MIN_TOUR_MINUTES = 10;
export const MAX_TOUR_MINUTES = 300;
export const MIN_EXTRA_HANDS = 1;
export const MAX_EXTRA_HANDS = 3;

export function isValidTournamentSettings(s: TournamentSettings): boolean {
  if (s.tours < MIN_TOURS || s.tours > MAX_TOURS) return false;
  if (s.tour_minutes < MIN_TOUR_MINUTES || s.tour_minutes > MAX_TOUR_MINUTES)
    return false;
  if (s.time_up === 'extra_hands')
    return s.extra_hands >= MIN_EXTRA_HANDS && s.extra_hands <= MAX_EXTRA_HANDS;
  return s.extra_hands === 0;
}

// TournamentParticipant mirrors models.TournamentParticipant.
export type TournamentParticipant = {
  user_id: number;
  name: string;
  // number: null before the draw and while it is out of date.
  number: number | null;
  withdrawn: boolean;
};

// TournamentParticipants mirrors models.TournamentParticipants
// (GET /api/event/:id/participants).
export type TournamentParticipants = {
  // drawn: the numbers belong to exactly the current players.
  drawn: boolean;
  participants: TournamentParticipant[];
};
