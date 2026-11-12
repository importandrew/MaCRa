import { PlayerScore, RoundData } from './GameData';
import type { GameRoundLogs } from './GameResult';

export interface GameView {
  id: number;
  status: 'created' | 'live' | 'ended' | 'aborted';
  round: RoundData;
  // log_size is the number of round log entries, as in GameData.
  log_size: number;
  player_scores: PlayerScore[];
  // log is the round log of the same state, as in GameData.
  log: GameRoundLogs;
}
