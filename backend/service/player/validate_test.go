package player

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

// errorIs returns an assertion that the error matches target.
func errorIs(target error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, _ ...any) bool {
		return assert.ErrorIs(t, err, target)
	}
}

func TestCheckName(t *testing.T) {
	testData := map[string]struct {
		name        string
		expectedErr assert.ErrorAssertionFunc
	}{
		"ok":                  {name: "Alice", expectedErr: assert.NoError},
		"empty":               {name: "", expectedErr: errorIs(models.ErrEmptyName)},
		"50 characters":       {name: strings.Repeat("a", 50), expectedErr: assert.NoError},
		"51 characters":       {name: strings.Repeat("a", 51), expectedErr: errorIs(models.ErrLongName)},
		"50 Cyrillic letters": {name: strings.Repeat("я", 50), expectedErr: assert.NoError},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			testCase.expectedErr(t, checkName(testCase.name))
		})
	}
}

func TestCheckEmail(t *testing.T) {
	// 50 characters in all: 38 + "@example.com".
	longest := strings.Repeat("a", 38) + "@example.com"

	testData := map[string]struct {
		email       string
		expectedErr assert.ErrorAssertionFunc
	}{
		"ok":                {email: "alice@example.com", expectedErr: assert.NoError},
		"subaddress":        {email: "alice+macra@example.com", expectedErr: assert.NoError},
		"50 characters":     {email: longest, expectedErr: assert.NoError},
		"51 characters":     {email: "a" + longest, expectedErr: errorIs(models.ErrLongEmail)},
		"empty":             {email: "", expectedErr: errorIs(models.ErrEmptyEmail)},
		"no at sign":        {email: "alice", expectedErr: errorIs(models.ErrInvalidEmail)},
		"no domain":         {email: "alice@", expectedErr: errorIs(models.ErrInvalidEmail)},
		"no local part":     {email: "@example.com", expectedErr: errorIs(models.ErrInvalidEmail)},
		"space inside":      {email: "ali ce@example.com", expectedErr: errorIs(models.ErrInvalidEmail)},
		"surrounding space": {email: " alice@example.com", expectedErr: errorIs(models.ErrInvalidEmail)},
		"display name":      {email: "Alice <alice@example.com>", expectedErr: errorIs(models.ErrInvalidEmail)},
		"two addresses":     {email: "a@example.com, b@example.com", expectedErr: errorIs(models.ErrInvalidEmail)},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			testCase.expectedErr(t, checkEmail(testCase.email))
		})
	}
}

func TestCheckNewPassword(t *testing.T) {
	testData := map[string]struct {
		password    string
		expectedErr assert.ErrorAssertionFunc
	}{
		"one character": {password: "a", expectedErr: assert.NoError},
		"72 bytes":      {password: strings.Repeat("a", 72), expectedErr: assert.NoError},
		"73 bytes":      {password: strings.Repeat("a", 73), expectedErr: errorIs(models.ErrLongPassword)},
		// 37 Cyrillic letters are 74 bytes: bcrypt counts bytes.
		"37 Cyrillic letters": {password: strings.Repeat("я", 37), expectedErr: errorIs(models.ErrLongPassword)},
		"empty":               {password: "", expectedErr: errorIs(models.ErrEmptyPassword)},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			testCase.expectedErr(t, checkNewPassword(testCase.password))

			// Every password that passes the check can be hashed.
			if checkNewPassword(testCase.password) == nil {
				_, err := hashPassword(testCase.password)
				assert.NoError(t, err)
			}
		})
	}
}

func TestResetPassword_refusesUnusablePassword(t *testing.T) {
	// The password is checked before the token, so nothing reaches the repo.
	service := NewService(NewMockPlayerRepo(t), nil)

	err := service.ResetPassword(context.Background(), "token", strings.Repeat("a", 73))

	assert.ErrorIs(t, err, models.ErrLongPassword)
}
