package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// CountPlayers returns the number of members with the player role.
func (r *Repo) CountPlayers(ctx context.Context, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT count(*) FROM event_user_roles
		WHERE event_id = $1 AND role = $2`,
		eventID, models.EventPlayer,
	)
	if err != nil {
		return 0, err
	}

	return pgx.CollectOneRow(rows, pgx.RowTo[int])
}
