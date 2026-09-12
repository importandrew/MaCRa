package event

import (
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func testEvent(typ models.EventType) models.Event {
	return models.Event{
		ID:           1,
		Type:         typ,
		OwnerUserID:  10,
		Visibility:   models.EventVisibilityPrivate,
		RatingsReach: models.ContentReachPublic,
		GamesReach:   models.ContentReachPublic,
	}
}

func TestHasCapability_staffBypassRatings(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.RatingsReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventReferee},
		IsMember: true,
	}

	assert.True(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_ownerNoStaffBypassWithoutRefereeRole(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.RatingsReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 10,
		IsOwner:  true,
		IsMember: true,
	}

	assert.False(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_playerHiddenRatings(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.RatingsReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventPlayer},
		IsMember: true,
	}

	assert.False(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_regularMemberWatchLive(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.GamesReach = models.ContentReachMembers

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventPlayer},
		IsMember: true,
	}

	assert.True(t, authCtx.hasCapability(models.CapWatchLive))
}

func TestHasCapability_tournamentPlayerNoWatchLive(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.GamesReach = models.ContentReachMembers

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventPlayer},
		IsMember: true,
	}

	assert.False(t, authCtx.hasCapability(models.CapWatchLive))
}

func TestHasCapability_tournamentObserverWatchLive(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.GamesReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventObserver},
		IsMember: true,
	}

	assert.True(t, authCtx.hasCapability(models.CapWatchLive))
}

func TestHasCapability_publicAnonymousWatchLive(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.Visibility = models.EventVisibilityPublic
	event.GamesReach = models.ContentReachPublic

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 0,
	}

	assert.True(t, authCtx.hasCapability(models.CapWatchLive))
}

func TestHasCapability_publicAnonymousStaffReachBlocked(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.Visibility = models.EventVisibilityPublic
	event.GamesReach = models.ContentReachStaff
	event.RatingsReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 0,
	}

	assert.True(t, authCtx.canAccessEvent())
	assert.False(t, authCtx.hasCapability(models.CapWatchLive))
	assert.False(t, authCtx.hasCapability(models.CapWatchResults))
	assert.False(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_publicAnonymousMembersReachBlocked(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.Visibility = models.EventVisibilityPublic
	event.GamesReach = models.ContentReachMembers

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 0,
	}

	assert.False(t, authCtx.hasCapability(models.CapWatchLive))
}

func TestHasCapability_privateAnonymousBlocked(t *testing.T) {
	event := testEvent(models.EventRegular)

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 0,
	}

	assert.False(t, authCtx.canAccessEvent())
}

func TestHasCapability_privatePublicReachMemberCanViewRatings(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.RatingsReach = models.ContentReachPublic

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 5,
		Roles:    []models.EventRole{models.EventPlayer},
		IsMember: true,
	}

	assert.True(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_privatePublicReachAnonymousBlocked(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.RatingsReach = models.ContentReachPublic

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 0,
	}

	assert.False(t, authCtx.canAccessEvent())
	assert.False(t, authCtx.hasCapability(models.CapViewRatings))
}

func TestHasCapability_manageEventOwnerOnly(t *testing.T) {
	event := testEvent(models.EventRegular)

	owner := EventAuthContext{Event: event, CallerID: 10, IsOwner: true, IsMember: true}
	assert.True(t, owner.hasCapability(models.CapManageEvent))

	member := EventAuthContext{Event: event, CallerID: 5, Roles: []models.EventRole{models.EventReferee}, IsMember: true}
	assert.False(t, member.hasCapability(models.CapManageEvent))
}

func TestHasCapability_regularOwnerWithDBRoles(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.RatingsReach = models.ContentReachMembers
	event.GamesReach = models.ContentReachMembers

	owner := EventAuthContext{
		Event:    event,
		CallerID: 10,
		IsOwner:  true,
		IsMember: true,
		Roles:    []models.EventRole{models.EventPlayer, models.EventReferee},
	}

	assert.True(t, owner.hasCapability(models.CapCreateGame))
	assert.True(t, owner.hasCapability(models.CapAbortGame))
}

func TestHasCapability_ownerStrippedPlayRoles(t *testing.T) {
	event := testEvent(models.EventRegular)

	owner := EventAuthContext{
		Event:    event,
		CallerID: 10,
		IsOwner:  true,
		IsMember: true,
		Roles:    nil,
	}

	assert.True(t, owner.hasCapability(models.CapManageEvent))
	assert.False(t, owner.hasCapability(models.CapCreateGame))
	assert.False(t, owner.hasCapability(models.CapAbortGame))
}

func TestHasCapability_tournamentOwnerNoPlayRoles(t *testing.T) {
	event := testEvent(models.EventTournament)

	owner := EventAuthContext{
		Event:    event,
		CallerID: 10,
		IsOwner:  true,
		IsMember: true,
	}

	assert.True(t, owner.hasCapability(models.CapManageEvent))
	assert.False(t, owner.hasCapability(models.CapCreateGame))
	assert.False(t, owner.hasCapability(models.CapAbortGame))
}

func TestHasCapability_finishGameSeatedOrReferee(t *testing.T) {
	event := testEvent(models.EventRegular)

	player := EventAuthContext{
		Event:          event,
		CallerID:       5,
		Roles:          []models.EventRole{models.EventPlayer},
		IsMember:       true,
		IsSeatedInGame: true,
	}
	assert.True(t, player.hasCapability(models.CapFinishGame))

	referee := EventAuthContext{
		Event:    event,
		CallerID: 6,
		Roles:    []models.EventRole{models.EventReferee},
		IsMember: true,
	}
	assert.True(t, referee.hasCapability(models.CapFinishGame))
}

func TestCanViewContent_publicReachLoggedInNonMemberOnPublicEvent(t *testing.T) {
	event := testEvent(models.EventRegular)
	event.Visibility = models.EventVisibilityPublic
	event.GamesReach = models.ContentReachPublic

	stranger := EventAuthContext{
		Event:    event,
		CallerID: 99,
	}

	assert.True(t, stranger.hasCapability(models.CapWatchLive))
}
