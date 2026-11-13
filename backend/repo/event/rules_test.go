package event

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

var shortRules = models.GameRules{Hands: 8, MinHand: 24, FalseMahjong: [2]int{30, 60}}

func TestUpdateRules(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	seating := [4]int{1, 2, 3, 4}
	newGame := func(t *testing.T, eventID int) int {
		gameID, err := game.NewRepo(db).NewGame(context.Background(), seating, eventID)
		require.NoError(t, err)
		return gameID
	}

	testData := map[string]struct {
		setup       func(t *testing.T, eventID int)
		eventID     func(eventID int) int
		expectedErr error
	}{
		"no games": {},
		"only aborted games": {
			setup: func(t *testing.T, eventID int) {
				gameID := newGame(t, eventID)
				require.NoError(t, game.NewRepo(db).Abort(context.Background(), gameID, false))
			},
		},
		"a live game": {
			setup:       func(t *testing.T, eventID int) { newGame(t, eventID) },
			expectedErr: models.ErrEventRulesFrozen,
		},
		"an ended game": {
			setup: func(t *testing.T, eventID int) {
				gameID := newGame(t, eventID)
				require.NoError(t, game.NewRepo(db).SetStatus(context.Background(), gameID, models.GameEnded))
			},
			expectedErr: models.ErrEventRulesFrozen,
		},
		"finished event": {
			setup: func(t *testing.T, eventID int) {
				require.NoError(t, NewRepo(db).Finish(context.Background(), eventID))
			},
			expectedErr: models.ErrEventAlreadyFinished,
		},
		"no such event": {
			eventID:     func(eventID int) int { return eventID + 1 },
			expectedErr: models.ErrEventNotFound,
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			ctx := context.Background()
			eventID := createTestEvent(t, db)
			if testCase.setup != nil {
				testCase.setup(t, eventID)
			}
			target := eventID
			if testCase.eventID != nil {
				target = testCase.eventID(eventID)
			}

			err := NewRepo(db).UpdateRules(ctx, target, shortRules)

			event, getErr := NewRepo(db).Get(ctx, eventID)
			require.NoError(t, getErr)
			frozen, frozenErr := NewRepo(db).RulesFrozen(ctx, eventID)
			require.NoError(t, frozenErr)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Equal(t, models.StandardGameRules(), event.Rules)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, shortRules, event.Rules)
			assert.False(t, frozen)
		})
	}
}

// TestNewGame_CopiesEventRules checks that a game keeps the rules its event
// had when it started.
func TestNewGame_CopiesEventRules(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "copies rules", func(t *testing.T) {
		ctx := context.Background()
		eventID := createTestEvent(t, db)
		require.NoError(t, NewRepo(db).UpdateRules(ctx, eventID, shortRules))

		gameID, err := game.NewRepo(db).NewGame(ctx, [4]int{1, 2, 3, 4}, eventID)
		require.NoError(t, err)

		rules, err := game.NewRepo(db).Rules(ctx, gameID)
		require.NoError(t, err)
		assert.Equal(t, shortRules, rules)

		frozen, err := NewRepo(db).RulesFrozen(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, frozen)
	})
}

// TestUpdateRules_ConcurrentNewGame checks that changing the rules and
// starting the first game at the same time never leave a game with rules
// other than the event's: either the change is refused, or the game has the
// new rules.
func TestUpdateRules_ConcurrentNewGame(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "concurrent rules change and new game", func(t *testing.T) {
		ctx := context.Background()
		eventID := createTestEvent(t, db)
		events := NewRepo(db)
		games := game.NewRepo(db)

		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for range 20 {
			_, err := conn.Exec(ctx, `DELETE FROM games WHERE event_id = $1`, eventID)
			require.NoError(t, err)
			_, err = conn.Exec(ctx, `UPDATE events SET rules = $1 WHERE id = $2`, models.StandardGameRules(), eventID)
			require.NoError(t, err)

			var rulesErr, newGameErr error
			var gameID int
			var wg sync.WaitGroup
			wg.Go(func() { rulesErr = events.UpdateRules(ctx, eventID, shortRules) })
			wg.Go(func() { gameID, newGameErr = games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID) })
			wg.Wait()

			require.NoError(t, newGameErr)
			rules, err := games.Rules(ctx, gameID)
			require.NoError(t, err)
			if rulesErr == nil {
				assert.Equal(t, shortRules, rules)
			} else {
				require.ErrorIs(t, rulesErr, models.ErrEventRulesFrozen)
				assert.Equal(t, models.StandardGameRules(), rules)
			}
		}
	})
}
