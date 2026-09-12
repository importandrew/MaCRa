package event

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListGamesByStatus(
	ctx context.Context,
	eventID int,
	status models.GameStatus,
	limit, offset int,
) ([]models.EventGameItemRow, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT g.id, g.status
		FROM games g
		WHERE g.event_id = $1
			AND g.status = $2
		ORDER BY g.id DESC
		LIMIT $3 OFFSET $4`,
		eventID, status, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.EventGameItemRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return games, nil
}
