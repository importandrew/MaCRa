package event

import (
	"context"
	"slices"
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

// Ratings returns the event's rating: the players with at least the event's
// minimum of ended games ranked, the others below without a place, in the
// same order (orderedPlayers). For members and staff, the players of the
// event who haven't played yet follow them with zero games, by name, so
// players see who is in the event; others don't get the member list
// (docs/user-scenarios.md, privacy).
func (s *Service) Ratings(ctx context.Context, eventID, callerID int) (models.EventRatings, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.EventRatings{}, err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, 0)
	if err != nil {
		return models.EventRatings{}, err
	}

	rows, err := s.repo.ListEndedPlayerResults(ctx, eventID)
	if err != nil {
		return models.EventRatings{}, err
	}

	minGames := event.Settings.MinRatingGames
	ratings := models.EventRatings{
		Ratings:      []models.EventRatingEntry{},
		BelowMinimum: []models.EventRatingEntry{},
		MinGames:     minGames,
	}
	for _, p := range orderedPlayers(rows) {
		entry := models.EventRatingEntry{PlayerID: p.playerID, Name: p.name, Games: p.games, TP: p.tp, GP: p.gp}
		if p.games < minGames {
			ratings.BelowMinimum = append(ratings.BelowMinimum, entry)
			continue
		}
		entry.Rank = len(ratings.Ratings) + 1
		ratings.Ratings = append(ratings.Ratings, entry)
	}

	if authCtx.isParticipating() || authCtx.isStaff() {
		notPlayed, err := s.playersWithoutGames(ctx, eventID, rows)
		if err != nil {
			return models.EventRatings{}, err
		}
		ratings.BelowMinimum = append(ratings.BelowMinimum, notPlayed...)
	}

	return ratings, nil
}

// playersWithoutGames returns the event's players who have no ended game
// (rows), by name, as rating rows with zero games.
func (s *Service) playersWithoutGames(
	ctx context.Context, eventID int, rows []models.EventPlayerResultRow,
) ([]models.EventRatingEntry, error) {
	members, err := s.repo.ListMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	played := make(map[int]bool, len(rows))
	for _, row := range rows {
		played[row.PlayerID] = true
	}

	out := []models.EventRatingEntry{}
	for _, member := range members {
		if played[member.UserID] || !slices.Contains(member.Roles, models.EventPlayer) {
			continue
		}
		out = append(out, models.EventRatingEntry{PlayerID: member.UserID, Name: member.Name})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })

	return out, nil
}

// orderedPlayers sums the results by player and orders the players by total
// TP, then total GP, then their best games by GP, then by player ID.
func orderedPlayers(rows []models.EventPlayerResultRow) []*playerAgg {
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

	return players
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
