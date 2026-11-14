package player

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/repo/database"
)

func TestSearchByName(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "search", func(t *testing.T) {
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()

		for id, name := range map[int]string{
			1: "Annabel", 2: "Anna", 3: "Joanna", 4: "anna", 5: "Bob", 6: "Ann_a", 7: "Ann%",
		} {
			createTestUserWithName(t, conn, id, name)
		}
		repo := NewRepo(db)
		names := func(query string, limit int) []string {
			players, err := repo.SearchByName(ctx, query, limit)
			require.NoError(t, err)
			out := make([]string, len(players))
			for i, p := range players {
				out[i] = p.Name
			}
			return out
		}

		// Exact (any case) first, then prefixes, then the rest; by name and ID.
		assert.Equal(t, []string{"Anna", "anna", "Annabel", "Joanna"}, names("anna", 20))
		assert.Equal(t, []string{"Anna", "anna"}, names("anna", 2), "the limit")
		// LIKE's wildcards are searched as typed.
		assert.Equal(t, []string{"Ann_a"}, names("n_a", 20))
		assert.Equal(t, []string{"Ann%"}, names("n%", 20))
		assert.Empty(t, names("zzz", 20))
	})
}
