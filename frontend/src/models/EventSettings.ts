// EventSettings mirrors models.EventSettings: who may finish or abort a game.
// Unlike the game rules they can change at any time.
export type FinishGameBy = 'seated' | 'referees';
export type AbortGameBy = 'referees' | 'seated_before_first_round' | 'seated';

export type EventSettings = {
  finish_game: FinishGameBy;
  abort_game: AbortGameBy;
};

// The options in the order the form shows them, from referees only to the
// most open; referees always may.
export const FINISH_GAME_OPTIONS: FinishGameBy[] = ['referees', 'seated'];
export const ABORT_GAME_OPTIONS: AbortGameBy[] = [
  'referees',
  'seated_before_first_round',
  'seated',
];

export const DEFAULT_EVENT_SETTINGS: EventSettings = {
  finish_game: 'seated',
  abort_game: 'seated_before_first_round',
};

export function sameSettings(a: EventSettings, b: EventSettings): boolean {
  return a.finish_game === b.finish_game && a.abort_game === b.abort_game;
}
