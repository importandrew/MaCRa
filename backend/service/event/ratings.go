package event

import (
	"context"
	"sort"

	"github.com/importandrew/MaCRa/models"
)

type playerAgg struct {
	playerID  int
	name      string
	games     int
	tp        float32
	gp        int
	sessionGP []int // sorted descending
}

func (s *Service) Ratings(ctx context.Context, eventID int) ([]models.EventRatingEntry, error) {
	rows, err := s.repo.ListEndedPlayerResults(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return buildEventRatings(rows), nil
}

func buildEventRatings(rows []models.EventPlayerResultRow) []models.EventRatingEntry {
	byPlayer := make(map[int]*playerAgg)
	for _, row := range rows {
		agg, ok := byPlayer[row.PlayerID]
		if !ok {
			agg = &playerAgg{
				playerID:  row.PlayerID,
				name:      row.Name,
				sessionGP: make([]int, 0),
			}
			byPlayer[row.PlayerID] = agg
		}
		agg.games++
		agg.tp += row.TP
		agg.gp += row.GP
		agg.sessionGP = append(agg.sessionGP, row.GP)
	}

	players := make([]*playerAgg, 0, len(byPlayer))
	for _, agg := range byPlayer {
		sort.Sort(sort.Reverse(sort.IntSlice(agg.sessionGP)))
		players = append(players, agg)
	}

	sort.SliceStable(players, func(i, j int) bool {
		a, b := players[i], players[j]
		if a.tp != b.tp {
			return a.tp > b.tp
		}

		if a.gp != b.gp {
			return a.gp > b.gp
		}

		if cmp := compareSessionGPs(a.sessionGP, b.sessionGP); cmp != 0 {
			return cmp > 0
		}

		// Deterministic lottery: lower player_id ranks higher.
		return a.playerID < b.playerID
	})

	out := make([]models.EventRatingEntry, len(players))
	for i, p := range players {
		out[i] = models.EventRatingEntry{
			Rank:     i + 1,
			PlayerID: p.playerID,
			Name:     p.name,
			Games:    p.games,
			TP:       p.tp,
			GP:       p.gp,
		}
	}
	return out
}

// compareSessionGPs compares descending GP lists lexicographically.
// Positive means a is better than b.
func compareSessionGPs(a, b []int) int {
	n := min(len(b), len(a))

	for i := range n {
		if a[i] != b[i] {
			return a[i] - b[i]
		}
	}

	return len(a) - len(b)
}
