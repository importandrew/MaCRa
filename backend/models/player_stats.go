package models

// PlayerStats is a player's statistics in one event, from its ended games.
// Percentages are left to the client: it has the counts and the totals.
type PlayerStats struct {
	PlayerID int    `json:"player_id"`
	Name     string `json:"name"`

	// Games is the number of ended games; Places counts the places taken
	// (1st to 4th; a tie in GP shares the better place).
	Games  int     `json:"games"`
	Places [4]int  `json:"places"`
	AvgTP  float64 `json:"avg_tp"`
	AvgGP  float64 `json:"avg_gp"`
	// BestGP and WorstGP are the best and the worst game's GP; nil without
	// games.
	BestGP  *int `json:"best_gp"`
	WorstGP *int `json:"worst_gp"`

	// Hands is the number of hands played in those games.
	Hands      int `json:"hands"`
	SelfDraws  int `json:"self_draws"`
	DealInWins int `json:"deal_in_wins"`
	// AvgWin and MaxWin are the average and the biggest winning hand
	// value; 0 without wins.
	AvgWin float64 `json:"avg_win"`
	MaxWin int     `json:"max_win"`
	// DealIns is how many times the player dealt in to someone's win.
	DealIns       int `json:"deal_ins"`
	Draws         int `json:"draws"`
	Penalties     int `json:"penalties"`
	FalseMahjongs int `json:"false_mahjongs"`
}
