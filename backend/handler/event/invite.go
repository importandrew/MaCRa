package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type inviteTokenResponse struct {
	// Token is null while the invite link is off.
	Token *string `json:"token"`
}

// InviteToken returns the event's invite token (GET /event/:id/invite),
// null while the link is off. It needs models.CapManageMembers.
func (h Handler) InviteToken(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	token, err := h.event.InviteToken(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, inviteTokenResponse{Token: token})
}

// NewInviteToken turns the invite link on with a new token, replacing the
// old link (POST /event/:id/invite). It needs models.CapManageMembers.
func (h Handler) NewInviteToken(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	token, err := h.event.NewInviteToken(c.Request.Context(), eventID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, inviteTokenResponse{Token: &token})
}

// TurnOffInvite turns the invite link off (DELETE /event/:id/invite). It
// needs models.CapManageMembers.
func (h Handler) TurnOffInvite(c *gin.Context) {
	eventID := reqctx.GetEventID(c)
	if err := h.event.HasCapability(c.Request.Context(), eventID, reqctx.GetUserID(c), 0, models.CapManageMembers); err != nil {
		httperr.Error(c, err)
		return
	}

	if err := h.event.TurnOffInvite(c.Request.Context(), eventID); err != nil {
		httperr.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Invite shows an invite link's event before joining (GET /invite/:token),
// also to a visitor, who signs in to join.
func (h Handler) Invite(c *gin.Context) {
	invite, err := h.event.Invite(c.Request.Context(), c.Param("token"), reqctx.GetUserID(c))
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, invite)
}

// JoinByInvite makes the signed-in user a player of the invite link's event
// (POST /invite/:token/join) and returns the event's ID.
func (h Handler) JoinByInvite(c *gin.Context) {
	eventID, err := h.event.JoinByInvite(c.Request.Context(), c.Param("token"), reqctx.GetUserID(c))
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"event_id": eventID})
}
