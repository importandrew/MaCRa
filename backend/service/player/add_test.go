package player

import (
	"context"
	"errors"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAdd(t *testing.T) {
	testData := map[string]struct {
		email    string
		name     string
		password string

		mockSetup func(*MockPlayerRepo, string, string, string)

		expectedErr    assert.ErrorAssertionFunc
		expectedPlayer models.Player
	}{
		"ok": {
			email:    "alice@example.com",
			name:     "Alice",
			password: "password123",
			mockSetup: func(m *MockPlayerRepo, email, name, hash string) {
				m.On("Add", mock.Anything, email, name, mock.MatchedBy(func(h string) bool {
					// Verify that the hash is a valid bcrypt hash
					err := bcrypt.CompareHashAndPassword([]byte(h), []byte("password123"))
					return err == nil
				})).Return(1, nil)
			},
			expectedErr: assert.NoError,
			expectedPlayer: models.Player{
				ID:   1,
				Name: "Alice",
			},
		},
		"repo error": {
			email:    "bob@example.com",
			name:     "Bob",
			password: "password123",
			mockSetup: func(m *MockPlayerRepo, email, name, hash string) {
				m.On("Add", mock.Anything, email, name, mock.Anything).Return(
					0,
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"duplicate email": {
			email:    "existing@example.com",
			name:     "Existing",
			password: "password123",
			mockSetup: func(m *MockPlayerRepo, email, name, hash string) {
				m.On("Add", mock.Anything, email, name, mock.Anything).Return(
					0,
					models.ErrPlayerExists,
				)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockRepo, testCase.email, testCase.name, "")

			service := NewService(mockRepo)

			ctx := context.Background()
			player, err := service.Add(ctx, testCase.email, testCase.name, testCase.password)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, testCase.expectedPlayer.ID, player.ID)
			assert.Equal(t, testCase.expectedPlayer.Name, player.Name)
		})
	}
}

func TestHashPassword(t *testing.T) {
	testData := map[string]struct {
		password      string
		expectedError assert.ErrorAssertionFunc
	}{
		"valid password": {
			password:      "password123",
			expectedError: assert.NoError,
		},
		"empty password": {
			password:      "",
			expectedError: assert.NoError,
		},
		"long password": {
			password:      "a very long password with many characters and symbols !@#$%^&*()",
			expectedError: assert.NoError,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			hash, err := hashPassword(testCase.password)
			testCase.expectedError(t, err)
			if err != nil {
				return
			}

			// Verify the hash can be used to validate the original password
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(testCase.password))
			assert.NoError(t, err, "hash should be valid for the original password")
		})
	}
}
