package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestListMembers(t *testing.T) {
	const eventID, ownerID = 1, 10

	all := []models.EventMember{
		{UserID: ownerID, Name: "Owner", Roles: []models.EventRole{models.EventPlayer, models.EventReferee}},
		{UserID: 2, Name: "Player", Roles: []models.EventRole{models.EventPlayer}},
		{UserID: 3, Name: "Observer", Roles: []models.EventRole{models.EventObserver}},
		{UserID: 4, Name: "Playing referee", Roles: []models.EventRole{models.EventPlayer, models.EventReferee}},
	}
	staff := []models.EventMember{
		{UserID: ownerID, Name: "Owner", Roles: []models.EventRole{models.EventReferee}},
		{UserID: 3, Name: "Observer", Roles: []models.EventRole{models.EventObserver}},
		{UserID: 4, Name: "Playing referee", Roles: []models.EventRole{models.EventReferee}},
	}

	testData := map[string]struct {
		callerID    int
		callerRoles []models.EventRole
		expected    []models.EventMember
	}{
		"member gets everyone":            {callerID: 2, callerRoles: []models.EventRole{models.EventPlayer}, expected: all},
		"owner gets everyone":             {callerID: ownerID, callerRoles: []models.EventRole{models.EventPlayer}, expected: all},
		"signed-in non-member gets staff": {callerID: 99, expected: staff},
		"anonymous gets staff":            {callerID: 0, expected: staff},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{
				ID:          eventID,
				Type:        models.EventRegular,
				OwnerUserID: ownerID,
				Visibility:  models.EventVisibilityPublic,
			}, nil)
			if testCase.callerID != 0 {
				repo.EXPECT().GetMemberRoles(mock.Anything, eventID, testCase.callerID).Return(testCase.callerRoles, nil)
			}
			repo.EXPECT().ListMembers(mock.Anything, eventID).Return(all, nil)

			members, err := NewService(repo, nil).ListMembers(context.Background(), eventID, testCase.callerID)

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, members)
		})
	}
}

func TestStaffOnly_ownerWithoutRoles(t *testing.T) {
	// A tournament owner may have no roles at all; they're still staff.
	members := []models.EventMember{
		{UserID: 10, Name: "Owner", Roles: []models.EventRole{}},
		{UserID: 2, Name: "Player", Roles: []models.EventRole{models.EventPlayer}},
	}

	assert.Equal(t, []models.EventMember{{UserID: 10, Name: "Owner", Roles: []models.EventRole{}}}, staffOnly(members, 10))
}

func TestDetails(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, Name: "Club"}, nil)
	repo.EXPECT().CountPlayers(mock.Anything, 1).Return(12, nil)

	details, err := NewService(repo, nil).Details(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, models.EventDetails{ID: 1, Name: "Club", PlayerCount: 12}, details)
}
