package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// dataName is the key of a player's name in an error's data.
const dataName = "name"

// Sitters returns who sits at each seat of the game now: the seat's player
// (game_players) → the sitter, the player themself unless a substitute sits
// there.
func (r *Repo) Sitters(ctx context.Context, gameID int) (map[int]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT seat_player_id, sitter_id FROM game_sitters WHERE game_id = $1`, gameID)
	if err != nil {
		return nil, fmt.Errorf("sitters: %w", err)
	}
	type sitterRow struct {
		Seat   int `db:"seat_player_id"`
		Sitter int `db:"sitter_id"`
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[sitterRow])
	if err != nil {
		return nil, fmt.Errorf("sitters: %w", err)
	}
	sitters := make(map[int]int, len(list))
	for _, s := range list {
		sitters[s.Seat] = s.Sitter
	}
	return sitters, nil
}

// SeatOf returns the seat (its player's ID) the user sits at in the game
// now, or models.ErrPlayerNotFound when they don't sit there.
func (r *Repo) SeatOf(ctx context.Context, gameID, userID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT seat_player_id FROM game_sitters WHERE game_id = $1 AND sitter_id = $2`, gameID, userID)
	if err != nil {
		return 0, fmt.Errorf("seat of user %d: %w", userID, err)
	}
	seat, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return 0, models.ErrPlayerNotFound
	case err != nil:
		return 0, fmt.Errorf("seat of user %d: %w", userID, err)
	}
	return seat, nil
}

// HasPlace reports whether the user has a place in the game: a seat of
// their own, or one they sat at as a substitute. It is their game to see.
func (r *Repo) HasPlace(ctx context.Context, gameID, userID int) (bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT EXISTS (SELECT 1 FROM game_players WHERE game_id = $1 AND player_id = $2)
			OR EXISTS (SELECT 1 FROM game_seats WHERE game_id = $1 AND sitter_id = $2)`,
		gameID, userID)
	if err != nil {
		return false, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[bool])
}

// SeatSitter seats a user at a seat of a tournament game whose tour is
// published or live, from the next log entry on (seenLogSize: the log the
// referee saw; models.ErrLogChanged if it has changed). The sitter is a
// substitute of the event or the seat's own player (back at the table).
// With forfeit, the seat's player gets 0 TP and 0 GP for the game; seating
// the player again clears it. It takes the same locks as starting a game
// (the event FOR SHARE, the sitter's advisory lock, the game row) and
// refuses a sitter at another table that isn't over
// (models.ErrPlayerInAnotherGame) or at another seat of this one.
func (r *Repo) SeatSitter(ctx context.Context, seat models.SeatChange) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockOpenEvent(ctx, seat.EventID); err != nil {
			return err
		}
		sitter := seat.SitterID
		if err := tx.lockPlayersInEvent(ctx, [4]int{sitter, sitter, sitter, sitter}, seat.EventID); err != nil {
			return err
		}
		if err := tx.lockUnchangedLog(ctx, seat.GameID, seat.SeenLogSize); err != nil {
			return err
		}
		if err := tx.checkSeatable(ctx, seat); err != nil {
			return err
		}

		forfeit := seat.Forfeit && sitter != seat.SeatPlayerID
		if _, err := tx.db.Exec(ctx, `
			INSERT INTO game_seats (game_id, seat_player_id, sitter_id, from_log_size, created_at)
			VALUES ($1, $2, $3, $4, $5)`,
			seat.GameID, seat.SeatPlayerID, sitter, seat.SeenLogSize, time.Now().UTC()); err != nil {
			return fmt.Errorf("seat user %d: %w", sitter, err)
		}
		if _, err := tx.db.Exec(ctx, `
			UPDATE game_players SET forfeit = $3 WHERE game_id = $1 AND player_id = $2`,
			seat.GameID, seat.SeatPlayerID, forfeit); err != nil {
			return fmt.Errorf("forfeit of seat %d: %w", seat.SeatPlayerID, err)
		}
		return nil
	})
}

// checkSeatable checks a seat change under SeatSitter's locks: a game of
// the event in a published or live tour, not over; a seat of it; a sitter
// who is a substitute or the seat's player, sitting nowhere else at a table
// that isn't over (this game's other seats included).
func (r *Repo) checkSeatable(ctx context.Context, seat models.SeatChange) error {
	rows, err := r.db.Query(ctx, `
		SELECT g.status, t.status AS tour_status
		FROM games g JOIN tours t ON t.id = g.tour_id
		WHERE g.id = $1 AND g.event_id = $2`, seat.GameID, seat.EventID)
	if err != nil {
		return fmt.Errorf("game %d: %w", seat.GameID, err)
	}
	type gameRow struct {
		Status     models.GameStatus `db:"status"`
		TourStatus models.TourStatus `db:"tour_status"`
	}
	game, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[gameRow])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrGameNotFound
	case err != nil:
		return fmt.Errorf("game %d: %w", seat.GameID, err)
	case game.Status != models.GameCreated && game.Status != models.GameLive,
		game.TourStatus != models.TourPublished && game.TourStatus != models.TourLive:
		return fmt.Errorf("game %d is %s, tour %s: %w", seat.GameID, game.Status, game.TourStatus, models.ErrGameNotLive)
	}

	sitters, err := r.Sitters(ctx, seat.GameID)
	if err != nil {
		return err
	}
	current, ok := sitters[seat.SeatPlayerID]
	switch {
	case !ok:
		return fmt.Errorf("seat %d in game %d: %w", seat.SeatPlayerID, seat.GameID, models.ErrPlayerNotFound)
	case current == seat.SitterID:
		return fmt.Errorf("user %d already sits there: %w", seat.SitterID, models.ErrAlreadySeated)
	}

	if seat.SitterID != seat.SeatPlayerID {
		rows, err := r.db.Query(ctx, `
			SELECT EXISTS (SELECT 1 FROM event_user_roles
				WHERE event_id = $1 AND user_id = $2 AND role = $3)`,
			seat.EventID, seat.SitterID, models.EventSubstitute)
		if err != nil {
			return fmt.Errorf("substitute role: %w", err)
		}
		isSubstitute, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
		if err != nil {
			return fmt.Errorf("substitute role: %w", err)
		}
		if !isSubstitute {
			return fmt.Errorf("user %d: %w", seat.SitterID, models.ErrNotSubstitute)
		}
	}

	rows, err = r.db.Query(ctx, `
		SELECT u.name
		FROM game_sitters gs
		JOIN games g ON g.id = gs.game_id
		JOIN users u ON u.id = gs.sitter_id
		WHERE gs.sitter_id = $1 AND g.event_id = $2 AND g.status IN ($3, $4)
		LIMIT 1`,
		seat.SitterID, seat.EventID, models.GameCreated, models.GameLive)
	if err != nil {
		return fmt.Errorf("sitter's other tables: %w", err)
	}
	name, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("sitter's other tables: %w", err)
	}
	return fmt.Errorf("user %d sits at a table: %w", seat.SitterID,
		models.WithData(models.ErrPlayerInAnotherGame, models.ErrorData{dataName: name}))
}
