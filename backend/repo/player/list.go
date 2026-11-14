package player

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// likeEscaper escapes LIKE's wildcards, so a name with % or _ is searched as
// typed.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// SearchByName returns up to limit users whose name contains query, ignoring
// case: exact matches first, then names starting with it, then the rest;
// within each, by name and ID.
func (r *Repo) SearchByName(ctx context.Context, query string, limit int) ([]models.Player, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name FROM users
		WHERE name ILIKE '%' || @pattern || '%'
		ORDER BY
			CASE
				WHEN lower(name) = lower(@query) THEN 0
				WHEN name ILIKE @pattern || '%' THEN 1
				ELSE 2
			END,
			name, id
		LIMIT @limit`,
		pgx.NamedArgs{
			"query":   query,
			"pattern": likeEscaper.Replace(query),
			"limit":   limit,
		},
	)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Player])
}
