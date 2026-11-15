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
	// Anyone who can open the event gets the full list (area 1); emails
	// aren't in it.
	all := []models.EventMember{
		{UserID: 10, Name: "Owner", Roles: []models.EventRole{models.EventPlayer, models.EventReferee}},
		{UserID: 2, Name: "Player", Roles: []models.EventRole{models.EventPlayer}},
	}
	repo := NewMockEventRepo(t)
	repo.EXPECT().ListMembers(mock.Anything, 1).Return(all, nil)

	members, err := NewService(repo, nil, nil).ListMembers(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, all, members)
}

func TestDetails(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1, Name: "Club"}, nil)
	repo.EXPECT().CountPlayers(mock.Anything, 1).Return(12, nil)
	repo.EXPECT().RulesFrozen(mock.Anything, 1).Return(true, nil)

	details, err := NewService(repo, nil, nil).Details(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, models.EventDetails{ID: 1, Name: "Club", PlayerCount: 12, RulesFrozen: true}, details)
}
