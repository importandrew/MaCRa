package player

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) Get(ctx context.Context, id int) (models.Player, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.Player{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT id, name FROM users WHERE id = $1`, id)
	if err != nil {
		return models.Player{}, err
	}

	player, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Player])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Player{}, models.ErrPlayerNotFound
		}
		return models.Player{}, err
	}

	return player, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (models.Player, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.Player{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT id, name FROM users WHERE email = $1`,
		email,
	)
	if err != nil {
		return models.Player{}, err
	}

	player, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Player])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Player{}, models.ErrPlayerNotFound
		}
		return models.Player{}, err
	}

	return player, nil
}
