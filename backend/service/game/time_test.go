package game

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

// drawsAt returns draw entries, one per time, hand by hand.
func drawsAt(times ...time.Time) []models.RoundLog {
	logs := make([]models.RoundLog, 0, len(times))
	for i, at := range times {
		logs = append(logs, models.RoundLog{ID: i + 1, Round: i + 1, EventSequence: 1, Outcome: models.Draw, CreatedAt: at})
	}
	return logs
}

func TestLoad_timeUp(t *testing.T) {
	endsAt := time.Now().Add(-time.Hour)
	before, after := endsAt.Add(-time.Minute), endsAt.Add(time.Minute)

	testData := map[string]struct {
		rule       string
		extraHands int
		logs       []models.RoundLog
		status     models.GameStatus
		lastHand   int
		stopped    bool
	}{
		"hand in progress finished after the time": {
			rule: models.TimeUpFinishHand, logs: drawsAt(before, before, after),
			status: models.GameEnded, lastHand: 3,
		},
		"hand in progress still to play": {
			rule: models.TimeUpFinishHand, logs: drawsAt(before, before),
			status: models.GameLive, lastHand: 3,
		},
		"one more hand": {
			rule: models.TimeUpExtraHands, extraHands: 1, logs: drawsAt(before, before, after),
			status: models.GameLive, lastHand: 4,
		},
		"one more hand played": {
			rule: models.TimeUpExtraHands, extraHands: 1, logs: drawsAt(before, before, after, after),
			status: models.GameEnded, lastHand: 4,
		},
		"the game ends at once": {
			rule: models.TimeUpStopNow, logs: drawsAt(before, before),
			status: models.GameEnded, lastHand: 16, stopped: true,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockGameRepo(t)
			expectLoadGame(repo, 100, testCase.logs)
			repo.EXPECT().TimeLimit(mock.Anything, 100).Return(&models.GameTimeLimit{
				EndsAt: endsAt, TimeUp: testCase.rule, ExtraHands: testCase.extraHands,
			}, nil)

			game, _, err := NewService(repo, NewMockPlayerRepo(t)).load(context.Background(), 100)

			require.NoError(t, err)
			assert.Equal(t, testCase.status, game.Status())
			assert.Equal(t, testCase.lastHand, game.LastHand())
			assert.Equal(t, testCase.stopped, game.StoppedAtTime())
		})
	}
}

func TestLoad_beforeTheTimeIsUp(t *testing.T) {
	repo := NewMockGameRepo(t)
	expectLoadGame(repo, 100, drawsAt(time.Now().Add(-time.Minute)))
	repo.EXPECT().TimeLimit(mock.Anything, 100).Return(&models.GameTimeLimit{
		EndsAt: time.Now().Add(time.Hour), TimeUp: models.TimeUpStopNow,
	}, nil)

	game, _, err := NewService(repo, NewMockPlayerRepo(t)).load(context.Background(), 100)

	require.NoError(t, err)
	assert.Equal(t, models.GameLive, game.Status())
	assert.Equal(t, 16, game.LastHand())
	assert.Nil(t, gameTimer(game, nil))
	timer := gameTimer(game, &models.GameTimeLimit{EndsAt: time.Now().Add(time.Hour)})
	assert.Zero(t, timer.LastHand, "no last hand before the time is up")
}
