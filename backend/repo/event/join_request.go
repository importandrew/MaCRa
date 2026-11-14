package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Joining by a request or the invite link counts the event's players against
// its limit, so it runs under one more lock than the other member changes
// (member_locks.go): after the event row FOR SHARE,
// pg_advisory_xact_lock(event_id), the one-key lock of the event's player
// count. Two users can't both take the last place. The one-key lock never
// meets the two-key member locks: PostgreSQL keeps the two key spaces apart.

// joinEvent is what joining checks about an event, read under lockJoinEvent.
type joinEvent struct {
	ID             int                    `db:"id"`
	OwnerUserID    int                    `db:"owner_user_id"`
	Visibility     models.EventVisibility `db:"visibility"`
	AcceptRequests bool                   `db:"accept_requests"`
	PlayerLimit    *int                   `db:"player_limit"`
	FinishedAt     *time.Time             `db:"finished_at"`
}

// Conditions lockJoinEvent finds the event by.
const (
	byEventID    = `id = $1`
	byInviteLink = `invite_token = $1`
)

// lockJoinEvent locks the event row FOR SHARE, then its player count, and
// returns the event, or notFound when no event matches. where is byEventID
// or byInviteLink. Call it on a Repo from InTx.
func (r *Repo) lockJoinEvent(ctx context.Context, where string, arg any, notFound error) (joinEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, owner_user_id, visibility, accept_requests, player_limit, finished_at
		FROM events WHERE `+where+` FOR SHARE`, arg)
	if err != nil {
		return joinEvent{}, fmt.Errorf("lock event: %w", err)
	}

	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[joinEvent])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return joinEvent{}, notFound
	case err != nil:
		return joinEvent{}, fmt.Errorf("lock event: %w", err)
	}

	if _, err := r.db.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, event.ID); err != nil {
		return joinEvent{}, fmt.Errorf("lock player count of event %d: %w", event.ID, err)
	}
	return event, nil
}

// checkRoom returns models.ErrEventFull if the event has as many players as
// its limit (nil: no limit). Call it after lockJoinEvent.
func (r *Repo) checkRoom(ctx context.Context, eventID int, limit *int) error {
	if limit == nil {
		return nil
	}
	count, err := r.CountPlayers(ctx, eventID)
	if err != nil {
		return err
	}
	if count >= *limit {
		return models.ErrEventFull
	}
	return nil
}

// UpdateJoining sets whether the event takes requests to join and its
// player limit, refused in a finished event (models.ErrEventAlreadyFinished).
func (r *Repo) UpdateJoining(ctx context.Context, eventID int, settings models.JoiningSettings) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, err := tx.lockOpenEvent(ctx, eventID); err != nil {
			return err
		}
		_, err := tx.db.Exec(ctx, `
			UPDATE events SET accept_requests = $1, player_limit = $2 WHERE id = $3`,
			settings.AcceptRequests, settings.PlayerLimit, eventID)
		return err
	})
}

// JoinRequestStatus returns the status of the user's request to join the
// event, empty without one.
func (r *Repo) JoinRequestStatus(ctx context.Context, eventID, userID int) (models.JoinRequestStatus, error) {
	rows, err := r.db.Query(ctx, `
		SELECT status FROM event_join_requests WHERE event_id = $1 AND user_id = $2`,
		eventID, userID)
	if err != nil {
		return "", err
	}

	status, err := pgx.CollectOneRow(rows, pgx.RowTo[models.JoinRequestStatus])
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return status, err
}

// CreateJoinRequest stores the user's pending request to join; asking again
// while it is pending changes nothing. It refuses an event that doesn't take
// requests (finished, private or switched off: models.ErrRequestsClosed), a
// member (models.ErrAlreadyMember), a declined user
// (models.ErrRequestDeclined) and a full event (models.ErrEventFull).
func (r *Repo) CreateJoinRequest(ctx context.Context, eventID, userID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		event, err := tx.lockJoinEvent(ctx, byEventID, eventID, models.ErrEventNotFound)
		if err != nil {
			return err
		}
		if event.FinishedAt != nil || event.Visibility != models.EventVisibilityPublic || !event.AcceptRequests {
			return models.ErrRequestsClosed
		}

		roles, err := tx.GetMemberRoles(ctx, eventID, userID)
		if err != nil {
			return err
		}
		if len(roles) > 0 || event.OwnerUserID == userID {
			return models.ErrAlreadyMember
		}

		status, err := tx.JoinRequestStatus(ctx, eventID, userID)
		switch {
		case err != nil:
			return err
		case status == models.JoinRequestDeclined:
			return models.ErrRequestDeclined
		case status == models.JoinRequestPending:
			return nil
		}

		if err := tx.checkRoom(ctx, eventID, event.PlayerLimit); err != nil {
			return err
		}

		_, err = tx.db.Exec(ctx, `
			INSERT INTO event_join_requests (event_id, user_id, status, created_at)
			VALUES ($1, $2, $3, $4)`,
			eventID, userID, models.JoinRequestPending, time.Now().UTC())
		return err
	})
}

// WithdrawJoinRequest deletes the user's pending request, or returns
// models.ErrRequestNotFound.
func (r *Repo) WithdrawJoinRequest(ctx context.Context, eventID, userID int) error {
	res, err := r.db.Exec(ctx, `
		DELETE FROM event_join_requests
		WHERE event_id = $1 AND user_id = $2 AND status = $3`,
		eventID, userID, models.JoinRequestPending)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return models.ErrRequestNotFound
	}
	return nil
}

// ListJoinRequests returns the event's pending requests, oldest first.
func (r *Repo) ListJoinRequests(ctx context.Context, eventID int) ([]models.JoinRequest, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.user_id, u.name, r.created_at
		FROM event_join_requests r
		JOIN users u ON u.id = r.user_id
		WHERE r.event_id = $1 AND r.status = $2
		ORDER BY r.created_at, r.user_id`,
		eventID, models.JoinRequestPending)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.JoinRequest])
}

