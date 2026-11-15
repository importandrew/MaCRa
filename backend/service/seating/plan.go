package seating

import (
	"math"
	mrand "math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/importandrew/MaCRa/models"
)

// Fixed seats the tour by its tables imported from a file (input.Tables),
// winds as there. It refuses without them (models.ErrFileSeatingMissing),
// and tables that seat numbers who don't play (models.ErrFileSeatingAbsent)
// or leave out participants (models.ErrFileSeatingUnseated).
func Fixed(input Input) ([][4]int, error) {
	if input.Tables == nil {
		return nil, models.WithData(models.ErrFileSeatingMissing, models.ErrorData{dataTour: input.Tour})
	}

	userByNumber := make(map[int]int, len(input.Participants))
	for _, p := range input.Participants {
		userByNumber[p.Number] = p.UserID
	}
	seated := map[int]bool{}
	var absent []int
	tables := make([][4]int, 0, len(input.Tables))
	for _, numbers := range input.Tables {
		var table [4]int
		for seat, number := range numbers {
			userID, ok := userByNumber[number]
			if !ok {
				absent = append(absent, number)
			}
			table[seat] = userID
			seated[number] = true
		}
		tables = append(tables, table)
	}
	if len(absent) > 0 {
		return nil, models.WithData(models.ErrFileSeatingAbsent,
			models.ErrorData{dataTour: input.Tour, "numbers": numberList(absent)})
	}

	var unseated []int
	for number := range userByNumber {
		if !seated[number] {
			unseated = append(unseated, number)
		}
	}
	if len(unseated) > 0 {
		return nil, models.WithData(models.ErrFileSeatingUnseated,
			models.ErrorData{dataTour: input.Tour, "numbers": numberList(unseated)})
	}
	return tables, nil
}

// Keys of the data the seating errors carry.
const (
	dataTour = "tour"
	dataLine = "line"
)

// numberList returns participant numbers in order, as "5, 12".
func numberList(numbers []int) string {
	slices.Sort(numbers)
	words := make([]string, len(numbers))
	for i, n := range numbers {
		words[i] = strconv.Itoa(n)
	}
	return strings.Join(words, ", ")
}

// Plan seats the participants for the given number of tours at once, with
// as few repeat meetings over all of them (and the meetings before) as the
// search finds before the deadline. It returns each tour's tables as
// participant numbers, winds by number (east the lowest).
//
// The search is simulated annealing over swaps of two players at different
// tables of one tour; a pair meeting T times in all costs repeatsOf(T), the
// repeat meetings as Repeats counts them.
func Plan(input Input, tours int, deadline time.Time) ([][][4]int, error) {
	rng, err := newRNG()
	if err != nil {
		return nil, err
	}
	tableCount := len(input.Participants) / models.SeatsPerTable
	if tours <= 0 || tableCount == 0 {
		return [][][4]int{}, nil
	}

	search := newPlanSearch(input, tours, rng)
	floor := 0 // the cost of the meetings before, which no plan avoids
	for _, m := range input.Meetings {
		floor += repeatsOf(m)
	}
	best, bestCost := cloneTours(search.plan), search.cost
	start := time.Now()
	total := deadline.Sub(start)
	temperature := 0.0
	for step := 0; bestCost > floor && tableCount > 1; step++ {
		if step%1024 == 0 {
			elapsed := time.Since(start)
			if elapsed >= total {
				break
			}
			// Cooling from 1 to 0 over the time given.
			temperature = 1 - float64(elapsed)/float64(total)
		}
		if search.try(temperature) && search.cost < bestCost {
			best, bestCost = cloneTours(search.plan), search.cost
		}
	}

	result := make([][][4]int, 0, tours)
	for _, seats := range best {
		tables := make([][]models.Participant, 0, tableCount)
		for table := range tableCount {
			players := make([]models.Participant, 0, models.SeatsPerTable)
			for _, i := range seats[table*models.SeatsPerTable : (table+1)*models.SeatsPerTable] {
				players = append(players, input.Participants[i])
			}
			tables = append(tables, players)
		}
		result = append(result, numbersByNumber(tables))
	}
	return result, nil
}

// repeatsOf is the repeat meetings of a pair that meets met times: each
// meeting after the first counts every earlier one.
func repeatsOf(met int) int {
	return met * (met - 1) / 2 //nolint:mnd // the pairs among met meetings
}

// planSearch is a plan in the making: plan[t][k] is the player (an index in
// the input's participants) at seat k%4 of table k/4 in tour t; met[i][j]
// how often players i and j meet, before and in the plan; cost the repeat
// meetings of met.
type planSearch struct {
	plan   [][]int
	met    [][]int
	cost   int
	tables int
	rng    *mrand.Rand
}

