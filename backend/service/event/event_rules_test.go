package event

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestCreate_validatesAndDefaultsToPrivate(t *testing.T) {
	testData := map[string]struct {
		params  models.EventCreateParams
		wantErr error
		want    models.Event
	}{
		"defaults: private, member reaches, trimmed texts": {
			params: models.EventCreateParams{Type: models.EventRegular, Name: "  Club  ", Description: " d ", OwnerID: 7},
			want: models.Event{
				Type: models.EventRegular, Name: "Club", Description: "d", OwnerUserID: 7,
				Visibility: models.EventVisibilityPrivate, RatingsReach: models.ContentReachMembers, GamesReach: models.ContentReachMembers,
				Rules: models.StandardGameRules(),
			},
		},
		"given rules": {
			params: models.EventCreateParams{
				Type: models.EventRegular, Name: "Club", OwnerID: 7,
				Rules: &models.GameRules{Hands: 8, MinHand: 8},
			},
			want: models.Event{
				Type: models.EventRegular, Name: "Club", OwnerUserID: 7,
				Visibility: models.EventVisibilityPrivate, RatingsReach: models.ContentReachMembers, GamesReach: models.ContentReachMembers,
				Rules: models.GameRules{Hands: 8, MinHand: 8},
			},
		},
		"invalid rules": {
			params:  models.EventCreateParams{Type: models.EventRegular, Name: "Club", Rules: &models.GameRules{Hands: 12, MinHand: 8}},
			wantErr: models.ErrInvalidGameRules,
		},
		"empty name after trimming": {
			params:  models.EventCreateParams{Type: models.EventRegular, Name: "   "},
			wantErr: models.ErrEmptyName,
		},
		"name of 101 characters": {
			params:  models.EventCreateParams{Type: models.EventRegular, Name: strings.Repeat("я", 101)},
			wantErr: models.ErrLongEventName,
		},
		"description of 1001 characters": {
			params:  models.EventCreateParams{Type: models.EventRegular, Name: "Club", Description: strings.Repeat("я", 1001)},
			wantErr: models.ErrLongDescription,
		},
		"bad visibility": {
			params:  models.EventCreateParams{Type: models.EventRegular, Name: "Club", Visibility: "secret"},
			wantErr: models.ErrInvalidVisibility,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			if testCase.wantErr == nil {
				repo.EXPECT().Create(mock.Anything, testCase.want).Return(1, nil)
			}

			_, err := NewService(repo, nil, nil).Create(context.Background(), testCase.params)
			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdate_finishedEventOnlyVisibility(t *testing.T) {
	finished := models.Event{
		ID: 1, Name: "Club", Description: "d", Visibility: models.EventVisibilityPublic,
		RatingsReach: models.ContentReachPublic, GamesReach: models.ContentReachPublic,
		FinishedAt: new(time.Now()),
	}

	testData := map[string]struct {
		params  models.UpdateEventParams
		wantErr error
	}{
		"visibility changes": {
			params: models.UpdateEventParams{ID: 1, Visibility: new(models.EventVisibilityPrivate)},
		},
		"the whole form sent, only visibility changed": {
			params: models.UpdateEventParams{
				ID: 1, Name: new("Club"), Description: new("d"), Visibility: new(models.EventVisibilityPrivate),
				RatingsReach: new(models.ContentReachPublic), GamesReach: new(models.ContentReachPublic),
			},
		},
		"name changes": {
			params:  models.UpdateEventParams{ID: 1, Name: new("Other")},
			wantErr: models.ErrEventAlreadyFinished,
		},
		"reach changes": {
			params:  models.UpdateEventParams{ID: 1, RatingsReach: new(models.ContentReachStaff)},
			wantErr: models.ErrEventAlreadyFinished,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, 1).Return(finished, nil)
			if testCase.wantErr == nil {
				repo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
			}

			err := NewService(repo, nil, nil).Update(context.Background(), testCase.params)
			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdate_validates(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, Name: "Club"}, nil).Times(3)
	svc := NewService(repo, nil, nil)

	assert.ErrorIs(t, svc.Update(context.Background(), models.UpdateEventParams{ID: 1, Name: new(" ")}), models.ErrEmptyName)
	assert.ErrorIs(t, svc.Update(context.Background(), models.UpdateEventParams{ID: 1, Visibility: new(models.EventVisibility("x"))}), models.ErrInvalidVisibility)

	// An empty description clears it.
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(e models.Event) bool { return e.Description == "" })).Return(nil)
	assert.NoError(t, svc.Update(context.Background(), models.UpdateEventParams{ID: 1, Description: new("  ")}))
}

func TestMembers_finishedEventIsReadOnly(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, FinishedAt: new(time.Now())}, nil).Twice()
	svc := NewService(repo, nil, nil)

	assert.ErrorIs(t, svc.AddMemberRoles(context.Background(), 1, 5, []models.EventRole{models.EventPlayer}), models.ErrEventAlreadyFinished)
	assert.ErrorIs(t, svc.RemoveMembers(context.Background(), 1, []int{5}), models.ErrEventAlreadyFinished)
}
