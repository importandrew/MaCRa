package player

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
