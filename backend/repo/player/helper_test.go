package player

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func createTestUserWithName(t *testing.T, conn *pgxpool.Conn, id int, name string) int {
	t.Helper()
	ctx := context.Background()
	var userID int
	err := conn.QueryRow(ctx, `
		INSERT INTO users(id, name, email, password)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			password = EXCLUDED.password
		RETURNING id`,
		id, name, fmt.Sprintf("%d@test.com", id), "password123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return userID
}

func createTestGame(t *testing.T, conn *pgxpool.Conn, gameID int, status models.GameStatus) int {
	t.Helper()
	createTestEvent(t, conn)
	ctx := context.Background()
	var id int
	err := conn.QueryRow(ctx, `
		INSERT INTO games(id, created_at, status, event_id)
		VALUES($1, CURRENT_TIMESTAMP, $2, $3)
		RETURNING id`, gameID, status, testEventID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return id
}

func createTestGamePlayer(t *testing.T, conn *pgxpool.Conn, gameID, playerID int, wind string) {
	t.Helper()
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		INSERT INTO game_players(game_id, player_id, starting_wind)
		VALUES($1, $2, $3)`,
		gameID, playerID, wind,
	)
	if err != nil {
		t.Fatalf("Failed to create test game player: %v", err)
	}
}
