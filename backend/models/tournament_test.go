package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTournamentSettings_Validate(t *testing.T) {
	ok := DefaultTournamentSettings()
	assert.NoError(t, ok.Validate())

	extra := ok
	extra.TimeUp, extra.ExtraHands = TimeUpExtraHands, 2
	assert.NoError(t, extra.Validate())

	for name, change := range map[string]func(s *TournamentSettings){
		"too few tours":           func(s *TournamentSettings) { s.Tours = MinTours - 1 },
		"too many tours":          func(s *TournamentSettings) { s.Tours = MaxTours + 1 },
		"too short":               func(s *TournamentSettings) { s.TourMinutes = MinTourMinutes - 1 },
		"too long":                func(s *TournamentSettings) { s.TourMinutes = MaxTourMinutes + 1 },
		"unknown time-up rule":    func(s *TournamentSettings) { s.TimeUp = "never" },
		"extra hands missing":     func(s *TournamentSettings) { s.TimeUp = TimeUpExtraHands },
		"too many extra hands":    func(s *TournamentSettings) { s.TimeUp, s.ExtraHands = TimeUpExtraHands, MaxExtraHands+1 },
		"extra hands without use": func(s *TournamentSettings) { s.ExtraHands = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			settings := DefaultTournamentSettings()
			change(&settings)
			assert.ErrorIs(t, settings.Validate(), ErrInvalidTournamentSettings)
		})
	}
}

func TestEvent_Status(t *testing.T) {
	now := time.Now()
	assert.Equal(t, EventPreparation, Event{Type: EventTournament}.Status())
	assert.Equal(t, EventInProgress, Event{Type: EventTournament, StartedAt: &now}.Status())
	assert.Equal(t, EventFinished, Event{Type: EventTournament, StartedAt: &now, FinishedAt: &now}.Status())
	assert.Equal(t, EventInProgress, Event{Type: EventRegular}.Status(), "a regular event needs no start")
	assert.Equal(t, EventFinished, Event{Type: EventRegular, FinishedAt: &now}.Status())
}
