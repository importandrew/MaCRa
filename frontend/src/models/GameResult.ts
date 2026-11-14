import { Outcome } from './GameData';
import { Player } from './Player';

export interface GameResult {
  id: number;
  finished_at: string; // ISO 8601 from backend
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
  // undone: the entries undone so far, in the order they were undone; the
  // log shows each one at the time of its undo.
  undone: UndoneEntry[];
}

// LogUser is who entered or undid an entry: a player or a referee. null when
// unknown (an entry from before authors were recorded, a deleted account).
export type LogUser = { id: number; name: string } | null;

// UndoneEntry mirrors models.UndoneEntry: an undone round log entry.
export interface UndoneEntry {
  round: number;
  outcome: Outcome;
  winner: number;
  loser: number;
  value: number;
  force_next_round: boolean;
  entered_by: LogUser;
  entered_at: string;
  undone_by: LogUser;
  undone_at: string;
  // scores_before and scores_after: the players' scores before and after
  // the entry, by starting seat (east, south, west, north).
  scores_before: [number, number, number, number];
  scores_after: [number, number, number, number];
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
  created_at: string;
  entered_by: LogUser;
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

export interface PlayerStatus {
  has_last: boolean;
  has_current: boolean;
}
