package player

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createTestUser(t *testing.T, conn *pgxpool.Conn, id int) {
	t.Helper()
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		INSERT INTO users(id, name, email, password)
		VALUES($1, $2, $3, $4)
		RETURNING id`,
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
	ctx := context.Background()
	var id int
	err := conn.QueryRow(ctx, `
		INSERT INTO games(id, created_at, status)
		VALUES($1, CURRENT_TIMESTAMP, $2)
		RETURNING id`, gameID, status,
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
