package event

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

// TestHandler_SetMemberRoles checks that the roles reach the service for the
// member of the URL.
func TestHandler_SetMemberRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	eventSvc := NewMockEventService(t)
	eventSvc.EXPECT().HasCapability(mock.Anything, 1, 1, 0, models.CapManageMembers).Return(nil)
	eventSvc.EXPECT().SetMemberRoles(mock.Anything, 1, 20, []models.EventRole{models.EventReferee}).Return(models.ErrPlayerRoleInGame)
	c, w := newEntryTestContext("/members/20", `{"roles":["referee"]}`)
	c.Request.Method = http.MethodPut
	c.Params = append(c.Params, gin.Param{Key: "playerId", Value: "20"})

	NewHandler(eventSvc, NewMockGameService(t)).SetMemberRoles(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}
