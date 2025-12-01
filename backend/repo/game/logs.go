package game

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

type RoundLog struct {
	ID             int            `db:"id"`
	Round          int            `db:"round_number"`
	EventSequence  int            `db:"event_sequence"`
	Outcome        models.Outcome `db:"outcome"`
	Winner         sql.NullInt64  `db:"winner_id"`
	Loser          sql.NullInt64  `db:"loser_id"`
	Value          sql.NullInt64  `db:"value"`
	ForceNextRound bool           `db:"force_next_round"`
	CreatedAt      time.Time      `db:"created_at"`
	EastScore      sql.NullInt64  `db:"east_score"`
	SouthScore     sql.NullInt64  `db:"south_score"`
	WestScore      sql.NullInt64  `db:"west_score"`
	NorthScore     sql.NullInt64  `db:"north_score"`
}

func (r *Repo) GameLog(ctx context.Context, gameID int) ([]models.RoundLog, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
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
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.RoundLog{}, nil
		}
		return nil, err
	}
	defer rows.Close()

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
		}

		if row.Winner.Valid {
			log.Winner = int(row.Winner.Int64)
		}
		if row.Loser.Valid {
			log.Loser = int(row.Loser.Int64)
		}
		if row.Value.Valid {
			log.Value = int(row.Value.Int64)
		}
		if row.EastScore.Valid {
			log.EastScore = int(row.EastScore.Int64)
		}
		if row.SouthScore.Valid {
			log.SouthScore = int(row.SouthScore.Int64)
		}
		if row.WestScore.Valid {
			log.WestScore = int(row.WestScore.Int64)
		}
		if row.NorthScore.Valid {
			log.NorthScore = int(row.NorthScore.Int64)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
