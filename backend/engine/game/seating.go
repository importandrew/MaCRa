package game

import (
	"slices"

	"github.com/importandrew/MaCRa/models"
)

type Seating [4]int

// E <-> S ; W <-> N.
func (s *Seating) ApplySouthNorthSeating() {
	s[models.East], s[models.South] =
		s[models.South], s[models.East]
	s[models.West], s[models.North] =
		s[models.North], s[models.West]
}

// E -> W -> S -> N -> E.
func (s *Seating) ApplyWestSeating() {
	s[models.East], s[models.South],
		s[models.West], s[models.North] =
		s[models.North], s[models.West],
		s[models.East], s[models.South]
}

// Usual dealership transfer: E -> S -> W -> N -> E.
func (s *Seating) NextRoundSeating() {
	s[models.East], s[models.South],
		s[models.West], s[models.North] =
		s[models.South], s[models.West],
		s[models.North], s[models.East]
}

func (s *Seating) Wind(id int) models.Wind {
	return models.Wind(slices.IndexFunc(s[:], func(i int) bool {
		return i == id
	}))
}

func (s *Seating) Player(wind models.Wind) int {
	return s[wind]
}

func (g *Game) CurrentWind(id int) models.Wind {
	return g.currentSeating.Wind(id)
}

func (g *Game) StartWind(id int) models.Wind {
	return g.startSeating.Wind(id)
}

func (g *Game) RoundStartWind(id int) models.Wind {
	return g.roundSeating.Wind(id)
}

func (g *Game) PlayerByStartingSeat(seat models.Wind) int {
	return g.startSeating[seat]
}

func (g *Game) PlayerBySeatInRound(round int, seat models.Wind) int {
	return g.startSeating.Player(models.RoundSeats[round][seat])
}

func (g *Game) PlayerByWindInRound(round int, wind models.Wind) int {
	return g.startSeating.Player(models.RoundWinds[round][wind])
}
