package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateRules sets the event's game rules. They can change only while the
// event has no games (models.ErrEventRulesFrozen), so all games of one rating
// follow the same rules.
func (s *Service) UpdateRules(ctx context.Context, eventID int, rules models.GameRules) error {
	if err := rules.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateRules(ctx, eventID, rules)
}
