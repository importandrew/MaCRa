package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

func TestSetMemberRoles(t *testing.T) {
	const eventID, ownerID, memberID = 1, 10, 20
	event := models.Event{ID: eventID, OwnerUserID: ownerID}
	referee := []models.EventRole{models.EventReferee}
	observer := []models.EventRole{models.EventObserver}

	testData := map[string]struct {
		event   models.Event
		userID  int
		roles   []models.EventRole
		current []models.EventRole
		// saved: the roles the repo is asked to set; nil if it isn't.
		saved       []models.EventRole
		expectedErr error
	}{
		"referee becomes an observer": {
			event: event, userID: memberID, roles: observer, current: referee, saved: observer,
		},
		"repeated roles are set once": {
			event: event, userID: memberID,
			roles:   []models.EventRole{models.EventObserver, models.EventObserver, models.EventReferee},
			current: referee,
			saved:   []models.EventRole{models.EventObserver, models.EventReferee},
		},
		"the owner drops every role": {
			event: event, userID: ownerID, roles: []models.EventRole{},
			current: referee, saved: []models.EventRole{},
		},
		"the owner without roles takes one": {
			event: event, userID: ownerID, roles: referee, current: nil, saved: referee,
		},
		"a member can't be left without roles": {
			event: event, userID: memberID, roles: []models.EventRole{}, current: referee,
			expectedErr: models.ErrInvalidRole,
		},
		"not a member": {
			event: event, userID: memberID, roles: referee, current: nil,
			expectedErr: models.ErrNotMember,
		},
		"finished event": {
			event:  models.Event{ID: eventID, OwnerUserID: ownerID, FinishedAt: new(time.Now())},
			userID: memberID, roles: referee,
			expectedErr: models.ErrEventAlreadyFinished,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, eventID).Return(testCase.event, nil)
			if testCase.event.FinishedAt == nil {
				repo.EXPECT().GetMemberRoles(mock.Anything, eventID, testCase.userID).Return(testCase.current, nil)
			}
			if testCase.saved != nil {
				repo.EXPECT().SetMemberRoles(mock.Anything, eventID, testCase.userID, testCase.saved).Return(nil)
			}

			err := NewService(repo, nil, nil).SetMemberRoles(context.Background(), eventID, testCase.userID, testCase.roles)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestSetMemberRoles_invalidRole(t *testing.T) {
	err := NewService(NewMockEventRepo(t), nil, nil).
		SetMemberRoles(context.Background(), 1, 20, []models.EventRole{"admin"})

	assert.ErrorIs(t, err, models.ErrInvalidRole)
}

func TestAddMemberRoles_needsValidRoles(t *testing.T) {
	for name, roles := range map[string][]models.EventRole{
		"unknown role": {"admin"},
		"no roles":     {},
	} {
		t.Run(name, func(t *testing.T) {
			err := NewService(NewMockEventRepo(t), nil, nil).AddMemberRoles(context.Background(), 1, 20, roles)

			assert.ErrorIs(t, err, models.ErrInvalidRole)
		})
	}
}

func TestRemoveMembers_notTheOwner(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, OwnerUserID: 10}, nil)

	err := NewService(repo, nil, nil).RemoveMembers(context.Background(), 1, []int{20, 10})

	assert.ErrorIs(t, err, models.ErrOwnerNotRemovable)
}