// CountJoinRequests returns the number of the event's pending requests.
func (r *Repo) CountJoinRequests(ctx context.Context, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT count(*) FROM event_join_requests WHERE event_id = $1 AND status = $2`,
		eventID, models.JoinRequestPending)
	if err != nil {
		return 0, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[int])
}

// AcceptJoinRequest makes the user a player and deletes their request. It
// refuses a finished event (models.ErrEventAlreadyFinished), a request that
// isn't pending (models.ErrRequestNotFound) and a full event
// (models.ErrEventFull). It doesn't matter whether the event still takes
// requests: the owner answers the pending ones anyway.
func (r *Repo) AcceptJoinRequest(ctx context.Context, eventID, userID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		event, err := tx.lockJoinEvent(ctx, byEventID, eventID, models.ErrEventNotFound)
		if err != nil {
			return err
		}
		if event.FinishedAt != nil {
			return models.ErrEventAlreadyFinished
		}

		status, err := tx.JoinRequestStatus(ctx, eventID, userID)
		if err != nil {
			return err
		}
		if status != models.JoinRequestPending {
			return fmt.Errorf("user %d: %w", userID, models.ErrRequestNotFound)
		}

		if err := tx.checkRoom(ctx, eventID, event.PlayerLimit); err != nil {
			return err
		}
		// insertRoles deletes the pending request.
		return tx.insertRoles(ctx, eventID, userID, []models.EventRole{models.EventPlayer})
	})
}

// DeclineJoinRequest marks the user's pending request declined (the user
// can't ask again), or returns models.ErrRequestNotFound.
func (r *Repo) DeclineJoinRequest(ctx context.Context, eventID, userID int) error {
	res, err := r.db.Exec(ctx, `
		UPDATE event_join_requests SET status = $1
		WHERE event_id = $2 AND user_id = $3 AND status = $4`,
		models.JoinRequestDeclined, eventID, userID, models.JoinRequestPending)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("user %d: %w", userID, models.ErrRequestNotFound)
	}
	return nil
}
