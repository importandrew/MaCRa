package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLastGame(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddPlayers []int
		preAddGames   []struct {
			gameID    int
			playerIDs [4]int
			status    models.GameStatus
			finished  bool // whether to set finished_at
		}

		playerID     int
		expectedErr  assert.ErrorAssertionFunc
		expectedGame int
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				playerIDs [4]int
				status    models.GameStatus
				finished  bool
			}{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameEnded,
					finished:  true,
				},
			},
			playerID:    1,
			expectedErr: assert.NoError,
			expectedGame: 1,
		},
		"ok multiple games": {
			preAddPlayers: []int{1, 2, 3, 4, 5, 6, 7, 8},
			preAddGames: []struct {
				gameID    int
				playerIDs [4]int
				status    models.GameStatus
				finished  bool
			}{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameEnded,
					finished:  true,
				},
				{
					gameID:    2,
					playerIDs: [4]int{1, 5, 6, 7},
					status:    models.GameEnded,
					finished:  true,
				},
			},
			playerID:    1,
			expectedErr: assert.NoError,
			expectedGame: 2, // Should return the most recent (highest ID)
		},
		"ok filters by status": {
			preAddPlayers: []int{1, 2, 3, 4, 5, 6, 7, 8},
			preAddGames: []struct {
				gameID    int
				playerIDs [4]int
				status    models.GameStatus
				finished  bool
			}{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameLive,
					finished:  false,
				},
				{
					gameID:    2,
					playerIDs: [4]int{1, 5, 6, 7},
					status:    models.GameEnded,
					finished:  true,
				},
			},
			playerID:    1,
			expectedErr: assert.NoError,
			expectedGame: 2, // Should return only ended game
		},
		"not found": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames:   []struct {
				gameID    int
				playerIDs [4]int
				status    models.GameStatus
				finished  bool
			}{},
			playerID:    1,
			expectedErr: assert.Error,
		},
		"player not in any game": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				playerIDs [4]int
				status    models.GameStatus
				finished  bool
			}{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
					status:    models.GameEnded,
					finished:  true,
				},
			},
			playerID:    999,
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

			for _, id := range testCase.preAddPlayers {
				createTestUser(t, conn, id)
			}

			for _, cfg := range testCase.preAddGames {
				createTestGameWithStatus(t, conn, cfg.gameID, cfg.status)
				createTestGamePlayers(t, conn, cfg.gameID, cfg.playerIDs)
				if cfg.finished {
					setGameFinishedAt(t, conn, cfg.gameID)
				}
			}

			gameID, err := repo.LastGame(ctx, testCase.playerID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Equal(t, testCase.expectedGame, gameID)
		})
	}
}
