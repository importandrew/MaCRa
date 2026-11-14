// Package seating makes a tournament tour's tables (docs/user-scenarios.md,
// area 8, "Seating"): random with as few repeat meetings as possible, Swiss
// (neighbours in the current standings) avoiding repeats, and interval (by
// the standings with a step). One generator for all of them: a method gives a
// starting seating and a cost, and a local search swaps players between
// tables while the cost goes down, until a deadline.
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
// of tables when k is 0): table i gets places i, i+k, i+2k, i+3k. No search:
// the method fixes the tables.
func Interval(input Input, k int) [][4]int {
	ranked := Ranked(input)
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
	return seatByNumber(tables)
}

// Random seats the participants at random with as few repeat meetings as
// the search finds before the deadline.
func Random(input Input, deadline time.Time) ([][4]int, error) {
	rng, err := newRNG()
	if err != nil {
		return nil, err
	}
	cost := func(tables [][]models.Participant) int {
		return Repeats(seatByNumber(tables), input.Meetings)
	}

	var best [][]models.Participant
	bestCost := -1
	for {
		start := slices.Clone(input.Participants)
		rng.Shuffle(len(start), func(i, j int) { start[i], start[j] = start[j], start[i] })
		tables, c := improve(chunk(start), cost, rng, deadline)
		if bestCost < 0 || c < bestCost {
			best, bestCost = tables, c
		}
		if bestCost == 0 || time.Now().After(deadline) {
			return seatByNumber(best), nil
		}
	}
}

// swissRepeatWeight makes a repeat meeting cost more than any spread of
// places at a table, so the search avoids repeats first and keeps
// neighbours together second.
const swissRepeatWeight = 1_000_000

// Swiss seats neighbours in the standings together (places 1–4, 5–8, …),
// moving players apart only to avoid repeat meetings, as far as the search
// finds before the deadline.
func Swiss(input Input, deadline time.Time) ([][4]int, error) {
	rng, err := newRNG()
	if err != nil {
		return nil, err
	}
	ranked := Ranked(input)
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
		return swissRepeatWeight*Repeats(seatByNumber(tables), input.Meetings) + spread
	}

	tables, _ := improve(chunk(ranked), cost, rng, deadline)
	return seatByNumber(tables), nil
}

// improve swaps players between tables while that lowers the cost (first
// improvement, in a random order), until no swap helps or the deadline.
func improve(tables [][]models.Participant, cost func([][]models.Participant) int,
	rng *mrand.Rand, deadline time.Time,
) ([][]models.Participant, int) {
	current := cost(tables)
	for improved := true; improved && current > 0; {
		improved = false
		for _, a := range rng.Perm(len(tables)) {
			for b := a + 1; b < len(tables); b++ {
				if time.Now().After(deadline) {
					return tables, current
				}
				var better bool
				current, better = swapBetween(tables, a, b, cost, current)
				improved = improved || better
			}
		}
	}
	return tables, current
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
