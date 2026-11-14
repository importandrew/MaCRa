package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// ListParticipants returns a tournament's players with their numbers. The
// numbers count (Drawn) only while they belong to exactly the current
// players; otherwise they are left out and the draw is to be redone.
func (r *Repo) ListParticipants(ctx context.Context, eventID int) (models.TournamentParticipants, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id AS user_id, u.name, tp.number, tp.withdrawn_at IS NOT NULL AS withdrawn
		FROM event_user_roles eur
		JOIN users u ON u.id = eur.user_id
		LEFT JOIN tournament_participants tp ON tp.event_id = eur.event_id AND tp.user_id = eur.user_id
		WHERE eur.event_id = $1 AND eur.role = $2
		ORDER BY tp.number NULLS LAST, u.name, u.id`,
		eventID, models.EventPlayer)
	if err != nil {
		return models.TournamentParticipants{}, err
	}
	participants, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TournamentParticipant])
	if err != nil {
		return models.TournamentParticipants{}, err
	}

	rows, err = r.db.Query(ctx, `SELECT count(*) FROM tournament_participants WHERE event_id = $1`, eventID)
	if err != nil {
		return models.TournamentParticipants{}, err
	}
	numbered, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return models.TournamentParticipants{}, err
	}

	drawn := len(participants) > 0 && numbered == len(participants)
	for _, p := range participants {
		drawn = drawn && p.Number != nil
	}
	if !drawn {
		for i := range participants {
			participants[i].Number = nil
		}
	}
	return models.TournamentParticipants{Drawn: drawn, Participants: participants}, nil
}

// tournamentState locks the event row (FOR UPDATE with forUpdate, FOR SHARE
// otherwise) and returns whether it is a tournament and when it started, or
// models.ErrEventNotFound / models.ErrEventAlreadyFinished.
func (r *Repo) tournamentState(ctx context.Context, eventID int, forUpdate bool) (bool, *time.Time, error) {
	lock := "FOR SHARE"
	if forUpdate {
		lock = "FOR UPDATE"
	}
	rows, err := r.db.Query(ctx, `
		SELECT type, started_at, finished_at FROM events WHERE id = $1 `+lock, eventID)
	if err != nil {
		return false, nil, fmt.Errorf("lock event: %w", err)
	}

	type eventRow struct {
		Type       models.EventType `db:"type"`
		StartedAt  *time.Time       `db:"started_at"`
		FinishedAt *time.Time       `db:"finished_at"`
	}
	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[eventRow])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil, models.ErrEventNotFound
	case err != nil:
		return false, nil, fmt.Errorf("lock event: %w", err)
	case event.FinishedAt != nil:
		return false, nil, models.ErrEventAlreadyFinished
	}
	return event.Type == models.EventTournament, event.StartedAt, nil
}

// Draw gives the tournament's current players the numbers 1 to N in a random
// order, replacing an earlier draw. It refuses a regular event
// (models.ErrNotTournament), a tournament whose first tour has started
// (models.ErrParticipantsClosed) or without players (models.ErrNoPlayers).
// The event row is locked FOR UPDATE, so two draws don't mix.
func (r *Repo) Draw(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		isTournament, startedAt, err := tx.tournamentState(ctx, eventID, true)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		case startedAt != nil:
			return models.ErrParticipantsClosed
		}

		if _, err := tx.db.Exec(ctx, `DELETE FROM tournament_participants WHERE event_id = $1`, eventID); err != nil {
			return fmt.Errorf("delete old draw: %w", err)
		}
		res, err := tx.db.Exec(ctx, `
			INSERT INTO tournament_participants (event_id, user_id, number)
			SELECT event_id, user_id, row_number() OVER (ORDER BY random())
			FROM event_user_roles WHERE event_id = $1 AND role = $2`,
			eventID, models.EventPlayer)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
		if res.RowsAffected() == 0 {
			return models.ErrNoPlayers
		}
		return nil
	})
}

// Withdraw marks a participant of a tournament in progress as withdrawn:
// they keep their games and number and aren't seated any more. It refuses a
// tournament not started yet (models.ErrTournamentNotStarted; before the
// start a player is simply removed) and a user without a number or already
// withdrawn (models.ErrNotParticipant).
func (r *Repo) Withdraw(ctx context.Context, eventID, userID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		isTournament, startedAt, err := tx.tournamentState(ctx, eventID, false)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		case startedAt == nil:
			return models.ErrTournamentNotStarted
		}

		res, err := tx.db.Exec(ctx, `
			UPDATE tournament_participants SET withdrawn_at = $1
			WHERE event_id = $2 AND user_id = $3 AND withdrawn_at IS NULL`,
			time.Now().UTC(), eventID, userID)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return fmt.Errorf("user %d: %w", userID, models.ErrNotParticipant)
		}
		return nil
	})
}

// checkPlayersOpen returns models.ErrParticipantsClosed if the event is a
// tournament whose first tour has started: its players can't change any
// more. Call it on a Repo from InTx, after the event row is locked.
func (r *Repo) checkPlayersOpen(ctx context.Context, eventID int) error {
	rows, err := r.db.Query(ctx, `
		SELECT type = $2 AND started_at IS NOT NULL FROM events WHERE id = $1`,
		eventID, models.EventTournament)
	if err != nil {
		return err
	}
	closed, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return err
	}
	if closed {
		return models.ErrParticipantsClosed
	}
	return nil
}
