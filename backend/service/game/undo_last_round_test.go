package game

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestUndoLastRound_passesCutoff(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	service := NewService(mockGameRepo, NewMockPlayerRepo(t))

	before := time.Now()
	mockGameRepo.EXPECT().
		UndoLastRound(mock.Anything, 100, 7, mock.MatchedBy(func(endedAfter *time.Time) bool {
			cutoff := before.Add(-models.UndoFinalRoundTimeout)
			return endedAfter != nil && !endedAfter.Before(cutoff) && endedAfter.Before(cutoff.Add(time.Minute))
		})).
		Return(nil)

	require.NoError(t, service.UndoLastRound(context.Background(), 100, 7, true))
}

func TestUndoLastRound_withoutReviving(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	service := NewService(mockGameRepo, NewMockPlayerRepo(t))
	mockGameRepo.EXPECT().UndoLastRound(mock.Anything, 100, 7, (*time.Time)(nil)).Return(nil)

	require.NoError(t, service.UndoLastRound(context.Background(), 100, 7, false))
}
