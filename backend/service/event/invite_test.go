package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestNewInviteToken(t *testing.T) {
	repo := NewMockEventRepo(t)
	var saved string
	repo.EXPECT().SetInviteToken(mock.Anything, 1, mock.Anything).
		RunAndReturn(func(_ context.Context, _ int, token *string) error { saved = *token; return nil })

	token, err := NewService(repo, nil, nil).NewInviteToken(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, saved, token)
	assert.Len(t, token, 22, "16 random bytes, base64url")
}

func TestInvite(t *testing.T) {
	event := models.Event{ID: 1, Name: "Club", Type: models.EventRegular, OwnerUserID: 10}

	testData := map[string]struct {
		userID int
		roles  []models.EventRole
		want   bool
	}{
		"visitor":    {userID: 0, want: false},
		"not member": {userID: 5, roles: []models.EventRole{}, want: false},
		"member":     {userID: 5, roles: []models.EventRole{models.EventObserver}, want: true},
		"owner":      {userID: 10, roles: []models.EventRole{}, want: true},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			repo := NewMockEventRepo(t)
			repo.EXPECT().EventByInviteToken(mock.Anything, "tok").Return(event, nil)
			if testCase.userID != 0 {
				repo.EXPECT().GetMemberRoles(mock.Anything, 1, testCase.userID).Return(testCase.roles, nil)
			}

			invite, err := NewService(repo, nil, nil).Invite(context.Background(), "tok", testCase.userID)

			require.NoError(t, err)
			assert.Equal(t, models.EventInvite{EventID: 1, Name: "Club", Type: models.EventRegular, IsMember: testCase.want}, invite)
		})
	}
}
