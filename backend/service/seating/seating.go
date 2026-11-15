// Package seating makes a tournament tour's tables (docs/user-scenarios.md,
// area 8, "Seating"): random with as few repeat meetings as possible (also
// with the random tours planned after it), Swiss (neighbours in the current
// standings) and interval (by the standings with a step), each avoiding
// repeats or not, and tables imported from a file. Swiss and interval give a
// starting seating and a cost, and a local search swaps players between
// tables while the cost goes down, until a deadline; random tours are
// planned together by simulated annealing (Plan).
package seating

import (
	"crypto/rand"
	mrand "math/rand/v2"
	"slices"
	"time"

	"github.com/importandrew/MaCRa/models"
)

// Pair, Standing and Input are the models' seating types, short here.
type (
	Pair     = models.SeatingPair
	Standing = models.SeatingStanding
	Input    = models.SeatingInput
	// PlanTours is tables of several tours by participant number.
	PlanTours = models.PlanTours
)

// NewPair returns the pair of a and b in order.
func NewPair(a, b int) Pair {
	if a > b {
		a, b = b, a
	}
	return Pair{a, b}
}

// Repeats counts the repeat meetings at the tables: each earlier meeting of
// each pair seated together counts once.
func Repeats(tables [][4]int, meetings map[Pair]int) int {
	total := 0
	for _, table := range tables {
		for i := range table {
			for j := i + 1; j < len(table); j++ {
				total += meetings[NewPair(table[i], table[j])]
			}
		}
	}
	return total
}

// Ranked returns the participants by the standings: TP, then GP, both
// descending, then by participant number.
func Ranked(input Input) []models.Participant {
	ranked := slices.Clone(input.Participants)
	slices.SortStableFunc(ranked, func(a, b models.Participant) int {
		sa, sb := input.Standings[a.UserID], input.Standings[b.UserID]
		switch {
		case sa.TP != sb.TP:
			if sa.TP > sb.TP {
				return -1
			}
			return 1
		case sa.GP != sb.GP:
			return sb.GP - sa.GP
		default:
			return a.Number - b.Number
		}
	})
	return ranked
}

// Interval seats the participants by the standings with step k (the number
// of tables when k is 0): table i gets places i, i+k, i+2k, i+3k. With
// avoidRepeats it then swaps players to avoid repeat meetings, keeping each
// as near to their place in the scheme as the search finds before the
// deadline.
func Interval(input Input, k int, avoidRepeats bool, deadline time.Time) ([][4]int, error) {
	ranked := Ranked(input)
	tables := intervalTables(ranked, k)
	if !avoidRepeats {
		return seatByNumber(tables), nil
	}
	rng, err := newRNG()
	if err != nil {
		return nil, err
	}

	// slot[t][s] is the place the scheme gives seat s of table t.
	place := make(map[int]int, len(ranked))
	for i, p := range ranked {
		place[p.UserID] = i
	}
	slot := make([][]int, len(tables))
	for t, table := range tables {
		for _, p := range table {
			slot[t] = append(slot[t], place[p.UserID])
		}
	}
	cost := func(tables [][]models.Participant) int {
		moved := 0
		for t, table := range tables {
			for s, p := range table {
				moved += abs(place[p.UserID] - slot[t][s])
			}
		}
		return repeatWeight*Repeats(seatByNumber(tables), input.Meetings) + moved
	}
	tables = improve(tables, cost, rng, deadline)
	return seatByNumber(tables), nil
}

// abs returns |n|.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// intervalTables makes the interval scheme's tables from the ranked
// participants (see Interval).
func intervalTables(ranked []models.Participant, k int) [][]models.Participant {
	tableCount := len(ranked) / models.SeatsPerTable
	if k <= 0 || k > tableCount {
		k = tableCount
	}

	tables := make([][]models.Participant, tableCount)
	next := 0 // the next table to fill, cycling
	for start := 0; start < len(ranked); start += k * models.SeatsPerTable {
		block := ranked[start:min(start+k*models.SeatsPerTable, len(ranked))]
		rows := len(block) / models.SeatsPerTable
		for t := range rows {
			table := make([]models.Participant, 0, models.SeatsPerTable)
			for seat := range models.SeatsPerTable {
				table = append(table, block[t+seat*rows])
			}
			tables[next] = table
			next++
		}
	}
	return tables
}

