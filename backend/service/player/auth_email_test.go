package player

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

// fakeMailer records which auth emails were sent, as "kind to".
type fakeMailer struct {
	sent []string
}

func (m *fakeMailer) record(kind, to string) error {
	m.sent = append(m.sent, kind+" "+to)
	return nil
}

func (m *fakeMailer) SendVerification(to, _ string) error  { return m.record("verify", to) }
func (m *fakeMailer) SendPasswordReset(to, _ string) error { return m.record("reset", to) }
func (m *fakeMailer) SendAccountExists(to, _ string) error { return m.record("exists", to) }

func TestSignUp(t *testing.T) {
	const email = "alice@example.com"

	testData := map[string]struct {
		mockSetup    func(*MockPlayerRepo)
		expectedErr  error
		expectedSent []string
	}{
		"new address gets a verification link": {
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().Add(mock.Anything, email, "Alice", mock.Anything).Return(1, nil)
				m.EXPECT().IssueAuthToken(mock.Anything, 1, models.AuthPurposeEmailVerify, mock.Anything, mock.Anything, models.AuthMailLimits).Return(nil)
			},
			expectedSent: []string{"verify " + email},
		},
		"registered address gets the same answer and an account-exists mail": {
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().Add(mock.Anything, email, "Alice", mock.Anything).Return(0, models.ErrPlayerExists)
				m.EXPECT().GetByEmail(mock.Anything, email).Return(models.Player{ID: 7}, nil)
				m.EXPECT().IssueAuthToken(mock.Anything, 7, models.AuthPurposePasswordReset, mock.Anything, mock.Anything, models.AuthMailLimits).Return(nil)
			},
			expectedSent: []string{"exists " + email},
		},
		"over the address's limit nothing is sent, the answer is the same": {
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().Add(mock.Anything, email, "Alice", mock.Anything).Return(0, models.ErrPlayerExists)
				m.EXPECT().GetByEmail(mock.Anything, email).Return(models.Player{ID: 7}, nil)
				m.EXPECT().IssueAuthToken(mock.Anything, 7, models.AuthPurposePasswordReset, mock.Anything, mock.Anything, models.AuthMailLimits).
					Return(fmt.Errorf("1 in 1m0s: %w", models.ErrAuthMailLimit))
			},
		},
		"other errors are returned": {
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().Add(mock.Anything, email, "Alice", mock.Anything).Return(0, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockPlayerRepo(t)
			testCase.mockSetup(repo)
			mail := &fakeMailer{}

			err := NewService(repo, mail).SignUp(context.Background(), email, "Alice", "secret")

			assert.ErrorIs(t, err, testCase.expectedErr)
			assert.Equal(t, testCase.expectedSent, mail.sent)
		})
	}
}

func TestForgotPassword_overLimit(t *testing.T) {
	repo := NewMockPlayerRepo(t)
	repo.EXPECT().GetByEmail(mock.Anything, "a@example.com").Return(models.Player{ID: 3}, nil)
	repo.EXPECT().IssueAuthToken(mock.Anything, 3, models.AuthPurposePasswordReset, mock.Anything, mock.Anything, models.AuthMailLimits).
		Return(models.ErrAuthMailLimit)
	mail := &fakeMailer{}

	err := NewService(repo, mail).ForgotPassword(context.Background(), "a@example.com")

	assert.NoError(t, err, "the response says sent either way")
	assert.Empty(t, mail.sent)
}
