package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

type Game struct {
	id             int
	startSeating   Seating // at the start of the game
	currentSeating Seating // at the moment
	roundSeating   Seating // at the round start (after 4, 8, 12 reseats). dont confuse with 'round' field
	status         models.GameStatus
	round          int
	rules          Rules
	score          models.GameScore
	logSize        int // round log entries applied: results and penalties
}

type Config struct {
	ID      int
	Players [4]int
	Rules   Rules
}

func NewGame(cfg Config) *Game {
	log := logger.WithField("id", cfg.ID)

	game := Game{
		id:             cfg.ID,
		startSeating:   cfg.Players,
		currentSeating: cfg.Players,
		roundSeating:   cfg.Players,
		status:         models.GameLive,
		round:          1,
		rules:          cfg.Rules,
		score:          NewScore(cfg.Players),
	}

	log.Debugf("playing: %v", cfg.Players)
	return &game
}

func (g *Game) End() error {
	if g.status != models.GameLive {
		return fmt.Errorf("GameID %d: game should be live to be able to end", g.id)
	}
	g.status = models.GameEnded
	g.score.CalculateTP()
	logger.WithField("final result", g.score.Final()).Debug("game ended")

	return nil
}

func (g *Game) IsLive() bool {
	return g.status == models.GameLive
}

func (g *Game) ID() int {
	return g.id
}

func (g *Game) Status() models.GameStatus {
	return g.status
}

func (g *Game) Round() int {
	return g.round
}

// LogSize returns the number of round log entries applied to the game: round
// results, penalties and false mahjongs. Clients send it back with a new entry
// so that an entry made from an outdated table is refused.
func (g *Game) LogSize() int {
	return g.logSize
}

func (g *Game) Score() models.GameScore {
	return g.score
}
