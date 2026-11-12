package game

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

type RoundLog struct {
	ID             int            `db:"id"`
	Round          int            `db:"round_number"`
	EventSequence  int            `db:"event_sequence"`
	Outcome        models.Outcome `db:"outcome"`
	Winner         *int           `db:"winner_id"`
	Loser          *int           `db:"loser_id"`
	Value          *int           `db:"value"`
	ForceNextRound bool           `db:"force_next_round"`
	CreatedAt      time.Time      `db:"created_at"`
	EastScore      *int           `db:"east_score"`
	SouthScore     *int           `db:"south_score"`
	WestScore      *int           `db:"west_score"`
	NorthScore     *int           `db:"north_score"`
}

func (r *Repo) GameLog(ctx context.Context, gameID int) ([]models.RoundLog, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			r.id,
			r.round_number,
			r.event_sequence,
			r.outcome,
			r.winner_id,
			r.loser_id,
			r.value,
			r.force_next_round,
			r.created_at,
			(SELECT rs.score FROM round_scores rs 
			 JOIN game_players gp ON rs.player_id = gp.player_id AND gp.game_id = r.game_id
			 WHERE rs.round_id = r.id AND gp.starting_wind = 'east') as east_score,
			(SELECT rs.score FROM round_scores rs 
			 JOIN game_players gp ON rs.player_id = gp.player_id AND gp.game_id = r.game_id
			 WHERE rs.round_id = r.id AND gp.starting_wind = 'south') as south_score,
			(SELECT rs.score FROM round_scores rs 
			 JOIN game_players gp ON rs.player_id = gp.player_id AND gp.game_id = r.game_id
			 WHERE rs.round_id = r.id AND gp.starting_wind = 'west') as west_score,
			(SELECT rs.score FROM round_scores rs 
			 JOIN game_players gp ON rs.player_id = gp.player_id AND gp.game_id = r.game_id
			 WHERE rs.round_id = r.id AND gp.starting_wind = 'north') as north_score
		FROM rounds r
		WHERE r.game_id = $1
		ORDER BY r.round_number, r.event_sequence`, gameID)
	if err != nil {
		return nil, err
	}

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[RoundLog])
	if err != nil {
		return nil, err
	}

	// Convert to RoundLog models
	logs := make([]models.RoundLog, 0, len(dbRows))
	for _, row := range dbRows {
		log := models.RoundLog{
			ID:             row.ID,
			Round:          row.Round,
			EventSequence:  row.EventSequence,
			Outcome:        row.Outcome,
			ForceNextRound: row.ForceNextRound,
			CreatedAt:      row.CreatedAt,
			Winner:         valueOrZero(row.Winner),
			Loser:          valueOrZero(row.Loser),
			Value:          valueOrZero(row.Value),
			EastScore:      valueOrZero(row.EastScore),
			SouthScore:     valueOrZero(row.SouthScore),
			WestScore:      valueOrZero(row.WestScore),
			NorthScore:     valueOrZero(row.NorthScore),
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// valueOrZero returns *v, or 0 for a NULL (nil) column.
func valueOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
