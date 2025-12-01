package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// type Log struct {
// 	ID        int        `json:"id"`
// 	Players   []Player   `json:"players"`
// 	RoundLogs []RoundLog `json:"round_logs"`
// }

// type Player struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// type RoundLog struct {
// 	ID         int       `json:"id"`
// 	Round      int       `json:"round"`
// 	Outcome    string    `json:"outcome"`
// 	Winner     int       `json:"winner"`
// 	Loser      int       `json:"loser"`
// 	Value      int       `json:"value"`
// 	EastScore  int       `json:"east_score"`
// 	SouthScore int       `json:"south_score"`
// 	WestScore  int       `json:"west_score"`
// 	NorthScore int       `json:"north_score"`
// 	CreatedAt  time.Time `json:"created_at"`
// }

func (s *Service) Logs(ctx context.Context, id int) (models.GameLog, error) {
	// maybe player header should be obtained by separate api call
	players, err := s.Players(ctx, id)
	if err != nil {
		return models.GameLog{}, err
	}

	roundLogs, err := s.gameRepo.GameLog(ctx, id)
	if err != nil {
		return models.GameLog{}, err
	}

	playerHeaders := make([]models.PlayerHeader, 0, len(players.Players))
	for _, player := range players.Players {
		playerHeaders = append(playerHeaders, models.PlayerHeader(player))
	}

	gameLog := models.GameLog{
		ID:        id,
		Players:   playerHeaders,
		RoundLogs: roundLogs,
	}

	return gameLog, nil
}
