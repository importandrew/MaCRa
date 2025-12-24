package player

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListByNameLimit(t *testing.T) {
	testData := map[string]struct {
		name  string
		limit int

		mockSetup func(*MockPlayerRepo, string, int)

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers []models.Player
	}{
		"ok": {
			name:  "Alice",
			limit: 10,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					[]models.Player{
						{ID: 1, Name: "Alice"},
						{ID: 5, Name: "Alice Smith"},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 1, Name: "Alice"},
				{ID: 5, Name: "Alice Smith"},
			},
		},
		"empty result": {
			name:  "Nonexistent",
			limit: 10,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					[]models.Player{},
					nil,
				)
			},
			expectedErr:     assert.NoError,
			expectedPlayers: []models.Player{},
		},
		"limit zero": {
			name:  "Alice",
			limit: 0,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					[]models.Player{},
					nil,
				)
			},
			expectedErr:     assert.NoError,
			expectedPlayers: []models.Player{},
		},
		"limit one": {
			name:  "Bob",
			limit: 1,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					[]models.Player{
						{ID: 2, Name: "Bob"},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 2, Name: "Bob"},
			},
		},
		"repo error": {
			name:  "Charlie",
			limit: 10,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					nil,
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"empty name": {
			name:  "",
			limit: 10,
			mockSetup: func(m *MockPlayerRepo, name string, limit int) {
				m.On("ListByNameLimit", mock.Anything, name, limit).Return(
					[]models.Player{},
					nil,
				)
			},
			expectedErr:     assert.NoError,
			expectedPlayers: []models.Player{},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockRepo, testCase.name, testCase.limit)

			service := NewService(mockRepo)

			ctx := context.Background()
			players, err := service.ListByNameLimit(ctx, testCase.name, testCase.limit)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Len(t, players, len(testCase.expectedPlayers))
			for i, expectedPlayer := range testCase.expectedPlayers {
				assert.Equal(t, expectedPlayer.ID, players[i].ID)
				assert.Equal(t, expectedPlayer.Name, players[i].Name)
			}
		})
	}
}
