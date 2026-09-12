package event

import (
	"context"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) Create(ctx context.Context, event models.Event) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	args := pgx.NamedArgs{
		"owner_id":     event.OwnerUserID,
		"type":         event.Type,
		"name":         event.Name,
		"desc":         event.Description,
		"visibility":    event.Visibility,
		"ratings_reach": event.RatingsReach,
		"games_reach":   event.GamesReach,
		"now":          time.Now().UTC(),
	}

	rows, err := tx.Query(ctx, `
		INSERT INTO events (
			type, name, description, owner_user_id,
			visibility, ratings_reach, games_reach, created_at
		)
		VALUES (
			@type, @name, @desc, @owner_id,
			@visibility, @ratings_reach, @games_reach, @now
		)
		RETURNING id`,
		args,
	)
	if err != nil {
		return 0, err
	}

	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}

	if event.Type == models.EventRegular {
		ownerDefaultRoles := []models.EventRole{models.EventPlayer, models.EventReferee}
		for _, role := range ownerDefaultRoles {
			_, err = tx.Exec(ctx, `
				INSERT INTO event_user_roles (event_id, user_id, role)
				VALUES ($1, $2, $3)`,
				id, event.OwnerUserID, role,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}