// newPlanSearch seats the participants at random in every tour.
func newPlanSearch(input Input, tours int, rng *mrand.Rand) *planSearch {
	n := len(input.Participants)
	s := &planSearch{
		plan:   make([][]int, tours),
		met:    make([][]int, n),
		tables: n / models.SeatsPerTable,
		rng:    rng,
	}
	for i, a := range input.Participants {
		s.met[i] = make([]int, n)
		for j, b := range input.Participants {
			if i != j {
				s.met[i][j] = input.Meetings[NewPair(a.UserID, b.UserID)]
			}
		}
	}
	for t := range s.plan {
		s.plan[t] = rng.Perm(n)[:s.tables*models.SeatsPerTable]
		for table := range s.tables {
			s.meet(s.table(t, table))
		}
	}
	for i := range n {
		for j := i + 1; j < n; j++ {
			s.cost += repeatsOf(s.met[i][j])
		}
	}
	return s
}

// meet counts the meetings of the players at a table.
func (s *planSearch) meet(table []int) {
	for _, a := range table {
		for _, b := range table {
			if a != b {
				s.met[a][b]++
			}
		}
	}
}

// table returns the players at a table of tour t.
func (s *planSearch) table(t, table int) []int {
	return s.plan[t][table*models.SeatsPerTable : (table+1)*models.SeatsPerTable]
}

// try picks two players at different tables of one tour and swaps them if
// that lowers the cost, or at the temperature's odds if it raises it. It
// reports whether it swapped.
func (s *planSearch) try(temperature float64) bool {
	t := s.rng.IntN(len(s.plan))
	a, b := s.rng.IntN(s.tables), s.rng.IntN(s.tables-1)
	if b >= a {
		b++
	}
	ka := a*models.SeatsPerTable + s.rng.IntN(models.SeatsPerTable)
	kb := b*models.SeatsPerTable + s.rng.IntN(models.SeatsPerTable)
	x, y := s.plan[t][ka], s.plan[t][kb]

	// Moving x to b and y to a: the pairs at a lose x and gain y, at b the
	// other way round. Losing a pair met T times saves T−1, gaining one met
	// T times costs T.
	delta := s.moveDelta(s.table(t, a), x, y) + s.moveDelta(s.table(t, b), y, x)
	if delta > 0 && (temperature == 0 || s.rng.Float64() >= math.Exp(-float64(delta)/temperature)) {
		return false
	}
	s.move(s.table(t, a), x, y)
	s.move(s.table(t, b), y, x)
	s.plan[t][ka], s.plan[t][kb] = y, x
	s.cost += delta
	return true
}

// moveDelta is the change of the cost when out leaves a table and in takes
// its place.
func (s *planSearch) moveDelta(table []int, out, in int) int {
	delta := 0
	for _, u := range table {
		if u != out {
			delta += s.met[in][u] - (s.met[out][u] - 1)
		}
	}
	return delta
}

// move counts the meetings of in, instead of out, at a table.
func (s *planSearch) move(table []int, out, in int) {
	for _, u := range table {
		if u != out {
			s.met[out][u]--
			s.met[u][out]--
			s.met[in][u]++
			s.met[u][in]++
		}
	}
}

// cloneTours copies a plan in the making.
func cloneTours(plan [][]int) [][]int {
	clone := make([][]int, len(plan))
	for t := range plan {
		clone[t] = slices.Clone(plan[t])
	}
	return clone
}

// numbersByNumber returns each table's participant numbers east to north,
// the lowest east (as seatByNumber).
func numbersByNumber(tables [][]models.Participant) [][4]int {
	seated := make([][4]int, 0, len(tables))
	for _, table := range tables {
		numbers := make([]int, 0, len(table))
		for _, p := range table {
			numbers = append(numbers, p.Number)
		}
		slices.Sort(numbers)
		seated = append(seated, [4]int(numbers))
	}
	return seated
}

// FileInfo describes tables imported from a file for the desk: their tours
// in order, the most tables in a tour, how many numbers they seat, and their
// repeat meetings among their own tours (each earlier meeting of each pair
// seated together) and the most times one pair meets in them.
func FileInfo(tours PlanTours) models.SeatingFileInfo {
	info := models.SeatingFileInfo{Tours: []int{}}
	numbers := map[int]bool{}
	meetings := map[Pair]int{}
	for tour, tables := range tours {
		info.Tours = append(info.Tours, tour)
		info.Tables = max(info.Tables, len(tables))
		for _, table := range tables {
			for i, a := range table {
				numbers[a] = true
				for _, b := range table[i+1:] {
					meetings[NewPair(a, b)]++
				}
			}
		}
	}
	slices.Sort(info.Tours)
	info.Players = len(numbers)
	for _, m := range meetings {
		info.Repeats += repeatsOf(m)
		info.MaxMet = max(info.MaxMet, m)
	}
	return info
}
