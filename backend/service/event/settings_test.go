package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestUpdateSettings(t *testing.T) {
	referees := models.EventSettings{FinishGame: models.FinishByReferees, AbortGame: models.AbortByReferees}

	t.Run("saved", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1}, nil)
		repo.EXPECT().UpdateSettings(mock.Anything, 1, referees).Return(nil)

		require.NoError(t, NewService(repo, nil, nil).UpdateSettings(context.Background(), 1, referees))
	})

	t.Run("finished event", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, FinishedAt: new(time.Now())}, nil)

		err := NewService(repo, nil, nil).UpdateSettings(context.Background(), 1, referees)

		assert.ErrorIs(t, err, models.ErrEventAlreadyFinished)
	})

	t.Run("invalid settings", func(t *testing.T) {
		repo := NewMockEventRepo(t)

		err := NewService(repo, nil, nil).UpdateSettings(context.Background(), 1, models.EventSettings{FinishGame: "anyone"})

		assert.ErrorIs(t, err, models.ErrInvalidEventSettings)
	})
}
