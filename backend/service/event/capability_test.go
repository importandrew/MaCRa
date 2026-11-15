package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func testEvent(typ models.EventType) models.Event {
	return models.Event{
		ID:           1,
		Type:         typ,
		OwnerUserID:  10,
		Visibility:   models.EventVisibilityPrivate,
		RatingsReach: models.ContentReachPublic,
		GamesReach:   models.ContentReachPublic,
		Settings:     models.DefaultEventSettings(),
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

// The owner sees everything in their event, as staff, without any role.
func TestHasCapability_ownerSeesStaffOnlyRatings(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.RatingsReach = models.ContentReachStaff

	authCtx := EventAuthContext{
		Event:    event,
		CallerID: 10,
		IsOwner:  true,
		IsMember: true,
	}

	assert.True(t, authCtx.hasCapability(models.CapViewRatings))
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
		Game:     &models.GameSummary{Status: models.GameLive, LogSize: 3},
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
		Game:           &models.GameSummary{Status: models.GameLive},
	}
	assert.True(t, player.hasCapability(models.CapFinishGame))

	referee := EventAuthContext{
		Event:    event,
		CallerID: 6,
		Roles:    []models.EventRole{models.EventReferee},
		IsMember: true,
		Game:     &models.GameSummary{Status: models.GameLive},
	}
	assert.True(t, referee.hasCapability(models.CapFinishGame))
}

func TestHasCapability_seatedPlayerWatchesOwnGame(t *testing.T) {
	event := testEvent(models.EventTournament)
	event.GamesReach = models.ContentReachStaff

	seated := EventAuthContext{
		Event:          event,
		CallerID:       5,
		Roles:          []models.EventRole{models.EventPlayer},
		IsMember:       true,
		IsSeatedInGame: true,
	}
	assert.True(t, seated.hasCapability(models.CapWatchLive))
	assert.True(t, seated.hasCapability(models.CapWatchResults))

	notSeated := seated
	notSeated.IsSeatedInGame = false
	assert.False(t, notSeated.hasCapability(models.CapWatchLive))
	assert.False(t, notSeated.hasCapability(models.CapWatchResults))
}

func TestHasCapability_seatedPlayerWithoutEventAccess(t *testing.T) {
	// A removed member of a private event loses the game too.
	authCtx := EventAuthContext{
		Event:          testEvent(models.EventRegular),
		CallerID:       5,
		IsSeatedInGame: true,
	}

	assert.False(t, authCtx.hasCapability(models.CapWatchLive))
	assert.False(t, authCtx.hasCapability(models.CapWatchResults))
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

func TestSeatedAbortOnlyWithoutRounds(t *testing.T) {
	for abort, want := range map[string]bool{
		models.AbortByReferees:               false,
		models.AbortBySeatedBeforeFirstRound: true,
		models.AbortBySeated:                 false,
	} {
		event := testEvent(models.EventRegular)
		event.Settings.AbortGame = abort
		assert.Equal(t, want, EventAuthContext{Event: event}.seatedAbortOnlyWithoutRounds(), abort)
	}
}

func TestHasCapability_tournamentDesk(t *testing.T) {
	tournament := models.Event{ID: 1, Type: models.EventTournament, OwnerUserID: 10,
		Visibility: models.EventVisibilityPrivate, RatingsReach: models.ContentReachMembers, GamesReach: models.ContentReachMembers}
	regular := tournament
	regular.Type = models.EventRegular

	testData := map[string]struct {
		event  models.Event
		caller int
		roles  []models.EventRole
		view   bool
		run    bool
	}{
		"referee":         {event: tournament, caller: 5, roles: []models.EventRole{models.EventReferee}, view: true, run: true},
		"observer":        {event: tournament, caller: 5, roles: []models.EventRole{models.EventObserver}, view: true},
		"owner, no roles": {event: tournament, caller: 10, view: true},
		"player":          {event: tournament, caller: 5, roles: []models.EventRole{models.EventPlayer}},
		"regular referee": {event: regular, caller: 5, roles: []models.EventRole{models.EventReferee}},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			authCtx := eventAuthContext(testCase.event, testCase.caller, testCase.roles)
			assert.Equal(t, testCase.view, authCtx.hasCapability(models.CapViewTournament))
			assert.Equal(t, testCase.run, authCtx.hasCapability(models.CapRunTournament))
		})
	}
}

