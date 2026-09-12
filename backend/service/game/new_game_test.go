package game

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGame(t *testing.T) {
	testData := map[string]struct {
		playerIDs []int

		mockSetup func(*MockGameRepo, []int)

		expectedErr assert.ErrorAssertionFunc
		expectedID  int
	}{
		"ok": {
			playerIDs: []int{1, 2, 3, 4},
			mockSetup: func(m *MockGameRepo, playerIDs []int) {
				m.EXPECT().CurrentGame(mock.Anything, mock.Anything, testEventID).Return(0, models.ErrGameNotFound).Times(4)
				m.EXPECT().NewGame(mock.Anything, [4]int{1, 2, 3, 4}, testEventID).Return(100, nil)
			},
			expectedErr: assert.NoError,
			expectedID:  100,
		},
		"repo error": {
			playerIDs: []int{1, 2, 3, 4},
			mockSetup: func(m *MockGameRepo, playerIDs []int) {
				m.EXPECT().CurrentGame(mock.Anything, mock.Anything, testEventID).Return(0, models.ErrGameNotFound).Times(4)
				m.EXPECT().NewGame(mock.Anything, [4]int{1, 2, 3, 4}, testEventID).Return(0, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"wrong number of players": {
			playerIDs:   []int{1, 2, 3},
			mockSetup:   func(m *MockGameRepo, playerIDs []int) {},
			expectedErr: assert.Error,
		},
		"duplicate players": {
			playerIDs:   []int{1, 2, 3, 1},
			mockSetup:   func(m *MockGameRepo, playerIDs []int) {},
			expectedErr: assert.Error,
		},
		"zero player ID": {
			playerIDs:   []int{1, 2, 3, 0},
			mockSetup:   func(m *MockGameRepo, playerIDs []int) {},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockGameRepo, testCase.playerIDs)

			service := NewService(mockGameRepo, mockPlayerRepo)

			gameID, err := service.NewGame(context.Background(), testCase.playerIDs, testEventID)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedID, gameID)
		})
	}
}

func TestValidatePlayers(t *testing.T) {
	testData := map[string]struct {
		playerIDs   []int
		mockSetup   func(*MockGameRepo)
		expectedErr assert.ErrorAssertionFunc
	}{
		"valid four unique players": {
			playerIDs: []int{1, 2, 3, 4},
			mockSetup: func(m *MockGameRepo) {
				m.EXPECT().CurrentGame(mock.Anything, mock.Anything, testEventID).Return(0, models.ErrGameNotFound).Times(4)
			},
			expectedErr: assert.NoError,
		},
		"wrong number: three players": {
			playerIDs:   []int{1, 2, 3},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"wrong number: five players": {
			playerIDs:   []int{1, 2, 3, 4, 5},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"wrong number: empty": {
			playerIDs:   []int{},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"duplicate players": {
			playerIDs:   []int{1, 2, 3, 1},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"all same players": {
			playerIDs:   []int{1, 1, 1, 1},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"zero player ID": {
			playerIDs:   []int{1, 2, 3, 0},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"zero at first position": {
			playerIDs:   []int{0, 1, 2, 3},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
		"multiple zeros": {
			playerIDs:   []int{0, 0, 0, 0},
			mockSetup:   func(m *MockGameRepo) {},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockGameRepo)

			service := NewService(mockGameRepo, mockPlayerRepo)
			err := service.validatePlayers(context.Background(), testCase.playerIDs, testEventID)
			testCase.expectedErr(t, err)
		})
	}
}
