package event

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) Create(ctx context.Context, event models.Event) (int, error) {
	var id int
	err := r.InTx(ctx, func(tx *Repo) error {
		now := time.Now().UTC()
		args := pgx.NamedArgs{
			"owner_id":      event.OwnerUserID,
			"type":          event.Type,
			"name":          event.Name,
			"desc":          event.Description,
			"visibility":    event.Visibility,
			"ratings_reach": event.RatingsReach,
			"games_reach":   event.GamesReach,
			"rules":         event.Rules,
			"now":           now,
		}

		rows, err := tx.db.Query(ctx, `
			INSERT INTO events (
				type, name, description, owner_user_id,
				visibility, ratings_reach, games_reach, rules, created_at
			)
			VALUES (
				@type, @name, @desc, @owner_id,
				@visibility, @ratings_reach, @games_reach, @rules, @now
			)
			RETURNING id`,
			args,
		)
		if err != nil {
			return err
		}

		id, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
		if err != nil {
			return err
		}

		if event.Type == models.EventRegular {
			ownerDefaultRoles := []models.EventRole{models.EventPlayer, models.EventReferee}
			for _, role := range ownerDefaultRoles {
				_, err = tx.db.Exec(ctx, `
					INSERT INTO event_user_roles (event_id, user_id, role, created_at)
					VALUES ($1, $2, $3, $4)`,
					id, event.OwnerUserID, role, now,
				)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
