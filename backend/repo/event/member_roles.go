package event

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) AddMemberRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	return r.InTx(ctx, func(tx *Repo) error {
		now := time.Now().UTC()
		for _, role := range roles {
			_, err := tx.db.Exec(ctx, `
				INSERT INTO event_user_roles (event_id, user_id, role, created_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING`,
				eventID, userID, role, now,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repo) GetMemberRoles(ctx context.Context, eventID, userID int) ([]models.EventRole, error) {
	rows, err := r.db.Query(ctx, `
		SELECT role	FROM event_user_roles WHERE event_id = $1 AND user_id = $2`,
		eventID, userID,
	)
	if err != nil {
		return nil, err
	}

	roles, err := pgx.CollectRows(rows, pgx.RowTo[models.EventRole])
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repo) RemoveMembers(ctx context.Context, eventID int, userIDs []int) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM event_user_roles
		WHERE event_id = $1 AND user_id = ANY($2)`,
		eventID, userIDs,
	)
	if err != nil {
		return err
	}

	return nil
}

// SetMemberRoles makes roles the member's roles in the event, in one
// transaction: roles they hold and keep stay as they were (with their
// creation time), the others are added or dropped. Empty roles leave the
// user without roles.
func (r *Repo) SetMemberRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, err := tx.db.Exec(ctx, `
			DELETE FROM event_user_roles
			WHERE event_id = $1 AND user_id = $2 AND NOT (role = ANY($3))`,
			eventID, userID, roles,
		); err != nil {
			return err
		}

		return tx.AddMemberRoles(ctx, eventID, userID, roles)
	})
}
