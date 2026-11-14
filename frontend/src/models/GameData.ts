import type { GameRoundLogs } from './GameResult';
import type { GameRules } from './GameRules';

export interface GameData {
  id: number;
  me: PlayerScore;
  right: PlayerScore;
  opp: PlayerScore;
  left: PlayerScore;
  round: RoundData;
  // log_size is the number of round log entries; a new entry sends it back so
  // that the server refuses an entry made from an outdated table.
  log_size: number;
  // rules: the rules the game is played by.
  rules: GameRules;
  status: 'created' | 'live' | 'ended' | 'aborted';
  // log is the round log of the same state as the table.
  log: GameRoundLogs;
  // tour: the tour and table of a tournament game; missing otherwise.
  tour?: { tour: number; table: number };
}

export interface PlayerScore {
  name: string;
  // A wind symbol, e.g. '東'.
  wind: string;
  score: number;
  id: number;
}

export interface RoundData {
  num: number;
  // Round wind and number, e.g. '東1'.
  kanji: string;
}

export enum Outcome {
  Draw = 'draw',
  Self = 'self',
  Deal = 'deal-in',
  Penalty = 'penalty',
  FalseMahjong = 'false mahjong',
}

export interface HistoryProps {
  playerId?: number;
}

export interface GamesHistoryData {
  game_ids: number[];
}
