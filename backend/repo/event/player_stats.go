package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// PlayerStats returns the player's statistics in the event, from its ended
// games, or models.ErrPlayerNotFound for a user who is neither a player of
// the event nor has played in it (their name isn't shown).
func (r *Repo) PlayerStats(ctx context.Context, eventID, playerID int) (models.PlayerStats, error) {
	var stats models.PlayerStats
	err := r.InTx(ctx, func(tx *Repo) error {
		if err := tx.playerInEvent(ctx, eventID, playerID, &stats); err != nil {
			return err
		}
		if err := tx.gameStats(ctx, eventID, playerID, &stats); err != nil {
			return err
		}
		return tx.handStats(ctx, eventID, playerID, &stats)
	})
	return stats, err
}

// playerInEvent fills the player's ID and name, or returns
// models.ErrPlayerNotFound.
func (r *Repo) playerInEvent(ctx context.Context, eventID, playerID int, stats *models.PlayerStats) error {
	rows, err := r.db.Query(ctx, `
		SELECT u.name FROM users u
		WHERE u.id = $2 AND (
			EXISTS (SELECT 1 FROM event_user_roles eur
				WHERE eur.event_id = $1 AND eur.user_id = u.id AND eur.role = $3)
			OR EXISTS (SELECT 1 FROM game_players gp JOIN games g ON g.id = gp.game_id
				WHERE g.event_id = $1 AND gp.player_id = u.id))`,
		eventID, playerID, models.EventPlayer)
	if err != nil {
		return fmt.Errorf("player: %w", err)
	}

	name, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("player %d in event %d: %w", playerID, eventID, models.ErrPlayerNotFound)
	case err != nil:
		return fmt.Errorf("player: %w", err)
	}

	stats.PlayerID = playerID
	stats.Name = name
	return nil
}

// gameStats fills the per-game figures: games, places, average TP and GP,
// best and worst GP.
func (r *Repo) gameStats(ctx context.Context, eventID, playerID int, stats *models.PlayerStats) error {
	rows, err := r.db.Query(ctx, `
		WITH ranked AS (
			SELECT pr.player_id, pr.game_points, pr.tournament_points,
				1 + (SELECT count(*) FROM player_results o
					WHERE o.game_id = pr.game_id AND o.game_points > pr.game_points) AS place
			FROM player_results pr
			JOIN games g ON g.id = pr.game_id
			WHERE g.event_id = $1 AND g.status = $3
		)
		SELECT count(*) AS games,
			count(*) FILTER (WHERE place = 1) AS first,
			count(*) FILTER (WHERE place = 2) AS second,
			count(*) FILTER (WHERE place = 3) AS third,
			count(*) FILTER (WHERE place = 4) AS fourth,
			COALESCE(avg(tournament_points), 0)::float8 AS avg_tp,
			COALESCE(avg(game_points), 0)::float8 AS avg_gp,
			max(game_points) AS best_gp,
			min(game_points) AS worst_gp
		FROM ranked WHERE player_id = $2`,
		eventID, playerID, models.GameEnded)
	if err != nil {
		return fmt.Errorf("game stats: %w", err)
	}

	type gameRow struct {
		Games   int     `db:"games"`
		First   int     `db:"first"`
		Second  int     `db:"second"`
		Third   int     `db:"third"`
		Fourth  int     `db:"fourth"`
		AvgTP   float64 `db:"avg_tp"`
		AvgGP   float64 `db:"avg_gp"`
		BestGP  *int    `db:"best_gp"`
		WorstGP *int    `db:"worst_gp"`
	}
	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[gameRow])
	if err != nil {
		return fmt.Errorf("game stats: %w", err)
	}

	stats.Games = row.Games
	stats.Places = [4]int{row.First, row.Second, row.Third, row.Fourth}
	stats.AvgTP = row.AvgTP
	stats.AvgGP = row.AvgGP
	stats.BestGP = row.BestGP
	stats.WorstGP = row.WorstGP
	return nil
}

// handStats fills the per-hand figures from the round logs of the ended
// games the player sat in.
func (r *Repo) handStats(ctx context.Context, eventID, playerID int, stats *models.PlayerStats) error {
	rows, err := r.db.Query(ctx, `
		SELECT count(DISTINCT (r.game_id, r.round_number)) AS hands,
			count(*) FILTER (WHERE r.outcome = $4 AND r.winner_id = $2) AS self_draws,
			count(*) FILTER (WHERE r.outcome = $5 AND r.winner_id = $2) AS deal_in_wins,
			COALESCE(avg(r.value) FILTER (WHERE r.outcome IN ($4, $5) AND r.winner_id = $2), 0)::float8 AS avg_win,
			COALESCE(max(r.value) FILTER (WHERE r.outcome IN ($4, $5) AND r.winner_id = $2), 0) AS max_win,
			count(*) FILTER (WHERE r.outcome = $5 AND r.loser_id = $2) AS deal_ins,
			count(*) FILTER (WHERE r.outcome = $6) AS draws,
			count(*) FILTER (WHERE r.outcome = $7 AND r.loser_id = $2) AS penalties,
			count(*) FILTER (WHERE r.outcome = $8 AND r.loser_id = $2) AS false_mahjongs
		FROM rounds r
		JOIN games g ON g.id = r.game_id
		JOIN game_players gp ON gp.game_id = g.id AND gp.player_id = $2
		WHERE g.event_id = $1 AND g.status = $3`,
		eventID, playerID, models.GameEnded,
		models.Self, models.Deal, models.Draw, models.Penalty, models.FalseMahjong)
	if err != nil {
		return fmt.Errorf("hand stats: %w", err)
	}

	type handRow struct {
		Hands         int     `db:"hands"`
		SelfDraws     int     `db:"self_draws"`
		DealInWins    int     `db:"deal_in_wins"`
		AvgWin        float64 `db:"avg_win"`
		MaxWin        int     `db:"max_win"`
		DealIns       int     `db:"deal_ins"`
		Draws         int     `db:"draws"`
		Penalties     int     `db:"penalties"`
		FalseMahjongs int     `db:"false_mahjongs"`
	}
	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[handRow])
	if err != nil {
		return fmt.Errorf("hand stats: %w", err)
	}

	stats.Hands = row.Hands
	stats.SelfDraws = row.SelfDraws
	stats.DealInWins = row.DealInWins
	stats.AvgWin = row.AvgWin
	stats.MaxWin = row.MaxWin
	stats.DealIns = row.DealIns
	stats.Draws = row.Draws
	stats.Penalties = row.Penalties
	stats.FalseMahjongs = row.FalseMahjongs
	return nil
}
