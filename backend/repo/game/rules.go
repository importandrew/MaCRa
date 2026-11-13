package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Rules returns the rules the game started with, or models.ErrGameNotFound.
func (r *Repo) Rules(ctx context.Context, gameID int) (models.GameRules, error) {
	rows, err := r.db.Query(ctx, `SELECT rules FROM games WHERE id = $1`, gameID)
	if err != nil {
		return models.GameRules{}, err
	}

	rules, err := pgx.CollectOneRow(rows, pgx.RowTo[models.GameRules])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.GameRules{}, models.ErrGameNotFound
		}
		return models.GameRules{}, err
	}

	return rules, nil
}
