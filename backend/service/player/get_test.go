package player

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	testData := map[string]struct {
		playerID int

		mockSetup func(*MockPlayerRepo, int)

		expectedErr    assert.ErrorAssertionFunc
		expectedPlayer models.Player
	}{
		"ok": {
			playerID: 1,
			mockSetup: func(m *MockPlayerRepo, id int) {
				m.On("Get", mock.Anything, id).Return(
					models.Player{ID: 1, Name: "Alice"},
					nil,
				)
			},
			expectedErr:    assert.NoError,
			expectedPlayer: models.Player{ID: 1, Name: "Alice"},
		},
		"repo error": {
			playerID: 1,
			mockSetup: func(m *MockPlayerRepo, id int) {
				m.On("Get", mock.Anything, id).Return(
					models.Player{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"player not found": {
			playerID: 999,
			mockSetup: func(m *MockPlayerRepo, id int) {
				m.On("Get", mock.Anything, id).Return(
					models.Player{},
					errors.New("player not found"),
				)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockRepo, testCase.playerID)

			service := NewService(mockRepo, nil)

			ctx := context.Background()
			player, err := service.Get(ctx, testCase.playerID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedPlayer.ID, player.ID)
			assert.Equal(t, testCase.expectedPlayer.Name, player.Name)
		})
	}
}
