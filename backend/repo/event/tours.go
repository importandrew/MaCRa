package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// GenerateTour generates the next tour's seating as a draft, seen only by
// the tournament's staff: it creates the tour and its games (status
// created), one per table that makeTables returns, replacing an earlier
// draft. The event row is locked FOR UPDATE, so tours are made one at a
// time. It refuses a regular or finished event, a tour published or live
// (models.ErrTourInProgress), all tours played (models.ErrAllToursPlayed),
// numbers not drawn for the current players (models.ErrNotDrawn) and active
// participants that don't fill whole tables (models.ErrTablesIncomplete,
// with the number missing).
func (r *Repo) GenerateTour(ctx context.Context, eventID int, method models.SeatingMethod, makeTables models.TablesFunc) (models.Tour, error) {
	var tour models.Tour
	err := r.InTx(ctx, func(tx *Repo) error {
		isTournament, _, err := tx.tournamentState(ctx, eventID, true)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		}

		// A new draft replaces the old one (its games go with it).
		if _, err := tx.db.Exec(ctx, `
			DELETE FROM tours WHERE event_id = $1 AND status = $2`,
			eventID, models.TourDraft); err != nil {
			return fmt.Errorf("delete old draft: %w", err)
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

		input, err := tx.seatingInput(ctx, eventID, participants)
		if err != nil {
			return err
		}
		tables, err := makeTables(input)
		if err != nil {
			return err
		}

		tour, err = tx.insertDraftTour(ctx, eventID, next, method, tables)
		return err
	})
	return tour, err
}

// insertDraftTour inserts a drafted tour with its games, one per table.
func (r *Repo) insertDraftTour(ctx context.Context, eventID, number int, method models.SeatingMethod, tables [][4]int) (models.Tour, error) {
	rows, err := r.db.Query(ctx, `
		INSERT INTO tours (event_id, number, status, seating)
		VALUES ($1, $2, $3, $4)
		RETURNING id, number, status, seating, published_at, minutes, started_at, finished_at`,
		eventID, number, models.TourDraft, method)
	if err != nil {
		return models.Tour{}, fmt.Errorf("insert tour: %w", err)
	}
	tour, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tour])
	if err != nil {
		return models.Tour{}, fmt.Errorf("insert tour: %w", err)
	}

	now := time.Now().UTC()
	for i, seats := range tables {
		if err := r.insertTourGame(ctx, eventID, tour.ID, i+1, seats, now); err != nil {
			return models.Tour{}, err
		}
	}
	return tour, nil
}

// missingPlayers is how many players the active participants lack for whole
// tables of four (at least one table).
func missingPlayers(active int) int {
	if active == 0 {
		return models.SeatsPerTable
	}
	return (models.SeatsPerTable - active%models.SeatsPerTable) % models.SeatsPerTable
}

// nextTourNumber returns the number of the next tour, refusing while a tour
// is published or live (a draft doesn't count: it is to be replaced or
// published) and when all tours are played.
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

// PublishTour publishes the drafted seating: the players see their tables
// from now on. It refuses without a draft (models.ErrTourNotFound).
func (r *Repo) PublishTour(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		tourID, err := tx.lockTour(ctx, eventID, models.TourDraft)
		if err != nil {
			return err
		}
		_, err = tx.db.Exec(ctx, `
			UPDATE tours SET status = $1, published_at = $2 WHERE id = $3`,
			models.TourPublished, time.Now().UTC(), tourID)
		return err
	})
}

