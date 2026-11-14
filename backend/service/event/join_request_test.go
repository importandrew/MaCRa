package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestMyCapabilities_joinRequest(t *testing.T) {
	limit := 4
	finished := time.Now()

	testData := map[string]struct {
		visibility models.EventVisibility
		accept     bool
		limit      *int
		finishedAt *time.Time
		callerID   int
		roles      []models.EventRole
		status     models.JoinRequestStatus
		players    int
		want       models.JoinRequestState
	}{
		"available":          {accept: true, callerID: 5, want: models.JoinStateAvailable},
		"visitor":            {accept: true, want: models.JoinStateAvailable},
		"member":             {accept: true, callerID: 5, roles: []models.EventRole{models.EventObserver}, want: ""},
		"owner":              {accept: true, callerID: 10, want: ""},
		"switched off":       {callerID: 5, want: ""},
		"private":            {visibility: models.EventVisibilityPrivate, accept: true, callerID: 5, want: ""},
		"finished":           {accept: true, finishedAt: &finished, callerID: 5, want: ""},
		"pending":            {accept: true, callerID: 5, status: models.JoinRequestPending, want: models.JoinStatePending},
		"pending, now off":   {callerID: 5, status: models.JoinRequestPending, want: models.JoinStatePending},
		"declined":           {accept: true, callerID: 5, status: models.JoinRequestDeclined, want: models.JoinStateDeclined},
		"declined, now off":  {callerID: 5, status: models.JoinRequestDeclined, want: ""},
		"full":               {accept: true, limit: &limit, callerID: 5, players: 4, want: models.JoinStateFull},
		"one place left":     {accept: true, limit: &limit, callerID: 5, players: 3, want: models.JoinStateAvailable},
		"visitor, full":      {accept: true, limit: &limit, players: 4, want: models.JoinStateFull},
		"declined over full": {accept: true, limit: &limit, callerID: 5, status: models.JoinRequestDeclined, want: models.JoinStateDeclined},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			visibility := testCase.visibility
			if visibility == "" {
				visibility = models.EventVisibilityPublic
			}
			event := models.Event{
				ID: 1, Type: models.EventRegular, OwnerUserID: 10, Visibility: visibility,
				RatingsReach: models.ContentReachPublic, GamesReach: models.ContentReachPublic,
				AcceptRequests: testCase.accept, PlayerLimit: testCase.limit, FinishedAt: testCase.finishedAt,
			}

			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, 1).Return(event, nil)
			if testCase.callerID != 0 {
				repo.EXPECT().GetMemberRoles(mock.Anything, 1, testCase.callerID).Return(testCase.roles, nil)
				repo.EXPECT().JoinRequestStatus(mock.Anything, 1, testCase.callerID).Return(testCase.status, nil).Maybe()
			}
			repo.EXPECT().CountPlayers(mock.Anything, 1).Return(testCase.players, nil).Maybe()
			repo.EXPECT().CountJoinRequests(mock.Anything, 1).Return(0, nil).Maybe()

			me, err := NewService(repo, nil, nil).MyCapabilities(context.Background(), 1, testCase.callerID)

			require.NoError(t, err)
			assert.Equal(t, testCase.want, me.JoinRequest)
		})
	}
}

func TestMyCapabilities_pendingRequests(t *testing.T) {
	finished := time.Now()
	testData := map[string]struct {
		callerID   int
		finishedAt *time.Time
		want       int
	}{
		"owner":          {callerID: 10, want: 3},
		"not the owner":  {callerID: 5, want: 0},
		"finished event": {callerID: 10, finishedAt: &finished, want: 0},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			event := models.Event{
				ID: 1, Type: models.EventRegular, OwnerUserID: 10, Visibility: models.EventVisibilityPublic,
				RatingsReach: models.ContentReachPublic, GamesReach: models.ContentReachPublic,
				FinishedAt: testCase.finishedAt,
			}

			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, 1).Return(event, nil)
			repo.EXPECT().GetMemberRoles(mock.Anything, 1, testCase.callerID).
				Return([]models.EventRole{models.EventPlayer}, nil)
			if testCase.want > 0 {
				repo.EXPECT().CountJoinRequests(mock.Anything, 1).Return(testCase.want, nil)
			}

			me, err := NewService(repo, nil, nil).MyCapabilities(context.Background(), 1, testCase.callerID)

			require.NoError(t, err)
			assert.Equal(t, testCase.want, me.PendingRequests)
		})
	}
}

func TestUpdateJoining_invalidLimit(t *testing.T) {
	zero := 0
	err := NewService(NewMockEventRepo(t), nil, nil).
		UpdateJoining(context.Background(), 1, models.JoiningSettings{PlayerLimit: &zero})

	assert.ErrorIs(t, err, models.ErrInvalidPlayerLimit)
}
