package game

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// NewGame creates a live game in the event with the given seating
// (east, south, west, north). It returns models.ErrPlayerInAnotherGame if
// any of the players already has a live game in the event.
//
// Two concurrent requests could both pass a check made before the insert.
// NewGame therefore takes a transaction-level advisory lock per (event,
// player) and checks inside the transaction, so a concurrent creation for the
// same player waits until the first one commits and then sees its game.
func (r *Repo) NewGame(ctx context.Context, seating [4]int, eventID int) (int, error) {
	var gameID int
	err := r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockOpenEvent(ctx, eventID); err != nil {
			return err
		}

		if err := tx.lockPlayersInEvent(ctx, seating, eventID); err != nil {
			return err
		}

		if err := tx.checkNoOtherLiveGame(ctx, seating, eventID, 0); err != nil {
			return err
		}

		// The game keeps a copy of the event's rules: it is replayed by them.
		rows, err := tx.db.Query(ctx, `
			INSERT INTO games(created_at, status, event_id, rules)
			SELECT $1, $2, id, rules FROM events WHERE id = $3
			RETURNING id`,
			time.Now().UTC(), models.GameLive, eventID)
		if err != nil {
			return err
		}
		gameID, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
		if err != nil {
			return err
		}

		_, err = tx.db.Exec(ctx, `
			INSERT INTO game_players(game_id, player_id, starting_wind) 
			VALUES 
				($1, $2, 'east'),
				($1, $3, 'south'),
				($1, $4, 'west'),
				($1, $5, 'north')`,
			gameID, seating[0], seating[1], seating[2], seating[3])
		return err
	})
	if err != nil {
		return 0, err
	}

	return gameID, nil
}

// lockOpenEvent locks the event row FOR SHARE until the transaction ends and returns
// models.ErrEventAlreadyFinished if the event is finished.
//
// The event repo's Finish locks the same row FOR UPDATE before it checks for
// live games, so a game can't become live in an event that is being finished:
// whichever transaction comes second waits for the first and sees its result.
//
// Call it on a Repo from InTx: outside a transaction the lock is released at
// once.
func (r *Repo) lockOpenEvent(ctx context.Context, eventID int) error {
	rows, err := r.db.Query(ctx, `SELECT finished_at FROM events WHERE id = $1 FOR SHARE`, eventID)
	if err != nil {
		return fmt.Errorf("lock event: %w", err)
	}

	finishedAt, err := pgx.CollectOneRow(rows, pgx.RowTo[*time.Time])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrEventNotFound
	case err != nil:
		return fmt.Errorf("lock event: %w", err)
	case finishedAt != nil:
		return models.ErrEventAlreadyFinished
	}

	return nil
}

// lockPlayersInEvent takes pg_advisory_xact_lock(eventID, playerID) for each
// player, in ascending player order so that overlapping seatings cannot
// deadlock. The locks are released when the transaction ends; call it on a
// Repo from InTx.
func (r *Repo) lockPlayersInEvent(ctx context.Context, seating [4]int, eventID int) error {
	playerIDs := slices.Clone(seating[:])
	slices.Sort(playerIDs)

	for _, playerID := range playerIDs {
		if _, err := r.db.Exec(ctx, `SELECT pg_advisory_xact_lock($1, $2)`, eventID, playerID); err != nil {
			return fmt.Errorf("lock player %d in event %d: %w", playerID, eventID, err)
		}
	}

	return nil
}

// checkNoOtherLiveGame returns models.ErrPlayerInAnotherGame, with the
// player's name for the client, if a player of the seating has a live game in
// the event other than exceptGameID (0 for none). Call it on a Repo from InTx
// after lockPlayersInEvent, so the answer holds until the transaction ends.
func (r *Repo) checkNoOtherLiveGame(ctx context.Context, seating [4]int, eventID, exceptGameID int) error {
	type liveGamePlayer struct {
		PlayerID int    `db:"player_id"`
		Name     string `db:"name"`
		GameID   int    `db:"game_id"`
	}

	rows, err := r.db.Query(ctx, `
		SELECT gp.player_id, u.name, g.id AS game_id
		FROM games g
		JOIN game_players gp ON gp.game_id = g.id
		JOIN users u ON u.id = gp.player_id
		WHERE gp.player_id = ANY($1) AND g.event_id = $2 AND g.status = $3 AND g.id <> $4
		ORDER BY u.name
		LIMIT 1`,
		seating[:], eventID, models.GameLive, exceptGameID)
	if err != nil {
		return fmt.Errorf("check live games: %w", err)
	}

	live, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[liveGamePlayer])
	switch {
	case err == nil:
		return fmt.Errorf("player %d in game %d: %w", live.PlayerID, live.GameID,
			models.WithData(models.ErrPlayerInAnotherGame, models.ErrorData{"name": live.Name}))
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	default:
		return fmt.Errorf("check live games: %w", err)
	}
}