// CancelTour cancels the seating of a tour not started yet, drafted or
// published: the tour and its games are deleted, so it can be seated again.
// It refuses when there is none (models.ErrTourNotFound).
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
			DELETE FROM tours WHERE event_id = $1 AND status IN ($2, $3)`,
			eventID, models.TourDraft, models.TourPublished)
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
			SELECT id, number, status, seating, published_at, minutes, started_at, finished_at
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
			if state.Repeats, err = tx.tourRepeats(ctx, eventID, tour); err != nil {
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
		SELECT g.table_number, g.id AS game_id, g.status, g.extra_minutes, gp.player_id, u.name,
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
		Table        int               `db:"table_number"`
		GameID       int               `db:"game_id"`
		Status       models.GameStatus `db:"status"`
		ExtraMinutes int               `db:"extra_minutes"`
		UserID       int               `db:"player_id"`
		Name         string            `db:"name"`
		Number       int               `db:"number"`
		Wind         string            `db:"starting_wind"`
	}
	seats, err := pgx.CollectRows(rows, pgx.RowToStructByName[seatRow])
	if err != nil {
		return nil, err
	}

	tables := []models.TourTable{}
	for _, s := range seats {
		if len(tables) == 0 || tables[len(tables)-1].Table != s.Table {
			tables = append(tables, models.TourTable{Table: s.Table, GameID: s.GameID, Status: s.Status, ExtraMinutes: s.ExtraMinutes})
		}
		last := &tables[len(tables)-1]
		last.Seats = append(last.Seats, models.Seat{UserID: s.UserID, Name: s.Name, Number: s.Number, Wind: s.Wind})
	}
	return tables, nil
}

// StartTour starts the published tour: it and its games become live, and the
// tournament is in progress from its first tour (events.started_at). It
// refuses when no tour is published (models.ErrTourNotFound).
func (r *Repo) StartTour(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		tourID, err := tx.lockTour(ctx, eventID, models.TourPublished)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		// The tour keeps the length it starts with: a later change of the
		// setting applies to the next tours.
		if _, err := tx.db.Exec(ctx, `
			UPDATE tours SET status = $1, started_at = $2,
				minutes = (SELECT (tournament->>'tour_minutes')::int FROM events WHERE id = $4)
			WHERE id = $3`,
			models.TourLive, now, tourID, eventID); err != nil {
			return fmt.Errorf("start tour: %w", err)
		}
		if _, err := tx.db.Exec(ctx, `
			UPDATE games SET status = $1 WHERE tour_id = $2 AND status = $3`,
			models.GameLive, tourID, models.GameCreated); err != nil {
			return fmt.Errorf("start games: %w", err)
		}
		if _, err := tx.db.Exec(ctx, `
			UPDATE events SET started_at = COALESCE(started_at, $1) WHERE id = $2`,
			now, eventID); err != nil {
			return fmt.Errorf("start tournament: %w", err)
		}
		return nil
	})
}

// FinishTour finishes the live tour once all its games have ended: its
// results count from then on and can't be undone any more. It refuses when
// no tour is live (models.ErrTourNotFound) and while tables still play
// (models.ErrTablesStillPlaying, with their number). The tour row is locked
// FOR UPDATE: an undo in one of its games locks it FOR SHARE, so a game can't
// be revived while the tour finishes.
func (r *Repo) FinishTour(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		tourID, err := tx.lockTour(ctx, eventID, models.TourLive)
		if err != nil {
			return err
		}

		rows, err := tx.db.Query(ctx, `
			SELECT count(*) FROM games WHERE tour_id = $1 AND status <> $2`,
			tourID, models.GameEnded)
		if err != nil {
			return err
		}
		playing, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
		if err != nil {
			return err
		}
		if playing > 0 {
			return models.WithData(models.ErrTablesStillPlaying, models.ErrorData{"playing": playing})
		}

		_, err = tx.db.Exec(ctx, `
			UPDATE tours SET status = $1, finished_at = $2 WHERE id = $3`,
			models.TourFinished, time.Now().UTC(), tourID)
		return err
	})
}

// AbortTour aborts the live tour: all its games become aborted, the tour
// doesn't count and its number is played again. It refuses when no tour is
// live (models.ErrTourNotFound).
func (r *Repo) AbortTour(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		tourID, err := tx.lockTour(ctx, eventID, models.TourLive)
		if err != nil {
			return err
		}

		if _, err := tx.db.Exec(ctx, `
			UPDATE games SET status = $1 WHERE tour_id = $2`,
			models.GameAborted, tourID); err != nil {
			return fmt.Errorf("abort games: %w", err)
		}
		_, err = tx.db.Exec(ctx, `
			UPDATE tours SET status = $1, finished_at = $2 WHERE id = $3`,
			models.TourAborted, time.Now().UTC(), tourID)
		return err
	})
}

// lockTour locks the tournament's event row and then its tour in the given
// status FOR UPDATE, and returns the tour's ID, or models.ErrTourNotFound.
func (r *Repo) lockTour(ctx context.Context, eventID int, status models.TourStatus) (int, error) {
	isTournament, _, err := r.tournamentState(ctx, eventID, true)
	switch {
	case err != nil:
		return 0, err
	case !isTournament:
		return 0, models.ErrNotTournament
	}

	rows, err := r.db.Query(ctx, `
		SELECT id FROM tours WHERE event_id = $1 AND status = $2 FOR UPDATE`,
		eventID, status)
	if err != nil {
		return 0, fmt.Errorf("lock tour: %w", err)
	}
	tourID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, models.ErrTourNotFound
	}
	return tourID, err
}

// AddTime adds minutes to a table of the live tour: the game's time is up
// that much later. It refuses a game that isn't a live game of the live tour
// (models.ErrTourNotFound).
func (r *Repo) AddTime(ctx context.Context, eventID, gameID, minutes int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		tourID, err := tx.lockTour(ctx, eventID, models.TourLive)
		if err != nil {
			return err
		}
		res, err := tx.db.Exec(ctx, `
			UPDATE games SET extra_minutes = extra_minutes + $1
			WHERE id = $2 AND tour_id = $3 AND status = $4`,
			minutes, gameID, tourID, models.GameLive)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return fmt.Errorf("game %d: %w", gameID, models.ErrTourNotFound)
		}
		return nil
	})
}

// seatingInput gathers what a seating is made from: the participants, how
// many times each pair has met in the tournament's live and finished tours
// (aborted ones don't count), and the totals from the finished tours.
func (r *Repo) seatingInput(ctx context.Context, eventID int, participants []models.Participant) (models.SeatingInput, error) {
	input := models.SeatingInput{
		Participants: participants,
		Meetings:     map[models.SeatingPair]int{},
		Standings:    map[int]models.SeatingStanding{},
	}

	rows, err := r.db.Query(ctx, `
		SELECT a.player_id AS a, b.player_id AS b, count(*) AS meetings
		FROM tours t
		JOIN games g ON g.tour_id = t.id
		JOIN game_players a ON a.game_id = g.id
		JOIN game_players b ON b.game_id = g.id AND a.player_id < b.player_id
		WHERE t.event_id = $1 AND t.status IN ($2, $3)
		GROUP BY a.player_id, b.player_id`,
		eventID, models.TourLive, models.TourFinished)
	if err != nil {
		return input, fmt.Errorf("meetings: %w", err)
	}
	type meetingRow struct {
		A        int `db:"a"`
		B        int `db:"b"`
		Meetings int `db:"meetings"`
	}
	meetings, err := pgx.CollectRows(rows, pgx.RowToStructByName[meetingRow])
	if err != nil {
		return input, fmt.Errorf("meetings: %w", err)
	}
	for _, m := range meetings {
		input.Meetings[models.SeatingPair{m.A, m.B}] = m.Meetings
	}

	rows, err = r.db.Query(ctx, `
		SELECT pr.player_id, sum(pr.tournament_points)::float8 AS tp, sum(pr.game_points)::int AS gp
		FROM tours t
		JOIN games g ON g.tour_id = t.id AND g.status = $3
		JOIN player_results pr ON pr.game_id = g.id
		WHERE t.event_id = $1 AND t.status = $2
		GROUP BY pr.player_id`,
		eventID, models.TourFinished, models.GameEnded)
	if err != nil {
		return input, fmt.Errorf("standings: %w", err)
	}
	type standingRow struct {
		PlayerID int     `db:"player_id"`
		TP       float64 `db:"tp"`
		GP       int     `db:"gp"`
	}
	standings, err := pgx.CollectRows(rows, pgx.RowToStructByName[standingRow])
	if err != nil {
		return input, fmt.Errorf("standings: %w", err)
	}
	for _, s := range standings {
		input.Standings[s.PlayerID] = models.SeatingStanding{TP: s.TP, GP: s.GP}
	}
	return input, nil
}

// tourRepeats counts the repeat meetings at a tour's tables: for each pair
// seated together, how many times they met in the tournament's earlier live
// or finished tours.
func (r *Repo) tourRepeats(ctx context.Context, eventID int, tour models.Tour) (int, error) {
	rows, err := r.db.Query(ctx, `
		WITH cur AS (
			SELECT g.id AS game_id, gp.player_id FROM games g
			JOIN game_players gp ON gp.game_id = g.id
			WHERE g.tour_id = $2
		), prev AS (
			SELECT g.id AS game_id, gp.player_id FROM tours t
			JOIN games g ON g.tour_id = t.id
			JOIN game_players gp ON gp.game_id = g.id
			WHERE t.event_id = $1 AND t.id <> $2 AND t.number < $3 AND t.status IN ($4, $5)
		)
		SELECT count(*) FROM cur a
		JOIN cur b ON b.game_id = a.game_id AND a.player_id < b.player_id
		JOIN prev pa ON pa.player_id = a.player_id
		JOIN prev pb ON pb.game_id = pa.game_id AND pb.player_id = b.player_id`,
		eventID, tour.ID, tour.Number, models.TourLive, models.TourFinished)
	if err != nil {
		return 0, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[int])
}
