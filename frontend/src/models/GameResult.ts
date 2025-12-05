import { Outcome } from './GameData';
import { Player } from './Player';

export interface GameResult {
  id: number;
  player_results: {
    name: string;
    id: number;
    gp: number;
    tp: number;
  }[];
}

export interface GameRoundLogs {
  id: number;
  players: Player[];
  round_logs: RoundLog[];
}

export interface RoundLog {
  id: number;
  round: number;
  outcome: Outcome;
  winner: number;
  loser: number;
  value: number;
  east_score: number;
  south_score: number;
  west_score: number;
  north_score: number;
}

export interface PreparedGameRoundLogs {
  players: Player[];
  rounds: PreparedRoundLog[];
}

export interface PreparedRoundLog {
  round: number;
  outcome: Outcome | null;
  initial_scores: [number, number, number, number];
  result: number | null;
  score_changes: [number, number, number, number] | null;
}
