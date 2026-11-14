package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// PublishTour publishes the next tour of a tournament: it creates the tour
// and its games (status created), one per table that makeTables returns.
// The event row is locked FOR UPDATE, so tours are published one at a time.
// It refuses a regular or finished event, a tour still published or live
// (models.ErrTourInProgress), all tours played (models.ErrAllToursPlayed),
// numbers not drawn for the current players (models.ErrNotDrawn) and active
// participants that don't fill whole tables (models.ErrTablesIncomplete,
// with the number missing).
func (r *Repo) PublishTour(ctx context.Context, eventID int, method models.SeatingMethod, makeTables models.TablesFunc) (models.Tour, error) {
	var tour models.Tour
	err := r.InTx(ctx, func(tx *Repo) error {
		isTournament, _, err := tx.tournamentState(ctx, eventID, true)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		}

		next, err := tx.nextTourNumber(ctx, eventID)
		if err != nil {
			return err
		}

		participants, err := tx.activeParticipants(ctx, eventID)
		if err != nil {
			return err
		}
		if missing := missingPlayers(len(participants)); missing > 0 {
			return models.WithData(models.ErrTablesIncomplete, models.ErrorData{"missing": missing})
		}

		tables, err := makeTables(participants)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		rows, err := tx.db.Query(ctx, `
			INSERT INTO tours (event_id, number, status, seating, published_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, number, status, seating, published_at, started_at, finished_at`,
			eventID, next, models.TourPublished, method, now)
		if err != nil {
			return fmt.Errorf("insert tour: %w", err)
		}
		if tour, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tour]); err != nil {
			return fmt.Errorf("insert tour: %w", err)
		}

		for i, seats := range tables {
			if err := tx.insertTourGame(ctx, eventID, tour.ID, i+1, seats, now); err != nil {
				return err
			}
		}
		return nil
	})
	return tour, err
}

// missingPlayers is how many players the active participants lack for whole
// tables of four (at least one table).
func missingPlayers(active int) int {
	if active == 0 {
		return models.SeatsPerTable
	}
	return (models.SeatsPerTable - active%models.SeatsPerTable) % models.SeatsPerTable
}

// nextTourNumber returns the number of the tour to publish, refusing while a
// tour is published or live and when all tours are played.
func (r *Repo) nextTourNumber(ctx context.Context, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			count(*) FILTER (WHERE t.status IN ($2, $3)) AS in_progress,
			count(*) FILTER (WHERE t.status = $4) AS finished,
			(SELECT (e.tournament->>'tours')::int FROM events e WHERE e.id = $1) AS tours
		FROM tours t WHERE t.event_id = $1`,
		eventID, models.TourPublished, models.TourLive, models.TourFinished)
	if err != nil {
		return 0, fmt.Errorf("count tours: %w", err)
	}

	type counts struct {
		InProgress int `db:"in_progress"`
		Finished   int `db:"finished"`
		Tours      int `db:"tours"`
	}
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[counts])
	switch {
	case err != nil:
		return 0, fmt.Errorf("count tours: %w", err)
	case c.InProgress > 0:
		return 0, models.ErrTourInProgress
	case c.Finished >= c.Tours:
		return 0, models.ErrAllToursPlayed
	}
	return c.Finished + 1, nil
}

// activeParticipants returns the participants to seat: the numbered players
// who haven't withdrawn, or models.ErrNotDrawn if the numbers aren't drawn
// for the current players.
func (r *Repo) activeParticipants(ctx context.Context, eventID int) ([]models.Participant, error) {
	list, err := r.ListParticipants(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if !list.Drawn {
		return nil, models.ErrNotDrawn
	}

	active := make([]models.Participant, 0, len(list.Participants))
	for _, p := range list.Participants {
		if !p.Withdrawn {
			active = append(active, models.Participant{UserID: p.UserID, Number: *p.Number})
		}
	}
	return active, nil
}

// insertTourGame creates a tour's game at a table (status created) with the
// event's rules, seats east to north.
func (r *Repo) insertTourGame(ctx context.Context, eventID, tourID, table int, seats [4]int, now time.Time) error {
	rows, err := r.db.Query(ctx, `
		INSERT INTO games (created_at, status, event_id, rules, tour_id, table_number)
		SELECT $1, $2, id, rules, $3, $4 FROM events WHERE id = $5
		RETURNING id`,
		now, models.GameCreated, tourID, table, eventID)
	if err != nil {
		return fmt.Errorf("insert game of table %d: %w", table, err)
	}
	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("insert game of table %d: %w", table, err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO game_players (game_id, player_id, starting_wind)
		VALUES ($1, $2, 'east'), ($1, $3, 'south'), ($1, $4, 'west'), ($1, $5, 'north')`,
		gameID, seats[0], seats[1], seats[2], seats[3])
	if err != nil {
		return fmt.Errorf("seat table %d: %w", table, err)
	}
	return nil
}

