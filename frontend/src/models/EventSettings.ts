// EventSettings mirrors models.EventSettings: who may finish or abort a game.
// Unlike the game rules they can change at any time.
export type FinishGameBy = 'seated' | 'referees';
export type AbortGameBy = 'referees' | 'seated_before_first_round' | 'seated';
// UndoRoundBy 'seated': the seated players undo while the game is live.
export type UndoRoundBy = 'referees' | 'seated';
// PenaltiesBy 'seated': the players at the table enter penalties and false
// mahjongs too; referees always may, from the game's view.
export type PenaltiesBy = 'seated' | 'referees';

export type EventSettings = {
  finish_game: FinishGameBy;
  abort_game: AbortGameBy;
  undo_round: UndoRoundBy;
  penalties: PenaltiesBy;
  // min_rating_games: ended games a player needs for a place in the rating;
  // players with fewer are listed below it. 1 lists everyone.
  min_rating_games: number;
};

// Limits of min_rating_games, as the backend checks them.
export const MIN_MIN_RATING_GAMES = 1;
export const MAX_MIN_RATING_GAMES = 100;

// The options in the order the form shows them, from referees only to the
// most open; referees always may.
export const FINISH_GAME_OPTIONS: FinishGameBy[] = ['referees', 'seated'];
export const PENALTIES_OPTIONS: PenaltiesBy[] = ['referees', 'seated'];
export const UNDO_ROUND_OPTIONS: UndoRoundBy[] = ['referees', 'seated'];
export const ABORT_GAME_OPTIONS: AbortGameBy[] = [
  'referees',
  'seated_before_first_round',
  'seated',
];

export const DEFAULT_EVENT_SETTINGS: EventSettings = {
  finish_game: 'seated',
  abort_game: 'seated_before_first_round',
  undo_round: 'referees',
  penalties: 'seated',
  min_rating_games: 1,
};

export function sameSettings(a: EventSettings, b: EventSettings): boolean {
  return (
    a.finish_game === b.finish_game &&
    a.abort_game === b.abort_game &&
    a.undo_round === b.undo_round &&
    a.penalties === b.penalties &&
    a.min_rating_games === b.min_rating_games
  );
}

export function isValidSettings(settings: EventSettings): boolean {
  return (
    Number.isInteger(settings.min_rating_games) &&
    settings.min_rating_games >= MIN_MIN_RATING_GAMES &&
    settings.min_rating_games <= MAX_MIN_RATING_GAMES
  );
}
