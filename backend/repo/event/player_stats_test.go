package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

func TestPlayerStats(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "player stats", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		games := game.NewRepo(db)
		eventID := createTestEvent(t, db) // players 1–4, seated east to north

		stats, err := repo.PlayerStats(ctx, eventID, 2)
		require.NoError(t, err)
		assert.Equal(t, "user2", stats.Name)
		assert.Zero(t, stats.Games, "a player without games")
		assert.Nil(t, stats.BestGP)

		createTestUser(t, db, 9)
		_, err = repo.PlayerStats(ctx, eventID, 9)
		assert.ErrorIs(t, err, models.ErrPlayerNotFound, "not in the event")

		entry := func(gameID, logSize int, data models.RoundData) {
			t.Helper()
			data.LogSize = logSize
			require.NoError(t, games.SaveRoundOutcome(ctx, gameID, data))
		}
		finish := func(gameID int, gp [4]int, tp [4]float32) {
			t.Helper()
			player := func(i int) models.FinalPlayerData {
				return models.FinalPlayerData{ID: i + 1, GP: gp[i], TP: tp[i]}
			}
			require.NoError(t, games.Finalize(ctx, models.FinalGameData{
				GameID: gameID, East: player(0), South: player(1), West: player(2), North: player(3),
			}))
		}

		// Game 1: player 2 wins by self-draw (10), then by a deal-in from 1
		// (20), then a draw; player 2 is 1st.
		g1, err := games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID)
		require.NoError(t, err)
		entry(g1, 0, models.RoundData{Round: 1, Outcome: models.Self, WinnerID: 2, Value: 10})
		entry(g1, 1, models.RoundData{Round: 2, Outcome: models.Deal, WinnerID: 2, LoserID: 1, Value: 20})
		entry(g1, 2, models.RoundData{Round: 3, Outcome: models.Draw})
		finish(g1, [4]int{-60, 120, -30, -30}, [4]float32{0, 4, 1.5, 1.5})

		// Game 2: player 2 deals in to 3 (8) and gets a penalty; 4th.
		g2, err := games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID)
		require.NoError(t, err)
		entry(g2, 0, models.RoundData{Round: 1, Outcome: models.Penalty, LoserID: 2, Value: 10})
		entry(g2, 1, models.RoundData{Round: 1, Outcome: models.Deal, WinnerID: 3, LoserID: 2, Value: 8})
		finish(g2, [4]int{8, -50, 50, -8}, [4]float32{2, 0, 4, 1})

		// A live game doesn't count.
		g3, err := games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID)
		require.NoError(t, err)
		entry(g3, 0, models.RoundData{Round: 1, Outcome: models.Self, WinnerID: 2, Value: 88})

		stats, err = repo.PlayerStats(ctx, eventID, 2)
		require.NoError(t, err)
		assert.Equal(t, 2, stats.Games)
		assert.Equal(t, [4]int{1, 0, 0, 1}, stats.Places)
		assert.InDelta(t, 2.0, stats.AvgTP, 0.001)
		assert.InDelta(t, 35.0, stats.AvgGP, 0.001)
		assert.Equal(t, 120, *stats.BestGP)
		assert.Equal(t, -50, *stats.WorstGP)

		assert.Equal(t, 4, stats.Hands, "game 1: 3 hands, game 2: 1 (a penalty and a deal-in in the same hand)")
		assert.Equal(t, 1, stats.SelfDraws)
		assert.Equal(t, 1, stats.DealInWins)
		assert.InDelta(t, 15.0, stats.AvgWin, 0.001)
		assert.Equal(t, 20, stats.MaxWin)
		assert.Equal(t, 1, stats.DealIns)
		assert.Equal(t, 1, stats.Draws)
		assert.Equal(t, 1, stats.Penalties)
		assert.Zero(t, stats.FalseMahjongs)

		// Shared places: west and north tied at -30 in game 1 both count as
		// 2nd (120, -30, -30, -60); north is 3rd in game 2.
		stats, err = repo.PlayerStats(ctx, eventID, 4)
		require.NoError(t, err)
		assert.Equal(t, [4]int{0, 1, 1, 0}, stats.Places)
	})
}
