package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func testEvent(typ models.EventType) models.Event {
	return models.Event{
		ID:          1,
		Type:        typ,
		OwnerUserID: 10,
		Visibility:  models.EventVisibilityPrivate,
		Settings:    models.DefaultEventSettings(),
	}
}

func TestHasCapability_visibility(t *testing.T) {
	// One setting: a public event shows everything but managing it to
	// everyone, a private one to its members; staff always (area 1).
	player := []models.EventRole{models.EventPlayer}
	observer := []models.EventRole{models.EventObserver}
	for name, testCase := range map[string]struct {
		visibility models.EventVisibility
		caller     int
		roles      []models.EventRole
		want       bool
	}{
		"public, anonymous":    {models.EventVisibilityPublic, 0, nil, true},
		"public, stranger":     {models.EventVisibilityPublic, 99, nil, true},
		"public, player":       {models.EventVisibilityPublic, 5, player, true},
		"private, anonymous":   {models.EventVisibilityPrivate, 0, nil, false},
		"private, stranger":    {models.EventVisibilityPrivate, 99, nil, false},
		"private, player":      {models.EventVisibilityPrivate, 5, player, true},
		"private, observer":    {models.EventVisibilityPrivate, 6, observer, true},
		"private, owner alone": {models.EventVisibilityPrivate, 10, nil, true},
	} {
		for _, typ := range []models.EventType{models.EventRegular, models.EventTournament} {
			t.Run(name+" "+string(typ), func(t *testing.T) {
				event := testEvent(typ)
				event.Visibility = testCase.visibility
				authCtx := eventAuthContext(event, testCase.caller, testCase.roles)
				for _, capability := range []models.Capability{models.CapViewRatings, models.CapWatchLive, models.CapWatchResults} {
					assert.Equal(t, testCase.want, authCtx.hasCapability(capability), capability)
				}
			})
		}
	}
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
		Visibility: models.EventVisibilityPrivate}
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
	// A public tournament: still, a tour's other tables are hidden while it
	// is played, and a drafted seating from all but the staff (area 8,
	// "Visibility").
	event := testEvent(models.EventTournament)
	event.Visibility = models.EventVisibilityPublic
	draft, live, published, finished := models.TourDraft, models.TourLive, models.TourPublished, models.TourFinished
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
		"another table, seating published": {game: models.GameSummary{Status: models.GameCreated, TourStatus: &published}, want: true},
		"own table, seating drafted":       {caller: 5, roles: player, seated: true, game: models.GameSummary{Status: models.GameCreated, TourStatus: &draft}},
		"referee, seating drafted":         {caller: 7, roles: []models.EventRole{models.EventReferee}, game: models.GameSummary{Status: models.GameCreated, TourStatus: &draft}, want: true},
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
