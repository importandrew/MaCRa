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

func TestLogin(t *testing.T) {
	validPassword := "password123"
	validHash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	wrongPassword := "wrongpassword"

	testData := map[string]struct {
		email    string
		password string

		mockSetup func(*MockPlayerRepo, string, string)

		expectedErr    assert.ErrorAssertionFunc
		expectedPlayer models.Player
	}{
		"ok": {
			email:    "alice@example.com",
			password: validPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					string(validHash),
					nil,
				)
				m.EXPECT().GetByEmail(mock.Anything, email).Return(
					models.Player{ID: 1, Name: "Alice"},
					nil,
				)
				m.EXPECT().IsEmailVerified(mock.Anything, 1).Return(true, nil)
			},
			expectedErr: assert.NoError,
			expectedPlayer: models.Player{
				ID:   1,
				Name: "Alice",
			},
		},
		"email not verified": {
			email:    "alice@example.com",
			password: validPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					string(validHash),
					nil,
				)
				m.EXPECT().GetByEmail(mock.Anything, email).Return(
					models.Player{ID: 1, Name: "Alice"},
					nil,
				)
				m.EXPECT().IsEmailVerified(mock.Anything, 1).Return(false, nil)
			},
			expectedErr: assert.Error,
		},
		"wrong password": {
			email:    "alice@example.com",
			password: wrongPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					string(validHash),
					nil,
				)
			},
			expectedErr: assert.Error,
		},
		"email not found": {
			email:    "nonexistent@example.com",
			password: validPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					"",
					models.ErrPlayerNotFound,
				)
			},
			expectedErr: assert.Error,
		},
		"get password error": {
			email:    "error@example.com",
			password: validPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					"",
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"get player error after successful password check": {
			email:    "alice@example.com",
			password: validPassword,
			mockSetup: func(m *MockPlayerRepo, email, password string) {
				m.EXPECT().GetPassByEmail(mock.Anything, email).Return(
					string(validHash),
					nil,
				)
				m.EXPECT().GetByEmail(mock.Anything, email).Return(
					models.Player{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockPlayerRepo(t)
			testCase.mockSetup(mockRepo, testCase.email, testCase.password)

			service := NewService(mockRepo, nil)

			ctx := context.Background()
			player, err := service.Login(ctx, testCase.email, testCase.password)

			testCase.expectedErr(t, err)
			if err != nil {
				// Verify that bad credentials error is returned for wrong password
				if testCase.password == wrongPassword {
					assert.Equal(t, models.ErrBadCredentials, err)
				}
				if name == "email not verified" {
					assert.Equal(t, models.ErrEmailNotVerified, err)
				}
				return
			}

			assert.Equal(t, testCase.expectedPlayer.ID, player.ID)
			assert.Equal(t, testCase.expectedPlayer.Name, player.Name)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validPassword := "password123"
	validHash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	wrongPassword := "wrongpassword"

	testData := map[string]struct {
		storedPass   string
		incomingPass string
		expectedErr  assert.ErrorAssertionFunc
	}{
		"valid password": {
			storedPass:   string(validHash),
			incomingPass: validPassword,
			expectedErr:  assert.NoError,
		},
		"wrong password": {
			storedPass:   string(validHash),
			incomingPass: wrongPassword,
			expectedErr:  assert.Error,
		},
		"empty password": {
			storedPass:   string(validHash),
			incomingPass: "",
			expectedErr:  assert.Error,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			err := validatePassword(testCase.storedPass, testCase.incomingPass)
			testCase.expectedErr(t, err)
		})
	}
}
