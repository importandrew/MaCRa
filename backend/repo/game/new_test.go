package game

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestNewGame(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		seating       [4]int
		preAddPlayers []int

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers map[string]int // wind -> player_id
	}{
		"ok": {
			seating:       [4]int{1, 2, 3, 4},
			preAddPlayers: []int{1, 2, 3, 4},
			expectedErr:   assert.NoError,
			expectedPlayers: map[string]int{
				"east":  1,
				"south": 2,
				"west":  3,
				"north": 4,
			},
		},
		"player not exists": {
			seating:       [4]int{1, 2, 3, 4},
			preAddPlayers: []int{1, 2, 3, 5},
			expectedErr:   assert.Error,
		},
		"duplicate player": {
			seating:       [4]int{1, 2, 3, 3},
			preAddPlayers: []int{1, 2, 3},
			expectedErr:   assert.Error,
		},
		"short seating": {
			seating:       [4]int{1, 2, 3},
			preAddPlayers: []int{1, 2, 3},
			expectedErr:   assert.Error,
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

			for _, id := range testCase.preAddPlayers {
				createTestUser(t, conn, id)
			}
			createTestEvent(t, conn)

			gameID, err := repo.NewGame(ctx, testCase.seating, testEventID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Positive(t, gameID, "Game ID should be greater than 0")

			status := getGameStatus(t, conn, gameID)
			assert.Equal(t, models.GameLive, status, "Game status should be %s", models.GameLive)

			players := getGamePlayers(t, conn, gameID)
			require.Len(t, players, 4, "Game should have exactly 4 players")
			for wind, id := range testCase.expectedPlayers {
				assert.Equal(t, players[wind], id)
			}
		})
	}
}

func TestNewGameAlreadyInGame(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.Test(t, "player already in live game", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for id := 1; id <= 7; id++ {
			createTestUser(t, conn, id)
		}
		createTestEvent(t, conn)

		_, err = repo.NewGame(ctx, [4]int{1, 2, 3, 4}, testEventID)
		require.NoError(t, err)

		_, err = repo.NewGame(ctx, [4]int{5, 6, 7, 4}, testEventID)
		assert.ErrorIs(t, err, models.ErrValidatePlayers)
	})

	db.Test(t, "concurrent creations", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for id := 1; id <= 7; id++ {
			createTestUser(t, conn, id)
		}
		createTestEvent(t, conn)

		// The same table twice (a double tap) plus a table sharing player 4,
		// seated in a different order to exercise lock ordering.
		seatings := [][4]int{
			{1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4},
			{4, 7, 6, 5}, {5, 6, 7, 4},
		}

		const rounds = 20
		for range rounds {
			_, err := conn.Exec(ctx, `DELETE FROM games`)
			require.NoError(t, err)

			errs := make(chan error, len(seatings))
			var wg sync.WaitGroup
			for _, seating := range seatings {
				wg.Go(func() {
					_, err := repo.NewGame(ctx, seating, testEventID)
					errs <- err
				})
			}
			wg.Wait()
			close(errs)

			created := 0
			for err := range errs {
				if err == nil {
					created++
					continue
				}
				require.ErrorIs(t, err, models.ErrValidatePlayers)
			}
			require.Equal(t, 1, created, "exactly one game must be created")

			var live int
			err = conn.QueryRow(ctx, `SELECT count(*) FROM games WHERE status = $1`, models.GameLive).Scan(&live)
			require.NoError(t, err)
			require.Equal(t, 1, live)
		}
	})
}
