package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// ListMembers returns the event's members with their roles: everyone with a
// role, and the owner even without one (a tournament owner, or an owner who
// dropped their roles), so they stay in the list and can take a role again.
func (r *Repo) ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id AS user_id, u.name,
			COALESCE(array_agg(eur.role ORDER BY eur.role) FILTER (WHERE eur.role IS NOT NULL), '{}') AS roles
		FROM users u
		LEFT JOIN event_user_roles eur ON eur.user_id = u.id AND eur.event_id = $1
		WHERE u.id IN (SELECT user_id FROM event_user_roles WHERE event_id = $1)
			OR u.id = (SELECT owner_user_id FROM events WHERE id = $1)
		GROUP BY u.id, u.name
		ORDER BY u.id`,
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
