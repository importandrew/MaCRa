package player

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestSearch(t *testing.T) {
	ivan := models.Player{ID: 42, Name: "Ivan"}

	testData := map[string]struct {
		query     string
		mockSetup func(*MockPlayerRepo)
		want      []models.Player
		wantErr   bool
	}{
		"an ID matches exactly": {
			query:     " 42 ",
			mockSetup: func(m *MockPlayerRepo) { m.EXPECT().Get(mock.Anything, 42).Return(ivan, nil) },
			want:      []models.Player{ivan},
		},
		"an unknown ID finds nobody": {
			query: "7",
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().Get(mock.Anything, 7).Return(models.Player{}, models.ErrPlayerNotFound)
			},
			want: []models.Player{},
		},
		"a name searches names, 20 at most": {
			query: " Iv ",
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().SearchByName(mock.Anything, "Iv", 20).Return([]models.Player{ivan}, nil)
			},
			want: []models.Player{ivan},
		},
		"one letter finds nobody": {query: "I", want: []models.Player{}},
		"empty finds nobody":      {query: "  ", want: []models.Player{}},
		"one Cyrillic letter":     {query: "И", want: []models.Player{}},
		"a signed number is a name": {query: "+4", mockSetup: func(m *MockPlayerRepo) {
			m.EXPECT().SearchByName(mock.Anything, "+4", 20).Return([]models.Player{}, nil)
		}, want: []models.Player{}},
		"repo error": {
			query: "Ivan",
			mockSetup: func(m *MockPlayerRepo) {
				m.EXPECT().SearchByName(mock.Anything, "Ivan", 20).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockPlayerRepo(t)
			if testCase.mockSetup != nil {
				testCase.mockSetup(repo)
			}

			got, err := NewService(repo, nil).Search(context.Background(), testCase.query)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}
