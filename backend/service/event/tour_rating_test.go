package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestTourRating(t *testing.T) {
	// Tours 1 and 2 finished, tour 4 live; Dan withdrew after tour 1, Eve
	// took a seat in tour 2 (as a new player would), Fay hasn't played.
	results := []models.TourResultRow{
		{Tour: 1, PlayerID: 1, Name: "Ann", TP: 4, GP: 30, Table: 1},
		{Tour: 1, PlayerID: 2, Name: "Bob", TP: 2, GP: 10, Table: 1},
		{Tour: 1, PlayerID: 3, Name: "Cid", TP: 1, GP: -10, Table: 1},
		{Tour: 1, PlayerID: 4, Name: "Dan", TP: 0, GP: -30, Table: 1},
		{Tour: 2, PlayerID: 1, Name: "Ann", TP: 0, GP: -20, Table: 1},
		{Tour: 2, PlayerID: 2, Name: "Bob", TP: 2, GP: 0, Table: 1},
		{Tour: 2, PlayerID: 3, Name: "Cid", TP: 1, GP: 0, Table: 1},
		{Tour: 2, PlayerID: 5, Name: "Eve", TP: 4, GP: 20, Table: 1},
	}
	participants := models.TournamentParticipants{Drawn: true, Participants: []models.TournamentParticipant{
		{UserID: 1, Name: "Ann", Number: new(3)},
		{UserID: 2, Name: "Bob", Number: new(1)},
		{UserID: 3, Name: "Cid", Number: new(2)},
		{UserID: 4, Name: "Dan", Number: new(4), Withdrawn: true},
		{UserID: 5, Name: "Eve", Number: new(5)},
		{UserID: 6, Name: "Fay", Number: new(6)},
	}}

	rating := tourRating(results, 4, participants, models.DefaultTournamentSettings())
	assert.Equal(t, []int{1, 2}, rating.Tours)
	assert.Equal(t, 4, rating.InProgress)

	names := []string{}
	ranks := []any{}
	for _, row := range rating.Rows {
		names = append(names, row.Name)
		if row.Rank == nil {
			ranks = append(ranks, nil)
		} else {
			ranks = append(ranks, *row.Rank)
		}
	}
	// Bob 4 TP / +10, Ann 4 / +10 (equal: shared place, by number), Eve 4 /
	// +20 first; Cid 2; Fay hasn't played; Dan withdrew.
	assert.Equal(t, []string{"Eve", "Bob", "Ann", "Cid", "Fay", "Dan"}, names)
	assert.Equal(t, []any{1, 2, 2, 4, nil, nil}, ranks)

	ann := rating.Rows[2]
	require.Len(t, ann.Tours, 2)
	assert.Equal(t, &models.TourResult{TP: 4, GP: 30, Table: 1}, ann.Tours[0])
	assert.Equal(t, &models.TourResult{TP: 0, GP: -20, Table: 1}, ann.Tours[1])
	assert.Equal(t, new(3), ann.Number)
	assert.Nil(t, rating.Rows[0].Tours[0], "Eve didn't play tour 1")
	assert.True(t, rating.Rows[5].Withdrawn)
	assert.Equal(t, -30, rating.Rows[5].GP, "Dan keeps his results")

	// Before the draw, no numbers.
	participants.Drawn = false
	rating = tourRating(nil, 0, participants, models.DefaultTournamentSettings())
	assert.Empty(t, rating.Tours)
	assert.Len(t, rating.Rows, 6)
	for _, row := range rating.Rows {
		assert.Nil(t, row.Number)
		if !row.Withdrawn {
			assert.Nil(t, row.Rank, "no place before a finished tour")
		}
	}
}

func TestTourRating_tieBreak(t *testing.T) {
	// Ann and Bob end equal: 4 TP, +10 GP. Head to head they are 1:1 (Ann
	// ahead in tour 1, Bob in tour 2); Ann has a first place, Bob none; Bob
	// has the lower draw number.
	results := []models.TourResultRow{
		{Tour: 1, PlayerID: 1, Name: "Ann", TP: 4, GP: 30, Table: 1},
		{Tour: 1, PlayerID: 2, Name: "Bob", TP: 2, GP: 10, Table: 1},
		{Tour: 1, PlayerID: 3, Name: "Cid", TP: 1, GP: -10, Table: 1},
		{Tour: 1, PlayerID: 4, Name: "Dan", TP: 0, GP: -30, Table: 1},
		{Tour: 2, PlayerID: 1, Name: "Ann", TP: 0, GP: -20, Table: 1},
		{Tour: 2, PlayerID: 2, Name: "Bob", TP: 2, GP: 0, Table: 1},
		{Tour: 2, PlayerID: 3, Name: "Cid", TP: 1, GP: 0, Table: 1},
		{Tour: 2, PlayerID: 4, Name: "Dan", TP: 4, GP: 20, Table: 1},
	}
	participants := models.TournamentParticipants{Drawn: true, Participants: []models.TournamentParticipant{
		{UserID: 1, Name: "Ann", Number: new(3)},
		{UserID: 2, Name: "Bob", Number: new(1)},
		{UserID: 3, Name: "Cid", Number: new(2)},
		{UserID: 4, Name: "Dan", Number: new(4)},
	}}
	ranks := func(settings models.TournamentSettings) map[string]int {
		out := map[string]int{}
		for _, row := range tourRating(results, 0, participants, settings).Rows {
			out[row.Name] = *row.Rank
		}
		return out
	}
	breaking := func(criteria ...string) models.TournamentSettings {
		s := models.DefaultTournamentSettings()
		s.Ties, s.TieBreakers = models.TiesBreak, criteria
		return s
	}

	assert.Equal(t, map[string]int{"Ann": 1, "Bob": 1, "Dan": 3, "Cid": 4},
		ranks(models.DefaultTournamentSettings()), "shared")
	assert.Equal(t, map[string]int{"Bob": 1, "Ann": 2, "Dan": 3, "Cid": 4},
		ranks(breaking()), "the draw number only")
	assert.Equal(t, map[string]int{"Bob": 1, "Ann": 2, "Dan": 3, "Cid": 4},
		ranks(breaking(models.TieHeadToHead)), "1:1, then the draw number")
	assert.Equal(t, map[string]int{"Ann": 1, "Bob": 2, "Dan": 3, "Cid": 4},
		ranks(breaking(models.TieHeadToHead, models.TieFirstPlaces)), "first places decide")
}

func TestTournamentSettings_ties(t *testing.T) {
	valid := models.DefaultTournamentSettings()
	require.NoError(t, valid.Validate())
	for name, change := range map[string]func(*models.TournamentSettings){
		"unknown ties":        func(s *models.TournamentSettings) { s.Ties = "coin" },
		"criteria with share": func(s *models.TournamentSettings) { s.TieBreakers = []string{models.TieHeadToHead} },
		"unknown criterion":   func(s *models.TournamentSettings) { s.Ties, s.TieBreakers = models.TiesBreak, []string{"age"} },
		"a criterion twice": func(s *models.TournamentSettings) {
			s.Ties, s.TieBreakers = models.TiesBreak, []string{models.TieFirstPlaces, models.TieFirstPlaces}
		},
	} {
		s := models.DefaultTournamentSettings()
		change(&s)
		assert.ErrorIs(t, s.Validate(), models.ErrInvalidTournamentSettings, name)
	}
}
