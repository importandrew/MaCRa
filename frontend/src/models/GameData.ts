import { WIND_SYMBOLS } from '../utils/constants';

export type GamePlayers = { id: number; players: PlayerOrder[] };

export interface GameData {
  id: number;
  me: PlayerScore;
  right: PlayerScore;
  opp: PlayerScore;
  left: PlayerScore;
  round: RoundData;
  status: 'created' | 'live' | 'ended' | 'aborted';
}

export interface PlayerScore {
  name: string;
  wind: typeof WIND_SYMBOLS;
  score: number;
  id: number;
}

export interface RoundData {
  num: number;
  kanji: typeof WIND_SYMBOLS & number;
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
  id: number; // player id
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
