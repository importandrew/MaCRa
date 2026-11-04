package player

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

func TestUpdate(t *testing.T) {
	testData := map[string]struct {
		name        string
		callsRepo   bool
		expectedErr error
	}{
		"ok": {
			name:      "Alice",
			callsRepo: true,
		},
		"cyrillic name over 50 bytes but under 50 characters": {
			name:      "Александр Константиновский", // 26 characters, 51 bytes
			callsRepo: true,
		},
		"exactly 50 cyrillic characters": {
			name:      strings.Repeat("Я", 50),
			callsRepo: true,
		},
		"51 characters": {
			name:        strings.Repeat("я", 51),
			expectedErr: models.ErrLongName,
		},
		"empty": {
			name:        "",
			expectedErr: models.ErrEmptyName,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockPlayerRepo(t)
			if testCase.callsRepo {
				mockRepo.On("Update", mock.Anything, 1, testCase.name).Return(nil)
			}

			service := NewService(mockRepo, nil)

			err := service.Update(context.Background(), 1, testCase.name)

			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
