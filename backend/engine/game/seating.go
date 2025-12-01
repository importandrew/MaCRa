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

func (g *Game) CurrentWind(id int) models.Wind {
	return g.currentSeating.Wind(id)
}

func (g *Game) StartWind(id int) models.Wind {
	return g.startSeating.Wind(id)
}
