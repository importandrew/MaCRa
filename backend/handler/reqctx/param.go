package reqctx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey  = "userID"
	ZeroUserID = 0 // unauthorized user, value must be less than ID of any valid user
)

// use only inside API scope (where UserIDKey value is assigned), fallback to ZeroUserID
func GetUserID(c *gin.Context) int {
	value, ok := c.Get(UserIDKey)
	if !ok {
		return ZeroUserID
	}

	id, ok := value.(int)
	if !ok {
		return ZeroUserID
	}

	return id
}

func GetGameID(c *gin.Context) int {
	gameID, _ := strconv.Atoi(c.Param("gameId"))
	return gameID
}

func GetEventID(c *gin.Context) int {
	eventID, _ := strconv.Atoi(c.Param("eventId"))
	return eventID
}

func GetPlayerID(c *gin.Context) int {
	playerID, _ := strconv.Atoi(c.Param("playerId"))
	return playerID
}
