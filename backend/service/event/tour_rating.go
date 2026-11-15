package event

import (
	"cmp"
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// TourRating returns the tournament's rating by tour: TP and GP of each
// finished tour and the totals. The active players come first by TP, then
// GP (players equal in both share the place), those who haven't played a
// finished tour yet without a place; the players who withdrew follow,
// without a place.
func (s *Service) TourRating(ctx context.Context, eventID int) (models.TourRating, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.TourRating{}, err
	}
	if event.Type != models.EventTournament {
		return models.TourRating{}, models.ErrNotTournament
	}
	results, inProgress, err := s.repo.TourResults(ctx, eventID)
	if err != nil {
		return models.TourRating{}, err
	}
	participants, err := s.repo.ListParticipants(ctx, eventID)
	if err != nil {
		return models.TourRating{}, err
	}
	return tourRating(results, inProgress, participants), nil
}

// tourRating builds the rating by tour from the results and the
// participants, those without results too.
func tourRating(
	results []models.TourResultRow, inProgress int,
	participants models.TournamentParticipants,
) models.TourRating {
	rating := models.TourRating{Tours: []int{}, InProgress: inProgress, Rows: []models.TourRatingRow{}}
	tourIndex := map[int]int{}
	for _, r := range results {
		if _, ok := tourIndex[r.Tour]; !ok {
			tourIndex[r.Tour] = len(rating.Tours)
			rating.Tours = append(rating.Tours, r.Tour)
		}
	}

	rows := map[int]*models.TourRatingRow{}
	order := []int{} // player IDs, first seen
	row := func(playerID int, name string) *models.TourRatingRow {
		if r, ok := rows[playerID]; ok {
			return r
		}
		r := &models.TourRatingRow{PlayerID: playerID, Name: name, Tours: make([]*models.TourResult, len(rating.Tours))}
		rows[playerID] = r
		order = append(order, playerID)
		return r
	}
	for _, res := range results {
		r := row(res.PlayerID, res.Name)
		r.TP += res.TP
		r.GP += res.GP
		r.Tours[tourIndex[res.Tour]] = &models.TourResult{TP: res.TP, GP: res.GP, Table: res.Table}
	}
	for _, p := range participants.Participants {
		r := row(p.UserID, p.Name)
		r.Withdrawn = p.Withdrawn
		if participants.Drawn {
			r.Number = p.Number
		}
	}

	list := make([]models.TourRatingRow, 0, len(order))
	for _, playerID := range order {
		list = append(list, *rows[playerID])
	}
	rankRows(list)
	rating.Rows = list
	return rating
}

// The groups of the rating by tour, in order.
const (
	groupRanked = iota
	groupNotPlayed
	groupWithdrew
)

// ratingGroup returns a row's group: active players who played a finished
// tour, those who haven't yet, and those who withdrew.
func ratingGroup(r models.TourRatingRow) int {
	switch {
	case r.Withdrawn:
		return groupWithdrew
	case !played(r):
		return groupNotPlayed
	default:
		return groupRanked
	}
}

// rankRows orders the rows by group, then TP and GP (within a group the
// number, then the name, keeps the order stable), and gives the ranked group
// its places: players equal in TP and GP share one.
func rankRows(list []models.TourRatingRow) {
	slices.SortStableFunc(list, func(a, b models.TourRatingRow) int {
		return cmp.Or(
			cmp.Compare(ratingGroup(a), ratingGroup(b)),
			cmp.Compare(b.TP, a.TP),
			cmp.Compare(b.GP, a.GP),
			cmp.Compare(numberOr(a), numberOr(b)),
			cmp.Compare(a.Name, b.Name),
		)
	})
	for i := range list {
		if ratingGroup(list[i]) != groupRanked {
			continue
		}
		rank := i + 1
		if i > 0 && list[i-1].Rank != nil && list[i-1].TP == list[i].TP && list[i-1].GP == list[i].GP {
			rank = *list[i-1].Rank
		}
		list[i].Rank = &rank
	}
}

// played reports whether a player has a result in a finished tour.
func played(r models.TourRatingRow) bool {
	for _, res := range r.Tours {
		if res != nil {
			return true
		}
	}
	return false
}

// numberOr returns a player's participant number, or a large one without it
// (they sort last).
func numberOr(r models.TourRatingRow) int {
	if r.Number == nil {
		return int(^uint(0) >> 1)
	}
	return *r.Number
}
