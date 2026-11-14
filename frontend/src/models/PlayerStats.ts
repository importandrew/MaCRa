// PlayerStats mirrors models.PlayerStats: a player's statistics in one event,
// from its ended games (GET /api/event/:id/player/:playerId/stats).
export type PlayerStats = {
  player_id: number;
  name: string;
  games: number;
  // places: how many times 1st to 4th (a tie in GP shares the better place).
  places: [number, number, number, number];
  avg_tp: number;
  avg_gp: number;
  // best_gp and worst_gp: null without games.
  best_gp: number | null;
  worst_gp: number | null;
  hands: number;
  self_draws: number;
  deal_in_wins: number;
  // avg_win and max_win: 0 without wins.
  avg_win: number;
  max_win: number;
  deal_ins: number;
  draws: number;
  penalties: number;
  false_mahjongs: number;
};
