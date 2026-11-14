package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func newTimedGame() *Game {
	return NewGame(Config{ID: 1, Players: [4]int{1, 2, 3, 4}, Rules: models.StandardGameRules()})
}

func TestLimitHands(t *testing.T) {
	g := newTimedGame()
	require.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
	require.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
	assert.Equal(t, 3, g.Round())

	g.LimitHands(g.Round()) // time is up during hand 3
	assert.Equal(t, 3, g.LastHand())
	g.LimitHands(10)
	assert.Equal(t, 3, g.LastHand(), "never later once limited")

	require.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
	assert.Equal(t, models.GameEnded, g.Status(), "the hand in progress was the last")
	assert.ErrorIs(t, g.EnterRoundResult(models.Draw, 0, 0, 0), models.ErrGameNotLive)
}

func TestLimitHands_neverBeforeTheHandInProgressNorAfterTheRules(t *testing.T) {
	g := newTimedGame()
	require.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
	g.LimitHands(1)
	assert.Equal(t, 2, g.LastHand(), "the hand in progress is played")

	g = newTimedGame()
	g.LimitHands(40)
	assert.Equal(t, 16, g.LastHand())
}

func TestStopAtTime(t *testing.T) {
	g := newTimedGame()
	require.NoError(t, g.EnterRoundResult(models.Draw, 0, 0, 0))
	require.NoError(t, g.StopAtTime())

	assert.True(t, g.StoppedAtTime())
	assert.Equal(t, models.GameEnded, g.Status())
	assert.ErrorIs(t, g.EnterRoundResult(models.Draw, 0, 0, 0), models.ErrTimeUp)
	assert.ErrorIs(t, g.Penalty(1, 10, false), models.ErrTimeUp)
	require.NoError(t, g.StopAtTime(), "stopping again changes nothing")
}
