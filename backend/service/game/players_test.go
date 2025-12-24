package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlayers(t *testing.T) {
	testData := map[string]struct {
		gameID int

		mockSetup func(*MockGameRepo, *MockPlayerRepo, int)

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers models.GamePlayers
	}{
		"ok": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{ID: 2, Name: "Bob"}, nil)
				mp.On("Get", mock.Anything, 3).Return(models.Player{ID: 3, Name: "Charlie"}, nil)
				mp.On("Get", mock.Anything, 4).Return(models.Player{ID: 4, Name: "David"}, nil)
			},
			expectedErr: assert.NoError,
			expectedPlayers: models.GamePlayers{
				ID: 100,
				Players: []models.Player{
					{ID: 1, Name: "Alice"},
					{ID: 2, Name: "Bob"},
					{ID: 3, Name: "Charlie"},
					{ID: 4, Name: "David"},
				},
			},
		},
		"game repo error": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"player repo error on east": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{}, errors.New("player not found"))
			},
			expectedErr: assert.Error,
		},
		"player repo error on south": {
			gameID: 100,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo, gameID int) {
				mg.On("Players", mock.Anything, gameID).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mp.On("Get", mock.Anything, 1).Return(models.Player{ID: 1, Name: "Alice"}, nil)
				mp.On("Get", mock.Anything, 2).Return(models.Player{}, errors.New("player not found"))
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo, testCase.gameID)

			ctx := context.Background()
			players, err := service.Players(ctx, testCase.gameID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, 100, players.ID)
			assert.Len(t, players.Players, len(testCase.expectedPlayers.Players))
			for i, expected := range testCase.expectedPlayers.Players {
				assert.Equal(t, expected.ID, players.Players[i].ID)
				assert.Equal(t, expected.Name, players.Players[i].Name)
			}
		})
	}
}
