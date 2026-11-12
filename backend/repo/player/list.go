package player

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error) {
	nameFilter := ""
	if name != "" {
		nameFilter = " WHERE name ILIKE '%' || @name || '%'"
	}

	query := fmt.Sprintf(`SELECT id, name FROM users %s ORDER BY name`, nameFilter)
	if limit > 0 {
		query += ` LIMIT @limit`
	}

	rows, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"name":  strings.TrimSpace(name),
			"limit": limit,
		},
	)
	if err != nil {
		return nil, err
	}

	players, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Player])
	if err != nil {
		return nil, err
	}

	return players, nil
}
