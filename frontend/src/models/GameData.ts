export type GamePlayers = { id: number; players: PlayerOrder[] };

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
  status: 'created' | 'live' | 'ended' | 'aborted';
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

export interface PlayerOrder {
  id: number;
  name: string;
}

export enum Outcome {
  Draw = 'draw',
  Self = 'self',
  Deal = 'deal-in',
  Penalty = 'penalty',
  FalseMahjong = 'false mahjong',
}

export interface AnimatedScoreProps {
  scoreData: ScoreDataByPlayerIds;
  playerId: number;
  className?: string;
}

export interface ScoreDataByPlayerIds {
  [playerId: number]: number;
}

export interface HistoryProps {
  playerId?: number;
}

export interface GamesHistoryData {
  game_ids: number[];
}

export interface GameResultsData {
  id: number;
  finished_at: string;
  can_undo_last_round: boolean;
  player_results: PlayerResults[];
}

export interface PlayerResults {
  name: string;
  id: number;
  tp: number;
  gp: number;
}
