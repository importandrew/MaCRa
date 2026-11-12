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

func TestUndoLastRound_passesCutoff(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	service := NewService(mockGameRepo, NewMockPlayerRepo(t))

	before := time.Now()
	mockGameRepo.EXPECT().
		UndoLastRound(mock.Anything, 100, 7, mock.MatchedBy(func(endedAfter time.Time) bool {
			cutoff := before.Add(-models.UndoFinalRoundTimeout)
			return !endedAfter.Before(cutoff) && endedAfter.Before(cutoff.Add(time.Minute))
		})).
		Return(nil)

	require.NoError(t, service.UndoLastRound(context.Background(), 100, 7))
}

func TestGameResult_canUndoLastRound(t *testing.T) {
	ago := func(d time.Duration) *time.Time {
		t := time.Now().Add(-d)
		return &t
	}

	testData := map[string]struct {
		finishedAt *time.Time
		want       bool
	}{
		"ended just now":     {finishedAt: ago(time.Minute), want: true},
		"ended too long ago": {finishedAt: ago(models.UndoFinalRoundTimeout + time.Minute), want: false},
		"no finish time":     {finishedAt: nil, want: false},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			service := NewService(mockGameRepo, NewMockPlayerRepo(t))

			mockGameRepo.EXPECT().Status(mock.Anything, 100).Return(models.GameEnded, nil)
			mockGameRepo.EXPECT().GameResult(mock.Anything, 100).
				Return(models.GameResult{ID: 100, FinishedAt: testCase.finishedAt}, nil)

			result, err := service.GameResult(context.Background(), 100)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, result.CanUndoLastRound)
		})
	}
}
