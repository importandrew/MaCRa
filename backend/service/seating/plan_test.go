package seating

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

// planOf numbers a generated plan's tours from 1.
func planOf(tours [][][4]int) PlanTours {
	plan := PlanTours{}
	for i, tables := range tours {
		plan[i+1] = tables
	}
	return plan
}

// assertSeatsNumbersOnce checks a tour seats the numbers 1 to n once each,
// winds by number.
func assertSeatsNumbersOnce(t *testing.T, n int, tables [][4]int) {
	t.Helper()
	seated := []int{}
	for _, table := range tables {
		assert.True(t, slices.IsSorted(table[:]), "winds by participant number: %v", table)
		seated = append(seated, table[:]...)
	}
	slices.Sort(seated)
	want := make([]int, n)
	for i := range want {
		want[i] = i + 1
	}
	assert.Equal(t, want, seated)
}

func TestPlan_noRepeatsWhenPossible(t *testing.T) {
	// 16 players can play 5 tours without meeting anyone twice.
	tours, err := Plan(Input{Participants: participants(16)}, 5, time.Now().Add(20*time.Second))
	require.NoError(t, err)
	require.Len(t, tours, 5)
	for _, tables := range tours {
		assertSeatsNumbersOnce(t, 16, tables)
	}
	assert.Zero(t, FileInfo(planOf(tours)).Repeats)
}

func TestPlan_beyondTheLimit(t *testing.T) {
	// 7 tours of 16 players make 168 meetings for 120 pairs: 48 repeats at
	// least. The plan still seats everyone each tour.
	tours, err := Plan(Input{Participants: participants(16)}, 7, time.Now().Add(200*time.Millisecond))
	require.NoError(t, err)
	require.Len(t, tours, 7)
	for _, tables := range tours {
		assertSeatsNumbersOnce(t, 16, tables)
	}
	assert.GreaterOrEqual(t, FileInfo(planOf(tours)).Repeats, 48)
}

func TestPlan_avoidsEarlierMeetings(t *testing.T) {
	// Tour 1 was played: 1–4, 5–8. The one remaining tour can't avoid
	// repeats (it runs to the deadline), but the fewest come from two
	// players of each tour 1 table at each table.
	input := Input{
		Participants: participants(8),
		Meetings:     meetingsOf([][4]int{{101, 102, 103, 104}, {105, 106, 107, 108}}),
	}
	tours, err := Plan(input, 1, time.Now().Add(300*time.Millisecond))
	require.NoError(t, err)
	require.Len(t, tours, 1)
	assertSeatsNumbersOnce(t, 8, tours[0])
	for _, table := range tours[0] {
		low := 0
		for _, number := range table {
			if number <= 4 {
				low++
			}
		}
		assert.Equal(t, 2, low, "two from each table of tour 1: %v", table)
	}
}

func TestPlan_oneTable(t *testing.T) {
	tours, err := Plan(Input{Participants: participants(4)}, 3, time.Now().Add(time.Second))
	require.NoError(t, err)
	assert.Equal(t, [][][4]int{{{1, 2, 3, 4}}, {{1, 2, 3, 4}}, {{1, 2, 3, 4}}}, tours)
}

func TestFixed(t *testing.T) {
	input := Input{Participants: participants(8), Tour: 2, Tables: [][4]int{{4, 3, 2, 1}, {5, 6, 7, 8}}}

	tables, err := Fixed(input)
	require.NoError(t, err)
	assert.Equal(t, [][4]int{{104, 103, 102, 101}, {105, 106, 107, 108}}, tables, "winds as in the file")

	_, err = Fixed(Input{Participants: participants(8), Tour: 3})
	require.ErrorIs(t, err, models.ErrFileSeatingMissing)
	assertData(t, err, models.ErrorData{"tour": 3})

	// Numbers 3 and 7 withdrew.
	input.Participants = slices.DeleteFunc(participants(8), func(p models.Participant) bool {
		return p.Number == 3 || p.Number == 7
	})
	_, err = Fixed(input)
	require.ErrorIs(t, err, models.ErrFileSeatingAbsent)
	assertData(t, err, models.ErrorData{"tour": 2, "numbers": "3, 7"})

	// Numbers 9 to 12 joined after the file.
	input.Participants = participants(12)
	_, err = Fixed(input)
	require.ErrorIs(t, err, models.ErrFileSeatingUnseated)
	assertData(t, err, models.ErrorData{"tour": 2, "numbers": "9, 10, 11, 12"})
}

// assertData checks the data an error carries for the client.
func assertData(t *testing.T, err error, want models.ErrorData) {
	t.Helper()
	var dataErr *models.DataError
	require.ErrorAs(t, err, &dataErr)
	assert.Equal(t, want, dataErr.Data)
}

func TestFileInfo(t *testing.T) {
	info := FileInfo(PlanTours{
		3: {{1, 2, 3, 4}, {5, 6, 7, 8}},
		1: {{1, 2, 5, 6}, {3, 4, 7, 8}},
		2: {{1, 2, 3, 4}},
	})
	// 1–2 meet three times (3 repeats), 3–4 three times (3), 1–3, 1–4, 2–3,
	// 2–4 twice (1 each), 5–6 and 7–8 twice (1 each).
	assert.Equal(t, models.SeatingFileInfo{
		Tours: []int{1, 2, 3}, Tables: 2, Players: 8, Repeats: 12, MaxMet: 3,
	}, info)
}

func TestFixed_substitutes(t *testing.T) {
	// Number 3 withdrew; substitute 900 takes it, winds as in the file.
	input := Input{Tour: 2, Tables: [][4]int{{4, 3, 2, 1}, {5, 6, 7, 8}}}
	input.Participants = slices.DeleteFunc(participants(8), func(p models.Participant) bool { return p.Number == 3 })
	input.Participants = append(input.Participants,
		models.Participant{UserID: 900, Number: models.SubstituteNumber, Substitute: true})

	tables, err := Fixed(input)
	require.NoError(t, err)
	assert.Equal(t, [][4]int{{104, 900, 102, 101}, {105, 106, 107, 108}}, tables)

	// A second substitute has no seat.
	input.Participants = append(input.Participants,
		models.Participant{UserID: 901, Number: models.SubstituteNumber + 1, Substitute: true})
	_, err = Fixed(input)
	assert.ErrorIs(t, err, models.ErrTooManySubstitutes)
}

func TestSwiss_substitutesLast(t *testing.T) {
	// Substitutes have stand-in numbers after the players: last in the
	// standings, and after the players at a table's winds.
	input := Input{Participants: participants(7), Standings: map[int]Standing{}}
	input.Participants = append(input.Participants,
		models.Participant{UserID: 900, Number: models.SubstituteNumber, Substitute: true})
	for _, p := range participants(7) {
		input.Standings[p.UserID] = Standing{TP: float64(10 - p.Number)} // place = number
	}

	tables, err := Swiss(input, false, time.Now())
	require.NoError(t, err)
	assert.Equal(t, [][4]int{{101, 102, 103, 104}, {105, 106, 107, 900}}, tables)
}
