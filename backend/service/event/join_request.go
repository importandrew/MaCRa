package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// UpdateJoining sets whether the event takes requests to join and its
// player limit (Manage → Members), refused in a finished event.
func (s *Service) UpdateJoining(ctx context.Context, eventID int, settings models.JoiningSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateJoining(ctx, eventID, settings)
}

// takesRequests reports whether the event takes new requests to join: it is
// public, not finished, and the owner has switched requests on.
func takesRequests(event models.Event) bool {
	return event.Visibility == models.EventVisibilityPublic &&
		event.FinishedAt == nil && event.AcceptRequests
}

// joinRequestState is what the overview offers the caller about joining by a
// request (models.EventMe.JoinRequest): nothing to a member; the state of a
// pending request even when the event no longer takes requests (the user
// may withdraw it); otherwise, while the event takes requests, whether the
// user was declined, the event is full, or they may ask (a visitor too, who
// signs in first).
func (s *Service) joinRequestState(ctx context.Context, authCtx EventAuthContext) (models.JoinRequestState, error) {
	if authCtx.isParticipating() {
		return "", nil
	}

	event := authCtx.Event
	var status models.JoinRequestStatus
	if authCtx.CallerID != 0 {
		var err error
		if status, err = s.repo.JoinRequestStatus(ctx, event.ID, authCtx.CallerID); err != nil {
			return "", err
		}
	}

	switch {
	case status == models.JoinRequestPending:
		return models.JoinStatePending, nil
	case !takesRequests(event):
		return "", nil
	case status == models.JoinRequestDeclined:
		return models.JoinStateDeclined, nil
	}

	full, err := s.isFull(ctx, event)
	if err != nil {
		return "", err
	}
	if full {
		return models.JoinStateFull, nil
	}
	return models.JoinStateAvailable, nil
}

// isFull reports whether the event has as many players as its limit allows.
func (s *Service) isFull(ctx context.Context, event models.Event) (bool, error) {
	if event.PlayerLimit == nil {
		return false, nil
	}
	count, err := s.repo.CountPlayers(ctx, event.ID)
	if err != nil {
		return false, err
	}
	return count >= *event.PlayerLimit, nil
}

// RequestToJoin sends the user's request to join the event.
func (s *Service) RequestToJoin(ctx context.Context, eventID, userID int) error {
	return s.repo.CreateJoinRequest(ctx, eventID, userID)
}

// WithdrawJoinRequest deletes the user's pending request.
func (s *Service) WithdrawJoinRequest(ctx context.Context, eventID, userID int) error {
	return s.repo.WithdrawJoinRequest(ctx, eventID, userID)
}

// JoinRequests returns the event's pending requests, oldest first.
func (s *Service) JoinRequests(ctx context.Context, eventID int) ([]models.JoinRequest, error) {
	return s.repo.ListJoinRequests(ctx, eventID)
}

// AcceptJoinRequest makes the user a player of the event.
func (s *Service) AcceptJoinRequest(ctx context.Context, eventID, userID int) error {
	return s.repo.AcceptJoinRequest(ctx, eventID, userID)
}

// DeclineJoinRequest declines the user's request; they can't ask again.
func (s *Service) DeclineJoinRequest(ctx context.Context, eventID, userID int) error {
	return s.repo.DeclineJoinRequest(ctx, eventID, userID)
}