// CancelTour cancels the published seating of a tour not started yet: the
// tour and its games are deleted, so it can be seated again. It refuses when
// no tour is published (models.ErrTourNotFound).
func (r *Repo) CancelTour(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		isTournament, _, err := tx.tournamentState(ctx, eventID, true)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		}

		// The games go with the tour (ON DELETE CASCADE), their players with
		// them.
		res, err := tx.db.Exec(ctx, `
			DELETE FROM tours WHERE event_id = $1 AND status = $2`,
			eventID, models.TourPublished)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return models.ErrTourNotFound
		}
		return nil
	})
}

// TourState returns the tournament's current tour (the last one not
// aborted) with its tables, and what the next tour would be.
func (r *Repo) TourState(ctx context.Context, eventID int) (models.TourState, error) {
	state := models.TourState{Tables: []models.TourTable{}}
	err := r.InTx(ctx, func(tx *Repo) error {
		rows, err := tx.db.Query(ctx, `
			SELECT id, number, status, seating, published_at, started_at, finished_at
			FROM tours WHERE event_id = $1 AND status <> $2
			ORDER BY number DESC LIMIT 1`,
			eventID, models.TourAborted)
		if err != nil {
			return err
		}
		tour, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tour])
		switch {
		case errors.Is(err, pgx.ErrNoRows):
		case err != nil:
			return err
		default:
			state.Current = &tour
			if state.Tables, err = tx.tourTables(ctx, tour.ID); err != nil {
				return err
			}
		}

		next, err := tx.nextTourNumber(ctx, eventID)
		switch {
		case errors.Is(err, models.ErrTourInProgress), errors.Is(err, models.ErrAllToursPlayed):
		case err != nil:
			return err
		default:
			state.Next = next
		}

		list, err := tx.ListParticipants(ctx, eventID)
		if err != nil {
			return err
		}
		for _, p := range list.Participants {
			if !p.Withdrawn {
				state.Active++
			}
		}
		state.Missing = missingPlayers(state.Active)
		return nil
	})
	return state, err
}

// tourTables returns a tour's tables, each with its seats east to north.
func (r *Repo) tourTables(ctx context.Context, tourID int) ([]models.TourTable, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.table_number, g.id AS game_id, gp.player_id, u.name,
			COALESCE(tp.number, 0) AS number, gp.starting_wind
		FROM games g
		JOIN game_players gp ON gp.game_id = g.id
		JOIN users u ON u.id = gp.player_id
		LEFT JOIN tournament_participants tp ON tp.event_id = g.event_id AND tp.user_id = gp.player_id
		WHERE g.tour_id = $1
		ORDER BY g.table_number,
			CASE gp.starting_wind WHEN 'east' THEN 1 WHEN 'south' THEN 2 WHEN 'west' THEN 3 ELSE 4 END`,
		tourID)
	if err != nil {
		return nil, err
	}

	type seatRow struct {
		Table  int    `db:"table_number"`
		GameID int    `db:"game_id"`
		UserID int    `db:"player_id"`
		Name   string `db:"name"`
		Number int    `db:"number"`
		Wind   string `db:"starting_wind"`
	}
	seats, err := pgx.CollectRows(rows, pgx.RowToStructByName[seatRow])
	if err != nil {
		return nil, err
	}

	tables := []models.TourTable{}
	for _, s := range seats {
		if len(tables) == 0 || tables[len(tables)-1].Table != s.Table {
			tables = append(tables, models.TourTable{Table: s.Table, GameID: s.GameID})
		}
		last := &tables[len(tables)-1]
		last.Seats = append(last.Seats, models.Seat{UserID: s.UserID, Name: s.Name, Number: s.Number, Wind: s.Wind})
	}
	return tables, nil
}
