package player

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	nameFilter := ""
	if name != "" {
		nameFilter = " WHERE name ILIKE '%' || @name || '%'"
	}

	query := fmt.Sprintf(`SELECT id, name FROM users %s ORDER BY name`, nameFilter)
	if limit > 0 {
		query += ` LIMIT @limit`
	}

	rows, err := conn.Query(ctx, query,
		pgx.NamedArgs{
			"name":  strings.TrimSpace(name),
			"limit": limit,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	players, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Player])
	if err != nil {
		return nil, err
	}

	return players, nil
}
