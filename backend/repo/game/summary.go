package game

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Summary returns what the capability rules need to know about a game: its
// status, when it ended, how many round log entries it has and, for a
// tournament game, its tour's status. It returns
// models.ErrGameNotFound for an unknown game.
func (r *Repo) Summary(ctx context.Context, gameID int) (models.GameSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			g.status,
			g.finished_at,
			(SELECT count(*) FROM rounds WHERE game_id = g.id) AS log_size,
			t.status AS tour_status,
			COALESCE(t.status = 'live' AND e.tournament->>'time_up' = 'stop_now'
				AND t.started_at + make_interval(mins => t.minutes + g.extra_minutes) < $2, false)
				AS stopped_at_time
		FROM games g
		JOIN events e ON e.id = g.event_id
		LEFT JOIN tours t ON t.id = g.tour_id
		WHERE g.id = $1`,
		gameID, time.Now().UTC())
	if err != nil {
		return models.GameSummary{}, err
	}

	summary, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.GameSummary])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.GameSummary{}, models.ErrGameNotFound
	}
	return summary, err
}

// Abort makes a live game aborted. With onlyWithoutRounds it does so only
// while the game has no round log entries, in the same statement, so a round
// entered at the same moment can't slip in between the check and the abort;
// otherwise it returns models.ErrGameHasRounds.
func (r *Repo) Abort(ctx context.Context, gameID int, onlyWithoutRounds bool) error {
	res, err := r.db.Exec(ctx, `
		UPDATE games SET status = $2
		WHERE id = $1 AND status = $3
			AND (NOT $4 OR NOT EXISTS (SELECT 1 FROM rounds WHERE game_id = $1))`,
		gameID, models.GameAborted, models.GameLive, onlyWithoutRounds)
	if err != nil {
		return err
	}
	if res.RowsAffected() > 0 {
		return nil
	}

	// Nothing changed: say why.
	status, err := r.Status(ctx, gameID)
	if err != nil {
		return err
	}
	if status != models.GameLive {
		return models.ErrGameNotLive
	}
	return models.ErrGameHasRounds
}

// TimeLimit returns the time limit of a game of a started tournament tour:
// when its time is up (the tour's start, plus its length, plus the time
// added to the table) and the tournament's rule for then. It returns nil for
// any other game.
func (r *Repo) TimeLimit(ctx context.Context, gameID int) (*models.GameTimeLimit, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.started_at + make_interval(mins => t.minutes + g.extra_minutes) AS ends_at,
			e.tournament->>'time_up' AS time_up,
			(e.tournament->>'extra_hands')::int AS extra_hands
		FROM games g
		JOIN tours t ON t.id = g.tour_id
		JOIN events e ON e.id = g.event_id
		WHERE g.id = $1 AND t.started_at IS NOT NULL AND t.minutes IS NOT NULL`,
		gameID)
	if err != nil {
		return nil, err
	}
	limit, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.GameTimeLimit])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //nolint:nilnil // no limit: not a started tournament game
	}
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

// TimedOutGames returns the live games of live tournament tours whose time
// is up at now with the rule "the game ends at once": they have ended in the
// replay and are to be finished.
func (r *Repo) TimedOutGames(ctx context.Context, now time.Time) ([]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id
		FROM games g
		JOIN tours t ON t.id = g.tour_id
		JOIN events e ON e.id = g.event_id
		WHERE g.status = $1 AND t.status = $2
			AND e.tournament->>'time_up' = $3
			AND t.started_at + make_interval(mins => t.minutes + g.extra_minutes) < $4
		ORDER BY g.id`,
		models.GameLive, models.TourLive, models.TimeUpStopNow, now)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[int])
}
