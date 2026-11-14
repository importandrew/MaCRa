package event

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

// TestHandler_PenaltyAtTable_needsCapability checks that a seated player's
// penalty is refused when the event leaves penalties to referees.
func TestHandler_PenaltyAtTable_needsCapability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	eventSvc := NewMockEventService(t)
	gameSvc := NewMockGameService(t)
	gameSvc.EXPECT().CurrentGameID(mock.Anything, 1, 1).Return(100, nil)
	eventSvc.EXPECT().HasCapability(mock.Anything, 1, 1, 100, models.CapEnterPenalty).Return(models.ErrEventForbidden)
	c, w := newEntryTestContext("/penalty", `{"loser":2,"value":10,"log_size":0}`)

	NewHandler(eventSvc, gameSvc).Penalty(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestHandler_PenaltyInGame checks that a referee's penalty goes to the game
// of the URL and answers with its live view.
func TestHandler_PenaltyInGame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	eventSvc := NewMockEventService(t)
	gameSvc := NewMockGameService(t)
	eventSvc.EXPECT().HasCapability(mock.Anything, 1, 1, 100, models.CapRefereeEntry).Return(nil)
	gameSvc.EXPECT().PenaltyInGame(mock.Anything, models.FalseMahjong, 100, 2, 30, true, 5).
		Return(models.GameView{ID: 100, Status: models.GameLive, LogSize: 6}, nil)
	c, w := newEntryTestContext("/game/100/false-mahjong", `{"loser":2,"value":30,"force_next_round":true,"log_size":5}`)
	c.Params = append(c.Params, gin.Param{Key: "gameId", Value: "100"})

	NewHandler(eventSvc, gameSvc).GameFalseMahjong(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"log_size":6`)
}
