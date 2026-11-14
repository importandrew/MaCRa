package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Get(ctx context.Context, eventID int) (models.Event, error) {
	return s.repo.Get(ctx, eventID)
}

// Details returns the event with its player count and whether its rules are
// frozen, for the overview.
func (s *Service) Details(ctx context.Context, eventID int) (models.EventDetails, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.EventDetails{}, err
	}

	count, err := s.repo.CountPlayers(ctx, eventID)
	if err != nil {
		return models.EventDetails{}, err
	}

	frozen, err := s.repo.RulesFrozen(ctx, eventID)
	if err != nil {
		return models.EventDetails{}, err
	}

	return models.EventDetails{Event: event, PlayerCount: count, RulesFrozen: frozen}, nil
}

// ListVisible returns the events the user can open, with their roles, in
// the repo's order: active ones by their last game, then finished ones.
func (s *Service) ListVisible(ctx context.Context, userID int) ([]models.EventListItem, error) {
	rows, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	visible := make([]models.EventListItem, 0, len(rows))
	for _, row := range rows {
		authCtx := eventAuthContext(row.Event, userID, row.Roles)
		if !authCtx.canAccessEvent() {
			continue
		}

		item := models.EventListItem{Event: row.Event, LastGameAt: row.LastGameAt}
		if userID != 0 {
			item.IsParticipating = authCtx.isParticipating()
			item.IsOwner = authCtx.IsOwner
			item.Roles = authCtx.Roles
		}

		visible = append(visible, item)
	}

	return visible, nil
}
