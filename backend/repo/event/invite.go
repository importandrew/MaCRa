package event

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// InviteToken returns the event's invite token, nil when the link is off.
func (r *Repo) InviteToken(ctx context.Context, eventID int) (*string, error) {
	rows, err := r.db.Query(ctx, `SELECT invite_token FROM events WHERE id = $1`, eventID)
	if err != nil {
		return nil, err
	}

	token, err := pgx.CollectOneRow(rows, pgx.RowTo[*string])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrEventNotFound
	}
	return token, err
}

// SetInviteToken sets the event's invite token: a new one replaces the old
// link, nil turns it off. A finished event refuses a new token
// (models.ErrEventAlreadyFinished); turning it off is always fine.
func (r *Repo) SetInviteToken(ctx context.Context, eventID int, token *string) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if token != nil {
			if _, err := tx.lockOpenEvent(ctx, eventID); err != nil {
				return err
			}
		}

		res, err := tx.db.Exec(ctx, `UPDATE events SET invite_token = $1 WHERE id = $2`, token, eventID)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return models.ErrEventNotFound
		}
		return nil
	})
}

// EventByInviteToken returns the event the invite token belongs to, or
// models.ErrInviteNotFound.
func (r *Repo) EventByInviteToken(ctx context.Context, token string) (models.Event, error) {
	rows, err := r.db.Query(ctx, `SELECT id FROM events WHERE invite_token = $1`, token)
	if err != nil {
		return models.Event{}, err
	}

	eventID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Event{}, models.ErrInviteNotFound
	}
	if err != nil {
		return models.Event{}, err
	}

	return r.Get(ctx, eventID)
}

// JoinByInvite makes the user a player of the event the invite token belongs
// to and returns the event's ID; a player already stays as they are. It
// refuses an unknown token (models.ErrInviteNotFound), a finished event
// (models.ErrEventAlreadyFinished) and a full one (models.ErrEventFull). It
// takes the locks of joining (join_request.go), so the link can't be used
// while the event is being finished, and two users can't both take its last
// place.
func (r *Repo) JoinByInvite(ctx context.Context, token string, userID int) (int, error) {
	var eventID int
	err := r.InTx(ctx, func(tx *Repo) error {
		event, err := tx.lockJoinEvent(ctx, byInviteLink, token, models.ErrInviteNotFound)
		if err != nil {
			return err
		}
		if event.FinishedAt != nil {
			return models.ErrEventAlreadyFinished
		}

		roles, err := tx.GetMemberRoles(ctx, event.ID, userID)
		if err != nil {
			return err
		}
		if !slices.Contains(roles, models.EventPlayer) {
			if err := tx.checkRoom(ctx, event.ID, event.PlayerLimit); err != nil {
				return err
			}
		}

		eventID = event.ID
		return tx.insertRoles(ctx, eventID, userID, []models.EventRole{models.EventPlayer})
	})
	if err != nil {
		return 0, err
	}
	return eventID, nil
}
