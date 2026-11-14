package event

import (
	"context"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/game"
)

// createTestEvent creates a regular event owned by user 1, with users 1-4 as
// its players, and returns its ID.
func createTestEvent(t *testing.T, db database.DB) int {
	t.Helper()

	for id := 1; id <= 4; id++ {
		createTestUser(t, db, id)
	}

	eventID, err := NewRepo(db).Create(context.Background(), models.Event{
		Type:         models.EventRegular,
		Name:         "Finish",
		OwnerUserID:  1,
		Visibility:   models.EventVisibilityPublic,
		RatingsReach: models.ContentReachPublic,
		GamesReach:   models.ContentReachPublic,
		Rules:        models.StandardGameRules(),
		Settings:     models.DefaultEventSettings(),
	})
	require.NoError(t, err)
	// The owner got player and referee at creation; games need the others
	// to be players too.
	for id := 2; id <= 4; id++ {
		require.NoError(t, NewRepo(db).AddMemberRoles(context.Background(), eventID, id, []models.EventRole{models.EventPlayer}))
	}

	return eventID
}

func TestFinish(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	seating := [4]int{1, 2, 3, 4}

	testData := map[string]struct {
		setup       func(t *testing.T, eventID int)
		eventID     func(eventID int) int
		expectedErr error
	}{
		"no games": {},
		"only ended games": {
			setup: func(t *testing.T, eventID int) {
				gameID, err := game.NewRepo(db).NewGame(context.Background(), seating, eventID)
				require.NoError(t, err)
				require.NoError(t, game.NewRepo(db).SetStatus(context.Background(), gameID, models.GameEnded))
			},
		},
		"a live game": {
			setup: func(t *testing.T, eventID int) {
				_, err := game.NewRepo(db).NewGame(context.Background(), seating, eventID)
				require.NoError(t, err)
			},
			expectedErr: models.ErrEventHasOngoingGames,
		},
		"already finished": {
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
			eventID := createTestEvent(t, db)
			if testCase.setup != nil {
				testCase.setup(t, eventID)
			}
			target := eventID
			if testCase.eventID != nil {
				target = testCase.eventID(eventID)
			}

			err := NewRepo(db).Finish(context.Background(), target)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				return
			}
			require.NoError(t, err)
			event, err := NewRepo(db).Get(context.Background(), eventID)
			require.NoError(t, err)
			assert.NotNil(t, event.FinishedAt)
		})
	}
}

// TestFinish_ConcurrentNewGame checks that finishing an event and starting a
// game in it at the same time never leave a live game in a finished event:
// exactly one of the two succeeds.
func TestFinish_ConcurrentNewGame(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "concurrent finish and new game", func(t *testing.T) {
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
			_, err = conn.Exec(ctx, `UPDATE events SET finished_at = NULL WHERE id = $1`, eventID)
			require.NoError(t, err)

			var finishErr, newGameErr error
			var wg sync.WaitGroup
			wg.Go(func() { finishErr = events.Finish(ctx, eventID) })
			wg.Go(func() { _, newGameErr = games.NewGame(ctx, [4]int{1, 2, 3, 4}, eventID) })
			wg.Wait()

			if finishErr == nil {
				assert.ErrorIs(t, newGameErr, models.ErrEventAlreadyFinished)
			} else {
				assert.ErrorIs(t, finishErr, models.ErrEventHasOngoingGames)
				assert.NoError(t, newGameErr)
			}

			rows, err := conn.Query(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM games g JOIN events e ON e.id = g.event_id
					WHERE e.id = $1 AND e.finished_at IS NOT NULL AND g.status = $2)`,
				eventID, models.GameLive)
			require.NoError(t, err)
			liveInFinished, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
			require.NoError(t, err)
			require.False(t, liveInFinished, "a finished event has a live game")
		}
	})
}
