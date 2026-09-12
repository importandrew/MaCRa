package game

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getGamePlayers(t *testing.T, conn *pgxpool.Conn, gameID int) map[string]int {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT player_id, starting_wind
		FROM game_players
		WHERE game_id = $1`,
		gameID,
	)
	if err != nil {
		t.Fatalf("Failed to query game players: %v", err)
	}
	defer rows.Close()

	players := make(map[string]int)
	for rows.Next() {
		var playerID int
		var wind string
		if err := rows.Scan(&playerID, &wind); err != nil {
			t.Fatalf("Failed to scan game player: %v", err)
		}
		players[wind] = playerID
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Error iterating game players: %v", err)
	}
	return players
}

const testEventID = 1

func createTestEvent(t *testing.T, conn *pgxpool.Conn) {
	t.Helper()
	ctx := context.Background()
	createTestUser(t, conn, 1)
	_, err := conn.Exec(ctx, `
		INSERT INTO events(id, type, name, owner_user_id, visibility, ratings_reach, games_reach)
		VALUES($1, 'regular', $2, 1, 'public', 'public', 'public')
		ON CONFLICT (id) DO NOTHING`,
		testEventID, "Test Event",
	)
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
}

func createTestUser(t *testing.T, conn *pgxpool.Conn, id int) {
	t.Helper()
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		INSERT INTO users(id, name, email, password)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			password = EXCLUDED.password`,
		id, strconv.Itoa(id), fmt.Sprintf("%d@test.com", id), "password123",
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}

func createTestGame(t *testing.T, conn *pgxpool.Conn, gameID int) {
	t.Helper()
	createTestEvent(t, conn)
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		INSERT INTO games(id, created_at, status, event_id)
		VALUES($1, CURRENT_TIMESTAMP, 'live', $2)
		RETURNING id`, gameID, testEventID,
	)
	if err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
}

func createTestGamePlayers(t *testing.T, conn *pgxpool.Conn, gameID int, userIDs [4]int) {
	ctx := context.Background()
	winds := [4]string{"east", "south", "west", "north"}
	for i, id := range userIDs {
		// dont create user
		if id == 0 {
			continue
		}

		_, err := conn.Exec(ctx, `
		INSERT INTO game_players(game_id, player_id, starting_wind)
		VALUES ($1, $2, $3)`,
			gameID, id, winds[i],
		)
		if err != nil {
			t.Fatalf("Failed to create game players: %v", err)
		}
	}
}

func createTestGameWithStatus(t *testing.T, conn *pgxpool.Conn, gameID int, status models.GameStatus) int {
	t.Helper()
	createTestEvent(t, conn)
	ctx := context.Background()
	var id int
	err := conn.QueryRow(ctx, `
		INSERT INTO games(id, created_at, status, event_id)
		VALUES($1, CURRENT_TIMESTAMP, $2, $3)
		RETURNING id`, gameID, string(status), testEventID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}

	return id
}

func getGameStatus(t *testing.T, conn *pgxpool.Conn, gameID int) models.GameStatus {
	t.Helper()
	ctx := context.Background()
	var status models.GameStatus
	err := conn.QueryRow(ctx, `
		SELECT status FROM games WHERE id = $1`,
		gameID,
	).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to get game status: %v", err)
	}
	return status
}

type roundRecord struct {
	ID             int           `db:"id"`
	RoundNumber    int           `db:"round_number"`
	Outcome        string        `db:"outcome"`
	WinnerID       sql.NullInt64 `db:"winner_id"`
	LoserID        sql.NullInt64 `db:"loser_id"`
	Value          sql.NullInt64 `db:"value"`
	ForceNextRound bool          `db:"force_next_round"`
}

func getRound(t *testing.T, conn *pgxpool.Conn, gameID, roundNumber int) *roundRecord {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT id, round_number, outcome, winner_id, loser_id, value, force_next_round
		FROM rounds
		WHERE game_id = $1 AND round_number = $2
		ORDER BY event_sequence DESC
		LIMIT 1`,
		gameID, roundNumber,
	)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roundRecord])
	if err != nil {
		t.Fatalf("Failed to collect rows: %v", err)
	}

	return &rec
}

