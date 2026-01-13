import { PlayerScore, RoundData } from './GameData';

export interface GameView {
  id: number;
  status: 'created' | 'live' | 'ended' | 'aborted';
  round: RoundData;
  player_scores: PlayerScore[];
}
