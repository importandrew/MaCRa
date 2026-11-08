import { PlayerScore, RoundData } from './GameData';

export interface GameView {
  id: number;
  status: 'created' | 'live' | 'ended' | 'aborted';
  round: RoundData;
  // log_size is the number of round log entries, as in GameData.
  log_size: number;
  player_scores: PlayerScore[];
}
