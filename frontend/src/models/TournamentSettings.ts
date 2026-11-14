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
};

// TIME_UP_RULES in the order the settings list them (labels:
// tournament.timeUp.<rule>).
export const TIME_UP_RULES: TimeUpRule[] = [
  'finish_hand',
  'extra_hands',
  'stop_now',
];

// Limits, mirroring the backend (models/tournament.go).
export const MIN_TOURS = 3;
export const MAX_TOURS = 15;
export const MIN_TOUR_MINUTES = 30;
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
