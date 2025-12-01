package game

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestSetStatus(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type preAddGame struct {
		gameID    int
		playerIDs [4]int
	}

	testData := map[string]struct {
		preAddGames   []preAddGame
		preAddPlayers []int
		gameID        int
		status        models.GameStatus

		expectedErr    assert.ErrorAssertionFunc
		expectedStatus models.GameStatus
	}{
		"ok": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID:         1,
			status:         models.GameAborted,
			expectedErr:    assert.NoError,
			expectedStatus: models.GameAborted,
		},
		"game not exists": {
			gameID:      1,
			status:      models.GameLive,
			expectedErr: assert.Error,
		},
		"status not exists": {
			preAddPlayers: []int{1, 2, 3, 4},
			preAddGames: []preAddGame{
				{
					gameID:    1,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			gameID:      1,
			status:      models.GameStatus("not exists"),
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
				createTestGame(t, conn, cfg.gameID)
				createTestGamePlayers(t, conn, cfg.gameID, cfg.playerIDs)
			}

			err = repo.SetStatus(ctx, testCase.gameID, testCase.status)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			status := getGameStatus(t, conn, testCase.gameID)
			assert.Equal(t, testCase.expectedStatus, status)
		})
	}
}
