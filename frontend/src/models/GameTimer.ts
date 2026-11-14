// GameTimer mirrors models.GameTimer: a started tournament game's timer.
export type GameTimer = {
  ends_at: string;
  time_up: 'finish_hand' | 'extra_hands' | 'stop_now';
  // last_hand: the last hand to play once the time is up; missing before.
  last_hand?: number;
  // stopped: the time is up and the game ends at once: no more entries, the
  // game is to be finished.
  stopped?: boolean;
};
