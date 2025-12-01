package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/repo/database"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		email       string
		name        string
		hash        string

		expectedErr assert.ErrorAssertionFunc
	}{
		"ok": {
			email:       "alice@test.com",
			name:        "Alice",
			hash:        "hashedpassword123",
			expectedErr: assert.NoError,
		},
		"already exists": {
			preAddUsers: []int{1},
			email:       "1@test.com",
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)

			ctx := context.Background()
			conn, err := db.Acquire(ctx)
			if err != nil {
				t.Fatalf("Failed to acquire connection: %v", err)
			}
			defer conn.Release()

			for _, id := range testCase.preAddUsers {
				createTestUser(t, conn, id)
			}

			playerID, err := repo.Add(ctx, testCase.email, testCase.name, testCase.hash)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Positive(t, playerID, "Player ID should be greater than 0")

			// Verify the player was actually created
			rows, err := conn.Query(ctx, "SELECT name, email, password FROM users WHERE id = $1", playerID)
			require.NoError(t, err, "Should be able to query created user")
			if err != nil {
				return
			}

			type playerWithPass struct {
				Name  string `db:"name"`
				Email string `db:"email"`
				Pass  string `db:"password"`
			}
			player, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[playerWithPass])
			require.NoError(t, err, "Should be able to collect rows")
			if err != nil {
				return
			}

			assert.Equal(t, testCase.name, player.Name)
			assert.Equal(t, testCase.email, player.Email)
			assert.Equal(t, testCase.hash, player.Pass)
		})
	}
}
