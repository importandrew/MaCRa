package event

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// AddMemberRoles adds roles to a member (or makes the user one), refused in
// a finished event (models.ErrEventAlreadyFinished).
func (r *Repo) AddMemberRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, err := tx.lockOpenEvent(ctx, eventID); err != nil {
			return err
		}
		return tx.insertRoles(ctx, eventID, userID, roles)
	})
}

// insertRoles adds the roles the user doesn't hold yet. A member has no
// pending request to join, so it deletes the user's one (a declined request
// stays: the user still can't ask again after leaving).
func (r *Repo) insertRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	if _, err := r.db.Exec(ctx, `
		DELETE FROM event_join_requests
		WHERE event_id = $1 AND user_id = $2 AND status = $3`,
		eventID, userID, models.JoinRequestPending,
	); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, role := range roles {
		_, err := r.db.Exec(ctx, `
			INSERT INTO event_user_roles (event_id, user_id, role, created_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING`,
			eventID, userID, role, now,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) GetMemberRoles(ctx context.Context, eventID, userID int) ([]models.EventRole, error) {
	rows, err := r.db.Query(ctx, `
		SELECT role	FROM event_user_roles WHERE event_id = $1 AND user_id = $2`,
		eventID, userID,
	)
	if err != nil {
		return nil, err
	}

	roles, err := pgx.CollectRows(rows, pgx.RowTo[models.EventRole])
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// RemoveMembers removes the users' roles in the event, in one transaction
// under the member locks (member_locks.go). It refuses the owner
// (models.ErrOwnerNotRemovable), a user in a live game of the event
// (models.ErrMemberInCurrentGame, with their name) and a finished event
// (models.ErrEventAlreadyFinished).
func (r *Repo) RemoveMembers(ctx context.Context, eventID int, userIDs []int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		ownerID, err := tx.lockOpenEvent(ctx, eventID)
		if err != nil {
			return err
		}
		if slices.Contains(userIDs, ownerID) {
			return fmt.Errorf("user %d: %w", ownerID, models.ErrOwnerNotRemovable)
		}
		if err := tx.lockMembers(ctx, eventID, userIDs); err != nil {
			return err
		}
		if err := tx.checkNotInLiveGame(ctx, eventID, userIDs, models.ErrMemberInCurrentGame); err != nil {
			return err
		}

		_, err = tx.db.Exec(ctx, `
			DELETE FROM event_user_roles
			WHERE event_id = $1 AND user_id = ANY($2)`,
			eventID, userIDs,
		)
		return err
	})
}

// SetMemberRoles makes roles the member's roles in the event, in one
// transaction under the member locks (member_locks.go): roles they keep stay
// as they were (with their creation time), the others are added or dropped.
// It refuses:
//   - a user who isn't a member (models.ErrNotMember); the owner always is;
//   - no roles for anyone but the owner (models.ErrInvalidRole);
//   - dropping the player role of a user in a live game of the event
//     (models.ErrPlayerRoleInGame, with their name);
//   - a finished event (models.ErrEventAlreadyFinished).
func (r *Repo) SetMemberRoles(ctx context.Context, eventID, userID int, roles []models.EventRole) error {
	return r.InTx(ctx, func(tx *Repo) error {
		ownerID, err := tx.lockOpenEvent(ctx, eventID)
		if err != nil {
			return err
		}
		if err := tx.lockMembers(ctx, eventID, []int{userID}); err != nil {
			return err
		}

		current, err := tx.GetMemberRoles(ctx, eventID, userID)
		if err != nil {
			return err
		}
		isOwner := userID == ownerID
		if len(current) == 0 && !isOwner {
			return fmt.Errorf("user %d: %w", userID, models.ErrNotMember)
		}
		if len(roles) == 0 && !isOwner {
			return fmt.Errorf("%w: a member needs a role", models.ErrInvalidRole)
		}
		if slices.Contains(current, models.EventPlayer) && !slices.Contains(roles, models.EventPlayer) {
			if err := tx.checkNotInLiveGame(ctx, eventID, []int{userID}, models.ErrPlayerRoleInGame); err != nil {
				return err
			}
		}

		if _, err := tx.db.Exec(ctx, `
			DELETE FROM event_user_roles
			WHERE event_id = $1 AND user_id = $2 AND NOT (role = ANY($3))`,
			eventID, userID, roles,
		); err != nil {
			return err
		}

		return tx.insertRoles(ctx, eventID, userID, roles)
	})
}
