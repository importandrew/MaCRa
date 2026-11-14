package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

// TestSetMemberRoles checks the roles the service passes on: each once. The
// repo checks the rest under its locks (repo/event member_roles_test.go).
func TestSetMemberRoles(t *testing.T) {
	repo := NewMockEventRepo(t)
	repo.EXPECT().SetMemberRoles(mock.Anything, 1, 20,
		[]models.EventRole{models.EventObserver, models.EventReferee}).Return(models.ErrPlayerRoleInGame)

	err := NewService(repo, nil, nil).SetMemberRoles(context.Background(), 1, 20,
		[]models.EventRole{models.EventObserver, models.EventObserver, models.EventReferee})

	assert.ErrorIs(t, err, models.ErrPlayerRoleInGame)
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
