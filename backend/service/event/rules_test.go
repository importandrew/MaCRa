package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

func TestUpdateRules(t *testing.T) {
	t.Run("valid rules are saved", func(t *testing.T) {
		rules := models.GameRules{Hands: 8, MinHand: 24, FalseMahjong: [2]int{30, 60}}
		repo := NewMockEventRepo(t)
		repo.EXPECT().UpdateRules(mock.Anything, 1, rules).Return(models.ErrEventRulesFrozen)

		err := NewService(repo, nil, nil).UpdateRules(context.Background(), 1, rules)

		assert.ErrorIs(t, err, models.ErrEventRulesFrozen)
	})

	t.Run("invalid rules are refused", func(t *testing.T) {
		repo := NewMockEventRepo(t)

		err := NewService(repo, nil, nil).UpdateRules(context.Background(), 1, models.GameRules{Hands: 12})

		assert.ErrorIs(t, err, models.ErrInvalidGameRules)
	})
}
