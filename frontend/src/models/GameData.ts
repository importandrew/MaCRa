import { WIND_SYMBOLS } from '../utils/constants';

export type GamePlayers = { id: number; players: PlayerOrder[] };

export interface GameData {
  id: number;
  me: PlayerScore;
  right: PlayerScore;
  opp: PlayerScore;
  left: PlayerScore;
  round: number;
}

export interface PlayerScore {
  id: number;
  name: string;
  wind: string;
  score: number;
}

export interface RoundData {
  num: number;
  wind_num: string;
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
