package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT eur.user_id, u.name, array_agg(eur.role ORDER BY eur.role) AS roles
		FROM event_user_roles eur
		JOIN users u ON u.id = eur.user_id
		WHERE eur.event_id = $1
		GROUP BY eur.user_id, u.name
		ORDER BY eur.user_id`,
		eventID,
	)
	if err != nil {
		return nil, err
	}

	members, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.EventMember])
	if err != nil {
		return nil, err
	}

	return members, nil
}