// Random seats the participants at random with as few repeat meetings as
// the search finds before the deadline, planning input.FutureRandom more
// random tours along with it (so it doesn't use up the pairs they need);
// only this tour's tables are returned.
func Random(input Input, deadline time.Time) ([][4]int, error) {
	tours, err := Plan(input, 1+input.FutureRandom, deadline)
	if err != nil || len(tours) == 0 {
		return nil, err
	}
	userByNumber := make(map[int]int, len(input.Participants))
	for _, p := range input.Participants {
		userByNumber[p.Number] = p.UserID
	}
	tables := make([][4]int, 0, len(tours[0]))
	for _, numbers := range tours[0] {
		var table [4]int
		for seat, number := range numbers {
			table[seat] = userByNumber[number]
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// repeatWeight makes a repeat meeting cost more than any spread or move of
// places, so the Swiss and interval searches avoid repeats first and keep
// to their scheme second.
const repeatWeight = 1_000_000

// Swiss seats neighbours in the standings together (places 1–4, 5–8, …).
// With avoidRepeats it moves players apart to avoid repeat meetings, as far
// as the search finds before the deadline.
func Swiss(input Input, avoidRepeats bool, deadline time.Time) ([][4]int, error) {
	ranked := Ranked(input)
	if !avoidRepeats {
		return seatByNumber(chunk(ranked)), nil
	}
	rng, err := newRNG()
	if err != nil {
		return nil, err
	}
	place := make(map[int]int, len(ranked))
	for i, p := range ranked {
		place[p.UserID] = i
	}
	cost := func(tables [][]models.Participant) int {
		spread := 0
		for _, table := range tables {
			lowest, highest := len(ranked), -1
			for _, p := range table {
				lowest, highest = min(lowest, place[p.UserID]), max(highest, place[p.UserID])
			}
			spread += highest - lowest
		}
		return repeatWeight*Repeats(seatByNumber(tables), input.Meetings) + spread
	}

	tables := improve(chunk(ranked), cost, rng, deadline)
	return seatByNumber(tables), nil
}

// improve swaps players between tables while that lowers the cost (first
// improvement, in a random order), until no swap helps or the deadline.
func improve(tables [][]models.Participant, cost func([][]models.Participant) int,
	rng *mrand.Rand, deadline time.Time,
) [][]models.Participant {
	current := cost(tables)
	for improved := true; improved && current > 0; {
		improved = false
		for _, a := range rng.Perm(len(tables)) {
			for b := a + 1; b < len(tables); b++ {
				if time.Now().After(deadline) {
					return tables
				}
				var better bool
				current, better = swapBetween(tables, a, b, cost, current)
				improved = improved || better
			}
		}
	}
	return tables
}

// swapBetween tries every swap of a player at table a with one at table b,
// keeping those that lower the cost, and returns the cost and whether it
// went down.
func swapBetween(tables [][]models.Participant, a, b int, cost func([][]models.Participant) int, current int) (int, bool) {
	improved := false
	for i := range tables[a] {
		for j := range tables[b] {
			tables[a][i], tables[b][j] = tables[b][j], tables[a][i]
			if c := cost(tables); c < current {
				current, improved = c, true
				continue
			}
			tables[a][i], tables[b][j] = tables[b][j], tables[a][i]
		}
	}
	return current, improved
}

// chunk groups the participants into tables of four, in order.
func chunk(participants []models.Participant) [][]models.Participant {
	tables := make([][]models.Participant, 0, len(participants)/models.SeatsPerTable)
	for start := 0; start+models.SeatsPerTable <= len(participants); start += models.SeatsPerTable {
		tables = append(tables, slices.Clone(participants[start:start+models.SeatsPerTable]))
	}
	return tables
}

// seatByNumber returns each table's user IDs east to north: the winds go by
// participant number, east the lowest (area 8: table and wind by participant
// number).
func seatByNumber(tables [][]models.Participant) [][4]int {
	seated := make([][4]int, 0, len(tables))
	for _, table := range tables {
		sorted := slices.Clone(table)
		slices.SortFunc(sorted, func(a, b models.Participant) int { return a.Number - b.Number })
		seated = append(seated, [4]int{sorted[0].UserID, sorted[1].UserID, sorted[2].UserID, sorted[3].UserID})
	}
	return seated
}

// newRNG returns a generator seeded from crypto/rand: nobody should be able
// to predict a seating.
func newRNG() (*mrand.Rand, error) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return nil, err
	}
	return mrand.New(mrand.NewChaCha8(seed)), nil //nolint:gosec // ChaCha8 seeded from crypto/rand
}

// Quality describes the tables of a tour: per table its repeat meetings
// (the pairs who met before, and how often) and the range of places in the
// standings (nil without standings); and for the tour the repeat meetings,
// the most times one pair has met before and how many players meet someone
// again.
func Quality(tables [][4]int, input Input) ([]TableQuality, models.TourState) {
	var place map[int]int
	if len(input.Standings) > 0 {
		place = map[int]int{}
		for i, p := range Ranked(input) {
			place[p.UserID] = i + 1
		}
	}

	quality := make([]TableQuality, 0, len(tables))
	var summary models.TourState
	withRepeats := map[int]bool{}
	for _, table := range tables {
		q := TableQuality{Pairs: []models.RepeatPair{}}
		for i := range table {
			for j := i + 1; j < len(table); j++ {
				met := input.Meetings[NewPair(table[i], table[j])]
				if met == 0 {
					continue
				}
				q.Repeats += met
				q.Pairs = append(q.Pairs, models.RepeatPair{A: table[i], B: table[j], Met: met})
				summary.MaxMet = max(summary.MaxMet, met)
				withRepeats[table[i]], withRepeats[table[j]] = true, true
			}
		}
		if place != nil {
			q.Places = &models.PlaceRange{From: len(place) + 1}
			for _, userID := range table {
				q.Places.From = min(q.Places.From, place[userID])
				q.Places.To = max(q.Places.To, place[userID])
			}
		}
		summary.Repeats += q.Repeats
		quality = append(quality, q)
	}
	summary.PlayersWithRepeats = len(withRepeats)
	return quality, summary
}

// TableQuality is one table's part of Quality.
type TableQuality struct {
	Repeats int
	Pairs   []models.RepeatPair
	Places  *models.PlaceRange
}
