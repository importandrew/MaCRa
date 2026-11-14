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
	EnteredBy      *int           `db:"entered_by"`
	EnteredByName  *string        `db:"entered_by_name"`
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
			r.entered_by,
			eu.name AS entered_by_name,
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
		LEFT JOIN users eu ON eu.id = r.entered_by
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
			EnteredBy:      playerHeader(row.EnteredBy, row.EnteredByName),
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// playerHeader returns the user of a nullable user column with their name,
// nil for NULL.
func playerHeader(id *int, name *string) *models.PlayerHeader {
	if id == nil {
		return nil
	}
	header := models.PlayerHeader{ID: *id}
	if name != nil {
		header.Name = *name
	}
	return &header
}

// UndoneEntries returns the game's undone round log entries, in the order
// they were undone.
func (r *Repo) UndoneEntries(ctx context.Context, gameID int) ([]models.UndoneEntry, error) {
	type undoneRow struct {
		Round          int            `db:"round_number"`
		Outcome        models.Outcome `db:"outcome"`
		Winner         *int           `db:"winner_id"`
		Loser          *int           `db:"loser_id"`
		Value          *int           `db:"value"`
		ForceNextRound bool           `db:"force_next_round"`
		EnteredBy      *int           `db:"entered_by"`
		EnteredByName  *string        `db:"entered_by_name"`
		EnteredAt      time.Time      `db:"entered_at"`
		UndoneBy       *int           `db:"undone_by"`
		UndoneByName   *string        `db:"undone_by_name"`
		UndoneAt       time.Time      `db:"undone_at"`
		ScoresBefore   []int          `db:"scores_before"`
		ScoresAfter    []int          `db:"scores_after"`
	}

	rows, err := r.db.Query(ctx, `
		SELECT u.round_number, u.outcome, u.winner_id, u.loser_id, u.value,
			u.force_next_round, u.entered_by, eu.name AS entered_by_name, u.entered_at,
			u.undone_by, uu.name AS undone_by_name, u.undone_at,
			u.scores_before, u.scores_after
		FROM undone_rounds u
		LEFT JOIN users eu ON eu.id = u.entered_by
		LEFT JOIN users uu ON uu.id = u.undone_by
		WHERE u.game_id = $1
		ORDER BY u.undone_at, u.id`, gameID)
	if err != nil {
		return nil, err
	}

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[undoneRow])
	if err != nil {
		return nil, err
	}

	entries := make([]models.UndoneEntry, 0, len(dbRows))
	for _, row := range dbRows {
		entries = append(entries, models.UndoneEntry{
			Round:          row.Round,
			Outcome:        row.Outcome,
			Winner:         valueOrZero(row.Winner),
			Loser:          valueOrZero(row.Loser),
			Value:          valueOrZero(row.Value),
			ForceNextRound: row.ForceNextRound,
			EnteredBy:      playerHeader(row.EnteredBy, row.EnteredByName),
			EnteredAt:      row.EnteredAt,
			UndoneBy:       playerHeader(row.UndoneBy, row.UndoneByName),
			UndoneAt:       row.UndoneAt,
			ScoresBefore:   [4]int(row.ScoresBefore),
			ScoresAfter:    [4]int(row.ScoresAfter),
		})
	}
	return entries, nil
}

// valueOrZero returns *v, or 0 for a NULL (nil) column.
func valueOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
