package seating

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

// participants returns n participants, user ID 100+number.
func participants(n int) []models.Participant {
	out := make([]models.Participant, 0, n)
	for number := 1; number <= n; number++ {
		out = append(out, models.Participant{UserID: 100 + number, Number: number})
	}
	return out
}

// meetingsOf counts the pairs of the tables as meetings.
func meetingsOf(tables ...[][4]int) map[Pair]int {
	meetings := map[Pair]int{}
	for _, tour := range tables {
		for _, table := range tour {
			for i := range 4 {
				for j := i + 1; j < 4; j++ {
					meetings[NewPair(table[i], table[j])]++
				}
			}
		}
	}
	return meetings
}

// assertSeatsEveryoneOnce checks every participant sits once, winds by number.
func assertSeatsEveryoneOnce(t *testing.T, input Input, tables [][4]int) {
	t.Helper()
	number := map[int]int{}
	for _, p := range input.Participants {
		number[p.UserID] = p.Number
	}
	seated := []int{}
	for _, table := range tables {
		for i := 1; i < 4; i++ {
			assert.Less(t, number[table[i-1]], number[table[i]], "winds by participant number")
		}
		seated = append(seated, table[:]...)
	}
	slices.Sort(seated)
	want := make([]int, 0, len(input.Participants))
	for _, p := range input.Participants {
		want = append(want, p.UserID)
	}
	slices.Sort(want)
	assert.Equal(t, want, seated)
}

func TestRepeats(t *testing.T) {
	meetings := map[Pair]int{NewPair(1, 2): 2, NewPair(3, 4): 1}
	assert.Equal(t, 3, Repeats([][4]int{{1, 2, 3, 4}}, meetings))
	assert.Equal(t, 0, Repeats([][4]int{{1, 3, 5, 7}, {2, 4, 6, 8}}, meetings))
}

func TestRandom_avoidsRepeats(t *testing.T) {
	// 16 players after one tour of tables 1–4, 5–8, …: a second tour without
	// repeats exists (e.g. 1,5,9,13 …), and the search should find one.
	input := Input{Participants: participants(16)}
	firstTour := [][4]int{{101, 102, 103, 104}, {105, 106, 107, 108}, {109, 110, 111, 112}, {113, 114, 115, 116}}
	input.Meetings = meetingsOf(firstTour)

	tables, err := Random(input, time.Now().Add(2*time.Second))

	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Zero(t, Repeats(tables, input.Meetings))
}

func TestSwiss(t *testing.T) {
	input := Input{Participants: participants(8), Standings: map[int]Standing{}}
	// Standings: 108 first, then 107, … 101 last.
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(p.Number)}
	}

	tables, err := Swiss(input, time.Now().Add(time.Second))
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.ElementsMatch(t, [][4]int{{105, 106, 107, 108}, {101, 102, 103, 104}}, tables, "neighbours together")

	// Once the tables have met, they are split 2+2: with 8 players two
	// tablemates sit together again at each table, 4 repeats at the least.
	input.Meetings = meetingsOf([][4]int{{105, 106, 107, 108}, {101, 102, 103, 104}})
	tables, err = Swiss(input, time.Now().Add(time.Second))
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Equal(t, 4, Repeats(tables, input.Meetings))
}

func TestSwiss_noRepeatsWhenPossible(t *testing.T) {
	input := Input{Participants: participants(16), Standings: map[int]Standing{}}
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(p.Number)}
	}
	input.Meetings = meetingsOf([][4]int{
		{101, 102, 103, 104}, {105, 106, 107, 108}, {109, 110, 111, 112}, {113, 114, 115, 116},
	})

	tables, err := Swiss(input, time.Now().Add(2*time.Second))

	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Zero(t, Repeats(tables, input.Meetings))
}

func TestRanked_tiesByGPThenNumber(t *testing.T) {
	input := Input{Participants: participants(3), Standings: map[int]Standing{
		101: {TP: 4, GP: 10}, 102: {TP: 4, GP: 20}, 103: {TP: 4, GP: 20},
	}}
	ranked := Ranked(input)
	assert.Equal(t, []int{102, 103, 101}, []int{ranked[0].UserID, ranked[1].UserID, ranked[2].UserID})
}

func TestInterval(t *testing.T) {
	input := Input{Participants: participants(16), Standings: map[int]Standing{}}
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(100 - p.Number)} // place = number
	}

	tables := Interval(input, 0)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Equal(t, [][4]int{
		{101, 105, 109, 113}, // places 1, 5, 9, 13
		{102, 106, 110, 114},
		{103, 107, 111, 115},
		{104, 108, 112, 116},
	}, tables, "step = number of tables")

	tables = Interval(input, 2)
	assert.Equal(t, [][4]int{
		{101, 103, 105, 107}, // places 1, 3, 5, 7
		{102, 104, 106, 108},
		{109, 111, 113, 115},
		{110, 112, 114, 116},
	}, tables, "step 2: blocks of two tables")
}
