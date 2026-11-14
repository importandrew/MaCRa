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
	rules          models.GameRules
	score          models.GameScore
	logSize        int // round log entries applied: results and penalties
	// lastHand is the hand after which the game ends: rules.Hands, or fewer
	// once a tournament tour's time is up (LimitHands); 0 means rules.Hands.
	lastHand int
	// stoppedAtTime: the tour's time ran out with the rule "the game ends at
	// once" (StopAtTime); no more entries.
	stoppedAtTime bool
}

type Config struct {
	ID      int
	Players [4]int
	Rules   models.GameRules
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

// Rules returns the rules the game is played by.
func (g *Game) Rules() models.GameRules {
	return g.rules
}

func (g *Game) Score() models.GameScore {
	return g.score
}

// LastHand returns the hand after which the game ends.
func (g *Game) LastHand() int {
	if g.lastHand == 0 || g.lastHand > g.rules.Hands {
		return g.rules.Hands
	}
	return g.lastHand
}

// LimitHands makes the game end after hand last at the latest (when a
// tournament tour's time is up: the hand in progress, and maybe a few more).
// It never moves the end later than the rules' number of hands, nor before
// the hand in progress.
func (g *Game) LimitHands(last int) {
	last = max(last, g.round)
	if last < g.LastHand() {
		g.lastHand = last
	}
}

// StopAtTime ends a live game at once because the tour's time is up with the
// rule "the game ends at once": the hand in progress doesn't count, and later
// entries are refused with models.ErrTimeUp.
func (g *Game) StopAtTime() error {
	if !g.IsLive() {
		return nil
	}
	g.stoppedAtTime = true
	return g.End()
}

// StoppedAtTime reports whether StopAtTime ended the game.
func (g *Game) StoppedAtTime() bool {
	return g.stoppedAtTime
}

// checkLive returns nil for a live game, models.ErrTimeUp for one stopped
// at time and models.ErrGameNotLive otherwise.
func (g *Game) checkLive() error {
	switch {
	case g.IsLive():
		return nil
	case g.stoppedAtTime:
		return models.ErrTimeUp
	default:
		return models.ErrGameNotLive
	}
}
