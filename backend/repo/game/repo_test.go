package game

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

// TestInTx checks that methods called on the Repo that InTx passes run in one
// transaction: they see each other's writes, and an error rolls all of them
// back, including those of a method with its own InTx (a savepoint).
func TestInTx(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	errStop := errors.New("stop")

	testData := map[string]struct {
		fnErr     error
		wantGames int
	}{
		"commit":   {fnErr: nil, wantGames: 1},
		"rollback": {fnErr: errStop, wantGames: 0},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)

			ctx := context.Background()
			conn, err := db.Acquire(ctx)
			require.NoError(t, err)
			defer conn.Release()

			for id := 1; id <= 4; id++ {
				createTestUser(t, conn, id)
			}
			createTestEvent(t, conn)

			err = repo.InTx(ctx, func(tx *Repo) error {
				gameID, err := tx.NewGame(ctx, [4]int{1, 2, 3, 4}, testEventID)
				require.NoError(t, err)

				// The uncommitted game is visible inside the transaction.
				status, err := tx.Status(ctx, gameID)
				require.NoError(t, err)
				assert.Equal(t, models.GameLive, status)

				return testCase.fnErr
			})
			assert.ErrorIs(t, err, testCase.fnErr)

			rows, err := conn.Query(ctx, `SELECT count(*) FROM games`)
			require.NoError(t, err)
			games, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
			require.NoError(t, err)
			assert.Equal(t, testCase.wantGames, games)
		})
	}
}
