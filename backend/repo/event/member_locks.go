package event

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Member changes run under the same locks as starting a game
// (repo/game.NewGame) and reviving one (repo/game.reviveGame): the event row
// FOR SHARE, then pg_advisory_xact_lock(event_id, user_id) per user in
// ascending order. A game can't start with a player who is losing the player
// role or being removed at the same moment, and a player can't lose the role
// or leave while a game with them is starting: whichever transaction comes
// second waits and sees the first one's result. Call these on a Repo from
// InTx: outside a transaction the locks are released at once.

// lockOpenEvent locks the event row FOR SHARE and returns its owner, or
// models.ErrEventNotFound / models.ErrEventAlreadyFinished (Finish locks the
// row FOR UPDATE, so members can't change while it finishes the event).
func (r *Repo) lockOpenEvent(ctx context.Context, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT owner_user_id, finished_at FROM events WHERE id = $1 FOR SHARE`, eventID)
	if err != nil {
		return 0, fmt.Errorf("lock event: %w", err)
	}

	type eventRow struct {
		OwnerUserID int        `db:"owner_user_id"`
		FinishedAt  *time.Time `db:"finished_at"`
	}
	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[eventRow])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return 0, models.ErrEventNotFound
	case err != nil:
		return 0, fmt.Errorf("lock event: %w", err)
	case event.FinishedAt != nil:
		return 0, models.ErrEventAlreadyFinished
	}

	return event.OwnerUserID, nil
}

// lockMembers takes pg_advisory_xact_lock(eventID, userID) for each user, in
// ascending order, the same keys and order as repo/game.lockPlayersInEvent.
func (r *Repo) lockMembers(ctx context.Context, eventID int, userIDs []int) error {
	ids := slices.Clone(userIDs)
	slices.Sort(ids)

	for _, userID := range slices.Compact(ids) {
		if _, err := r.db.Exec(ctx, `SELECT pg_advisory_xact_lock($1, $2)`, eventID, userID); err != nil {
			return fmt.Errorf("lock user %d in event %d: %w", userID, eventID, err)
		}
	}

	return nil
}

// checkNotInLiveGame returns errInGame, with the user's name for the client,
// if one of the users is in a live game of the event. Call it after
// lockMembers, so the answer holds until the transaction ends.
func (r *Repo) checkNotInLiveGame(ctx context.Context, eventID int, userIDs []int, errInGame error) error {
	rows, err := r.db.Query(ctx, `
		SELECT gp.player_id, u.name
		FROM games g
		JOIN game_players gp ON gp.game_id = g.id
		JOIN users u ON u.id = gp.player_id
		WHERE g.event_id = $1 AND g.status = $2 AND gp.player_id = ANY($3)
		ORDER BY u.name
		LIMIT 1`,
		eventID, models.GameLive, userIDs)
	if err != nil {
		return fmt.Errorf("check live games: %w", err)
	}

	type playing struct {
		PlayerID int    `db:"player_id"`
		Name     string `db:"name"`
	}
	player, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[playing])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("check live games: %w", err)
	default:
		return fmt.Errorf("user %d: %w", player.PlayerID,
			models.WithData(errInGame, models.ErrorData{"name": player.Name}))
	}
}
