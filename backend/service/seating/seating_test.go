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

	tables, err := Swiss(input, true, time.Now().Add(time.Second))
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.ElementsMatch(t, [][4]int{{105, 106, 107, 108}, {101, 102, 103, 104}}, tables, "neighbours together")

	// Once the tables have met, they are split 2+2: with 8 players two
	// tablemates sit together again at each table, 4 repeats at the least.
	input.Meetings = meetingsOf([][4]int{{105, 106, 107, 108}, {101, 102, 103, 104}})
	tables, err = Swiss(input, true, time.Now().Add(time.Second))
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Equal(t, 4, Repeats(tables, input.Meetings))

	// Not avoiding repeats, the neighbours stay together.
	tables, err = Swiss(input, false, time.Now().Add(time.Second))
	require.NoError(t, err)
	assert.Equal(t, [][4]int{{105, 106, 107, 108}, {101, 102, 103, 104}}, tables)
}

func TestSwiss_noRepeatsWhenPossible(t *testing.T) {
	input := Input{Participants: participants(16), Standings: map[int]Standing{}}
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(p.Number)}
	}
	input.Meetings = meetingsOf([][4]int{
		{101, 102, 103, 104}, {105, 106, 107, 108}, {109, 110, 111, 112}, {113, 114, 115, 116},
	})

	tables, err := Swiss(input, true, time.Now().Add(2*time.Second))

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

	tables, err := Interval(input, 0, false, time.Now())
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Equal(t, [][4]int{
		{101, 105, 109, 113}, // places 1, 5, 9, 13
		{102, 106, 110, 114},
		{103, 107, 111, 115},
		{104, 108, 112, 116},
	}, tables, "step = number of tables")

	tables, err = Interval(input, 2, false, time.Now())
	require.NoError(t, err)
	assert.Equal(t, [][4]int{
		{101, 103, 105, 107}, // places 1, 3, 5, 7
		{102, 104, 106, 108},
		{109, 111, 113, 115},
		{110, 112, 114, 116},
	}, tables, "step 2: blocks of two tables")
}

func TestInterval_avoidsRepeats(t *testing.T) {
	input := Input{Participants: participants(16), Standings: map[int]Standing{}}
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(100 - p.Number)} // place = number
	}
	// 101 and 105 met: the scheme's table 1 (places 1, 5, 9, 13) repeats.
	input.Meetings = map[Pair]int{NewPair(101, 105): 1}

	tables, err := Interval(input, 0, false, time.Now())
	require.NoError(t, err)
	assert.Equal(t, 1, Repeats(tables, input.Meetings), "the scheme as it is")

	tables, err = Interval(input, 0, true, time.Now().Add(time.Second))
	require.NoError(t, err)
	assertSeatsEveryoneOnce(t, input, tables)
	assert.Zero(t, Repeats(tables, input.Meetings))
	// The cheapest fix swaps neighbours in places: 105 and 106 (or 101 and
	// 102) trade tables; the rest keep the scheme.
	moved := 0
	scheme := [][4]int{{101, 105, 109, 113}, {102, 106, 110, 114}, {103, 107, 111, 115}, {104, 108, 112, 116}}
	for i := range tables {
		for _, userID := range tables[i] {
			if !slices.Contains(scheme[i][:], userID) {
				moved++
			}
		}
	}
	assert.Equal(t, 2, moved, "two players moved: %v", tables)
}

func TestRandom_plansAhead(t *testing.T) {
	// 16 players, 5 random tours in all, seated one at a time: planning each
	// with the tours after it leaves no repeats at the end.
	input := Input{Participants: participants(16), Meetings: map[Pair]int{}}
	all := [][][4]int{}
	for tour := range 5 {
		input.FutureRandom = 4 - tour
		tables, err := Random(input, time.Now().Add(10*time.Second))
		require.NoError(t, err)
		assertSeatsEveryoneOnce(t, input, tables)
		all = append(all, tables)
		input.Meetings = meetingsOf(all...)
	}
	for pair, met := range input.Meetings {
		assert.Equal(t, 1, met, "pair %v", pair)
	}
}

func TestQuality(t *testing.T) {
	input := Input{Participants: participants(8), Standings: map[int]Standing{}}
	for _, p := range input.Participants {
		input.Standings[p.UserID] = Standing{TP: float64(100 - p.Number)} // place = number
	}
	input.Meetings = map[Pair]int{NewPair(101, 102): 2, NewPair(105, 108): 1}
	tables := [][4]int{{101, 102, 103, 104}, {105, 106, 107, 108}}

	quality, summary := Quality(tables, input)

	require.Len(t, quality, 2)
	assert.Equal(t, 2, quality[0].Repeats)
	assert.Equal(t, []models.RepeatPair{{A: 101, B: 102, Met: 2}}, quality[0].Pairs)
	assert.Equal(t, &models.PlaceRange{From: 1, To: 4}, quality[0].Places)
	assert.Equal(t, 1, quality[1].Repeats)
	assert.Equal(t, &models.PlaceRange{From: 5, To: 8}, quality[1].Places)
	assert.Equal(t, 3, summary.Repeats)
	assert.Equal(t, 2, summary.MaxMet)
	assert.Equal(t, 4, summary.PlayersWithRepeats)

	input.Standings = nil
	quality, _ = Quality(tables, input)
	assert.Nil(t, quality[0].Places, "no places before any tour is finished")
	assert.Empty(t, quality[0].Pairs[1:], "only the pairs who met")
}
