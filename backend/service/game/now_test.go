package game

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

func TestNow(t *testing.T) {
	testData := map[string]struct {
		playerID    int
		mockSetup   func(*MockGameRepo, *MockPlayerRepo)
		expectedErr assert.ErrorAssertionFunc
		checkNow    func(*testing.T, models.GameNow)
	}{
		"ok": {
			playerID: 1,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				expectCurrentGame(mg, 1, 100, nil, 1)
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, 1, now.Round.Num)
				assert.Equal(t, "Alice", now.Me.Name)
				assert.Equal(t, models.East.Kanji(), now.Me.Wind)
			},
		},
		"log of the same state, players in starting seat order": {
			playerID: 3,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				expectCurrentGame(mg, 3, 100, drawRoundLogs(3), 1)
				expectPlayerGets(mp)
			},
			expectedErr: assert.NoError,
			checkNow: func(t *testing.T, now models.GameNow) {
				assert.Equal(t, 3, now.LogSize)
				assert.Len(t, now.Log.RoundLogs, now.LogSize)
				assert.Equal(t, []models.PlayerHeader{
					{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}, {ID: 3, Name: "Charlie"}, {ID: 4, Name: "David"},
				}, now.Log.Players)
			},
		},
		"player not in game": {
			playerID: 999,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().CurrentGame(mock.Anything, 999, testEventID).Return(0, models.ErrGameNotFound)
			},
			expectedErr: assert.Error,
		},
		"player repo error": {
			playerID: 1,
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				expectCurrentGame(mg, 1, 100, nil, 1)
				mp.On("Get", mock.Anything, 1).Return(models.Player{}, errors.New("player not found"))
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo)

			now, err := service.Now(context.Background(), testCase.playerID, testEventID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.checkNow != nil {
				testCase.checkNow(t, now)
			}
		})
	}
}
