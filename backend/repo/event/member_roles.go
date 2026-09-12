package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) AddMemberRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, role := range roles {
		_, err = tx.Exec(ctx, `
			INSERT INTO event_user_roles (event_id, user_id, role)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING`,
			eventID, userID, role,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repo) GetMemberRoles(ctx context.Context, eventID, userID int) ([]models.EventRole, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
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
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		DELETE FROM event_user_roles
		WHERE event_id = $1 AND user_id = ANY($2)`,
		eventID, userIDs,
	)
	if err != nil {
		return err
	}

	return nil
}
