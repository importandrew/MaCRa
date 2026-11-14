package event

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestRandomTables(t *testing.T) {
	participants := make([]models.Participant, 0, 12)
	for n := 1; n <= 12; n++ {
		participants = append(participants, models.Participant{UserID: 100 + n, Number: n})
	}

	tables, err := randomTables(participants)

	require.NoError(t, err)
	require.Len(t, tables, 3)
	seated := []int{}
	for _, table := range tables {
		for i := 1; i < len(table); i++ {
			assert.Less(t, table[i-1], table[i], "winds by participant number, east the lowest")
		}
		seated = append(seated, table[:]...)
	}
	slices.Sort(seated)
	want := make([]int, 0, 12)
	for _, p := range participants {
		want = append(want, p.UserID)
	}
	assert.Equal(t, want, seated, "everyone seated once")
	assert.Len(t, participants, 12, "the input is left as it was")
	assert.Equal(t, 1, participants[0].Number)
}

func TestSeatingFunc_unknown(t *testing.T) {
	_, err := seatingFunc("swiss-ish")
	assert.ErrorIs(t, err, models.ErrInvalidSeating)
}
