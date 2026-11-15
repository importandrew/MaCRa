package event

import (
	"cmp"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// tieBreak orders players equal in TP and GP (docs/user-scenarios.md, area
// 8, "Rating"): by the tournament's criteria in order, then by the draw
// number, the lower first.
type tieBreak struct {
	criteria []string
	// games are the finished tours' games: each game's players and GP.
	games []map[int]int
}

// newTieBreak gathers the games the criteria look at from the results.
func newTieBreak(criteria []string, results []models.TourResultRow) tieBreak {
	type gameKey struct{ tour, table int }
	index := map[gameKey]int{}
	var games []map[int]int
	for _, r := range results {
		key := gameKey{r.Tour, r.Table}
		i, ok := index[key]
		if !ok {
			i = len(games)
			index[key] = i
			games = append(games, map[int]int{})
		}
		games[i][r.PlayerID] = r.GP
	}
	return tieBreak{criteria: criteria, games: games}
}

// order sorts a run of players equal in TP and GP.
func (tb tieBreak) order(run []models.TourRatingRow) {
	scores := make([]map[int]int, len(tb.criteria))
	for i, criterion := range tb.criteria {
		switch criterion {
		case models.TieHeadToHead:
			scores[i] = tb.headToHead(run)
		case models.TieFirstPlaces:
			scores[i] = tb.firstPlaces(run)
		}
	}
	slices.SortStableFunc(run, func(a, b models.TourRatingRow) int {
		for _, score := range scores {
			if c := cmp.Compare(score[b.PlayerID], score[a.PlayerID]); c != 0 {
				return c
			}
		}
		return cmp.Compare(numberOr(a), numberOr(b))
	})
}

// headToHead counts, for each player of the run, how often they scored more
// GP than another player of the run in a game they played together.
func (tb tieBreak) headToHead(run []models.TourRatingRow) map[int]int {
	wins := map[int]int{}
	for _, game := range tb.games {
		for _, a := range run {
			for _, b := range run {
				gpA, okA := game[a.PlayerID]
				gpB, okB := game[b.PlayerID]
				if okA && okB && gpA > gpB {
					wins[a.PlayerID]++
				}
			}
		}
	}
	return wins
}

// firstPlaces counts, for each player of the run, the games where nobody at
// the table scored more GP.
func (tb tieBreak) firstPlaces(run []models.TourRatingRow) map[int]int {
	firsts := map[int]int{}
	for _, game := range tb.games {
		top := 0
		first := true
		for _, gp := range game {
			if first || gp > top {
				top, first = gp, false
			}
		}
		for _, p := range run {
			if gp, ok := game[p.PlayerID]; ok && gp == top {
				firsts[p.PlayerID]++
			}
		}
	}
	return firsts
}
