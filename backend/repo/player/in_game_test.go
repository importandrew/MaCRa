package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
)

func TestIsInLiveGame(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []int
		preAddGames []struct {
			gameID    int
			status    models.GameStatus
			playerIDs [4]int
		}
		playerID int

		expectedErr    assert.ErrorAssertionFunc
		expectedInGame bool
	}{
		"ok in live game": {
			preAddUsers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameLive,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			playerID:       1,
			expectedErr:    assert.NoError,
			expectedInGame: true,
		},
		"ok in created game": {
			preAddUsers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameCreated,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			playerID:       1,
			expectedErr:    assert.NoError,
			expectedInGame: true,
		},
		"ok not in game": {
			preAddUsers: []int{1, 2, 3, 4, 5},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameLive,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			playerID:       5,
			expectedErr:    assert.NoError,
			expectedInGame: false,
		},
		"ok in ended game": {
			preAddUsers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameEnded,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			playerID:       1,
			expectedErr:    assert.NoError,
			expectedInGame: false,
		},
		"ok in aborted game": {
			preAddUsers: []int{1, 2, 3, 4},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameAborted,
					playerIDs: [4]int{1, 2, 3, 4},
				},
			},
			playerID:       1,
			expectedErr:    assert.NoError,
			expectedInGame: false,
		},
		"ok multiple games": {
			preAddUsers: []int{1, 2, 3, 4, 5, 6, 7, 8},
			preAddGames: []struct {
				gameID    int
				status    models.GameStatus
				playerIDs [4]int
			}{
				{
					gameID:    1,
					status:    models.GameEnded,
					playerIDs: [4]int{1, 2, 3, 4},
				},
				{
					gameID:    2,
					status:    models.GameLive,
					playerIDs: [4]int{5, 6, 7, 8},
				},
			},
			playerID:       5,
			expectedErr:    assert.NoError,
			expectedInGame: true,
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

			for _, cfg := range testCase.preAddGames {
				createTestGame(t, conn, cfg.gameID, cfg.status)
				winds := []string{"east", "south", "west", "north"}
				for i, playerID := range cfg.playerIDs {
					createTestGamePlayer(t, conn, cfg.gameID, playerID, winds[i])
				}
			}

			inGame, err := repo.IsInLiveGame(ctx, testCase.playerID)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedInGame, inGame)
		})
	}
}
