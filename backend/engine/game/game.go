package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

type Game struct {
	id             int
	startSeating   Seating
	currentSeating Seating
	status         models.GameStatus
	round          int
	rules          Rules
	score          models.GameScore
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
		status:         models.GameLive,
		round:          1,
		rules:          cfg.Rules,
		score:          NewScore(cfg.Players),
	}

	log.Infof("playing: %v", cfg.Players)
	return &game
}

func (g *Game) End() error {
	if g.status != models.GameLive {
		return fmt.Errorf("GameID %d: game should be live to be able to end", g.id)
	}
	g.status = models.GameEnded
	g.score.CalculateTP()
	logger.WithField("final result", g.score.Final()).Info("game ended")

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

func (g *Game) Score() models.GameScore {
	return g.score
}