func TestCanUndoInGame_tournament(t *testing.T) {
	live, finished := models.TourLive, models.TourFinished
	long := time.Now().Add(-24 * time.Hour)
	for name, testCase := range map[string]struct {
		game models.GameSummary
		want bool
	}{
		"live game, live tour":        {models.GameSummary{Status: models.GameLive, TourStatus: &live}, true},
		"ended long ago, live tour":   {models.GameSummary{Status: models.GameEnded, FinishedAt: &long, TourStatus: &live}, true},
		"ended, tour finished":        {models.GameSummary{Status: models.GameEnded, FinishedAt: &long, TourStatus: &finished}, false},
		"regular game ended long ago": {models.GameSummary{Status: models.GameEnded, FinishedAt: &long}, false},
	} {
		t.Run(name, func(t *testing.T) {
			authCtx := EventAuthContext{Game: &testCase.game, Now: time.Now()}
			assert.Equal(t, testCase.want, authCtx.canUndoInGame())
		})
	}
}

func TestCanFinishGame_tournament(t *testing.T) {
	tournament := testEvent(models.EventTournament)
	live := models.TourLive
	for name, testCase := range map[string]struct {
		stopped bool
		seated  bool
		roles   []models.EventRole
		want    bool
	}{
		"seated, time not up":     {seated: true, roles: []models.EventRole{models.EventPlayer}},
		"seated, stopped at time": {stopped: true, seated: true, roles: []models.EventRole{models.EventPlayer}, want: true},
		"referee, stopped":        {stopped: true, roles: []models.EventRole{models.EventReferee}, want: true},
		"another player, stopped": {stopped: true, roles: []models.EventRole{models.EventPlayer}},
	} {
		t.Run(name, func(t *testing.T) {
			authCtx := eventAuthContext(tournament, 5, testCase.roles)
			authCtx.IsSeatedInGame = testCase.seated
			authCtx.Game = &models.GameSummary{Status: models.GameLive, TourStatus: &live, StoppedAtTime: testCase.stopped}
			assert.Equal(t, testCase.want, authCtx.canFinishGame())
		})
	}
}

func TestHasCapability_tourTablesHidden(t *testing.T) {
	// A public tournament whose games are public: still, a tour's other
	// tables are hidden while they play (area 8, "Visibility").
	event := testEvent(models.EventTournament)
	event.Visibility = models.EventVisibilityPublic
	live, published, finished := models.TourLive, models.TourPublished, models.TourFinished
	player := []models.EventRole{models.EventPlayer}
	for name, testCase := range map[string]struct {
		caller   int
		roles    []models.EventRole
		seated   bool
		finished bool
		game     models.GameSummary
		want     bool
	}{
		"anonymous, live table":            {game: models.GameSummary{Status: models.GameLive, TourStatus: &live}},
		"player still playing":             {caller: 5, roles: player, game: models.GameSummary{Status: models.GameLive, TourStatus: &live}},
		"player whose game ended":          {caller: 5, roles: player, finished: true, game: models.GameSummary{Status: models.GameLive, TourStatus: &live}, want: true},
		"the table itself":                 {caller: 5, roles: player, seated: true, game: models.GameSummary{Status: models.GameLive, TourStatus: &live}, want: true},
		"observer":                         {caller: 6, roles: []models.EventRole{models.EventObserver}, game: models.GameSummary{Status: models.GameLive, TourStatus: &live}, want: true},
		"owner":                            {caller: 10, game: models.GameSummary{Status: models.GameLive, TourStatus: &live}, want: true},
		"ended table of the live tour":     {caller: 5, roles: player, game: models.GameSummary{Status: models.GameEnded, TourStatus: &live}, want: true},
		"another table, seating published": {caller: 5, roles: player, game: models.GameSummary{Status: models.GameCreated, TourStatus: &published}},
		"finished tour":                    {game: models.GameSummary{Status: models.GameEnded, TourStatus: &finished}, want: true},
	} {
		t.Run(name, func(t *testing.T) {
			authCtx := eventAuthContext(event, testCase.caller, testCase.roles)
			authCtx.IsSeatedInGame = testCase.seated
			authCtx.FinishedInTour = testCase.finished
			game := testCase.game
			authCtx.Game = &game
			for _, capability := range []models.Capability{models.CapViewGame, models.CapWatchLive, models.CapWatchResults} {
				assert.Equal(t, testCase.want, authCtx.hasCapability(capability), capability)
			}
		})
	}
}
