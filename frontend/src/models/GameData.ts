import { WIND_SYMBOLS } from '../utils/constants';

export type GamePlayers = { id: number; players: PlayerOrder[] };

export interface GameData {
  id: number;
  me: PlayerScore;
  right: PlayerScore;
  opp: PlayerScore;
  left: PlayerScore;
  round: RoundData;
}

export interface PlayerScore {
  name: string;
  wind: typeof WIND_SYMBOLS;
  score: number;
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
