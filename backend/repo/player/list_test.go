package player

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type idName struct {
	id   int
	name string
}

func TestListByNameLimit(t *testing.T) {
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testData := map[string]struct {
		preAddUsers []idName
		searchName  string
		limit       int

		expectedErr     assert.ErrorAssertionFunc
		expectedPlayers []models.Player
	}{
		"ok match single": {
			preAddUsers: []idName{
				{1, "Alice"},
				{2, "Bob"},
			},
			searchName:  "alice",
			limit:       0,
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 1, Name: "Alice"},
			},
		},
		"ok match multiple": {
			preAddUsers: []idName{
				{1, "Alice"},
				{2, "Bob"},
				{3, "Charlie"},
			},
			searchName:  "a",
			limit:       0,
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 1, Name: "Alice"},
				{ID: 3, Name: "Charlie"},
			},
		},
		"ok with limit": {
			preAddUsers: []idName{
				{1, "Alice"},
				{2, "Bob"},
				{3, "Charlie"},
			},
			searchName:  "a",
			limit:       1,
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 1, Name: "Alice"},
			},
		},
		"ok no match": {
			preAddUsers: []idName{
				{1, "Alice"},
			},
			searchName:      "xyz",
			limit:           0,
			expectedErr:     assert.NoError,
			expectedPlayers: []models.Player{},
		},
		"ok empty": {
			searchName:      "alice",
			limit:           0,
			expectedErr:     assert.NoError,
			expectedPlayers: []models.Player{},
		},
		"empty name zero limit": {
			preAddUsers: []idName{
				{1, "Alice"},
				{2, "Bob"},
				{3, "Charlie"},
			},
			searchName:  "",
			limit:       0,
			expectedErr: assert.NoError,
			expectedPlayers: []models.Player{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 3, Name: "Charlie"},
			},
		},
	}

	for name, testCase := range testData {
		db.Test(t, name, func(t *testing.T) {
			repo := NewRepo(db)

			ctx := context.Background()
			conn, err := db.Acquire(ctx)
			if err != nil {
				t.Fatalf("Failed to acquire connection: %v", err)
			}
			defer conn.Release()

			for _, user := range testCase.preAddUsers {
				createTestUserWithName(t, conn, user.id, user.name)
			}

			players, err := repo.ListByNameLimit(ctx, testCase.searchName, testCase.limit)
			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			require.Len(
				t,
				players,
				len(testCase.expectedPlayers),
				"Expected %d players, got %d",
				len(testCase.expectedPlayers),
				len(players),
			)
			for i, expected := range testCase.expectedPlayers {
				assert.Equal(t, expected.ID, players[i].ID, "Player %d ID should match", i)
				assert.Equal(t, expected.Name, players[i].Name, "Player %d Name should match", i)
			}
		})
	}
}