func getRoundScores(t *testing.T, conn *pgxpool.Conn, roundID int) map[string]int {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT gp.starting_wind, rs.score
		FROM round_scores rs
		JOIN game_players gp ON rs.player_id = gp.player_id AND rs.game_id = gp.game_id
		WHERE rs.round_id = $1`,
		roundID,
	)
	if err != nil {
		t.Fatalf("Failed to query round scores: %v", err)
	}
	defer rows.Close()

	scores := make(map[string]int)
	for rows.Next() {
		var wind string
		var score int
		if err := rows.Scan(&wind, &score); err != nil {
			t.Fatalf("Failed to scan round score: %v", err)
		}
		scores[wind] = score
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Error iterating round scores: %v", err)
	}
	return scores
}

type preAddRound struct {
	roundNumber int
	outcome     models.Outcome
	winnerID    int
	loserID     int
	value       int
	scores      map[string]int // wind -> score
}

func createTestRound(t *testing.T, conn *pgxpool.Conn, gameID int, round preAddRound) {
	t.Helper()
	ctx := context.Background()

	// Insert round - handle NULL for winner_id, loser_id, and value when they're 0
	var winnerID, loserID, value any
	if round.winnerID != 0 {
		winnerID = round.winnerID
	}
	if round.loserID != 0 {
		loserID = round.loserID
	}
	if round.value != 0 {
		value = round.value
	}

	rows, err := conn.Query(ctx, `
		INSERT INTO rounds(
			game_id, round_number, event_sequence, outcome, winner_id, loser_id, 
			value, force_next_round, created_at)
		VALUES($1, $2, NULL, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		RETURNING id, event_sequence`,
		gameID, round.roundNumber, string(round.outcome), winnerID, loserID, value, false)
	if err != nil {
		t.Fatalf("Failed to create test round: %v", err)
	}

	type roundResult struct {
		ID            int `db:"id"`
		EventSequence int `db:"event_sequence"`
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roundResult])
	if err != nil {
		t.Fatalf("Failed to collect round result: %v", err)
	}

	// Get players for this game
	players := getGamePlayers(t, conn, gameID)

	// Insert round scores
	winds := []string{"east", "south", "west", "north"}
	playerIDs := []int{players["east"], players["south"], players["west"], players["north"]}

	for i, wind := range winds {
		score := round.scores[wind]
		_, err = conn.Exec(ctx, `
			INSERT INTO round_scores(round_id, game_id, round_number, event_sequence, player_id, score)
			VALUES($1, $2, $3, $4, $5, $6)
			ON CONFLICT (round_id, player_id) DO UPDATE SET score = EXCLUDED.score`,
			result.ID, gameID, round.roundNumber, result.EventSequence, playerIDs[i], score)
		if err != nil {
			t.Fatalf("Failed to create round score for %s: %v", wind, err)
		}
	}
}

type playerResult struct {
	ID int     `db:"player_id"`
	GP int     `db:"game_points"`
	TP float32 `db:"tournament_points"`
}

func getPlayerResults(t *testing.T, conn *pgxpool.Conn, gameID int) []playerResult {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT player_id, game_points, tournament_points
		FROM player_results
		WHERE game_id = $1`,
		gameID)
	if err != nil {
		t.Fatalf("Failed to query player results: %v", err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[playerResult])
	if err != nil {
		t.Fatalf("Failed to collect player results: %v", err)
	}

	return results
}

func setGameFinishedAt(t *testing.T, conn *pgxpool.Conn, gameID int) {
	t.Helper()
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		UPDATE games 
		SET finished_at = CURRENT_TIMESTAMP 
		WHERE id = $1`,
		gameID)
	if err != nil {
		t.Fatalf("Failed to set finished_at: %v", err)
	}
}
