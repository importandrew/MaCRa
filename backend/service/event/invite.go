package event

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/importandrew/MaCRa/models"
)

// inviteTokenBytes is the amount of random data in an invite token (128
// bits): it can't be guessed, and the link stays short for a QR code.
const inviteTokenBytes = 16

func generateInviteToken() string {
	buf := make([]byte, inviteTokenBytes)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

// InviteToken returns the event's invite token (Manage → Members), nil when
// the link is off.
func (s *Service) InviteToken(ctx context.Context, eventID int) (*string, error) {
	return s.repo.InviteToken(ctx, eventID)
}

// NewInviteToken turns the invite link on with a new token, replacing an old
// link (a leaked one stops working), and returns it.
func (s *Service) NewInviteToken(ctx context.Context, eventID int) (string, error) {
	token := generateInviteToken()
	if err := s.repo.SetInviteToken(ctx, eventID, &token); err != nil {
		return "", err
	}
	return token, nil
}

// TurnOffInvite turns the invite link off.
func (s *Service) TurnOffInvite(ctx context.Context, eventID int) error {
	return s.repo.SetInviteToken(ctx, eventID, nil)
}

// Invite shows the invite link's event before joining: its name and type,
// and whether the user is already a member (false for a visitor, userID 0).
// The link is the access: its holder sees the event's name even in a
// private event.
func (s *Service) Invite(ctx context.Context, token string, userID int) (models.EventInvite, error) {
	event, err := s.repo.EventByInviteToken(ctx, token)
	if err != nil {
		return models.EventInvite{}, err
	}

	invite := models.EventInvite{EventID: event.ID, Name: event.Name, Type: event.Type}
	if userID != 0 {
		roles, err := s.repo.GetMemberRoles(ctx, event.ID, userID)
		if err != nil {
			return models.EventInvite{}, err
		}
		invite.IsMember = len(roles) > 0 || event.OwnerUserID == userID
	}
	return invite, nil
}

// JoinByInvite makes the user a player of the invite link's event and
// returns its ID.
func (s *Service) JoinByInvite(ctx context.Context, token string, userID int) (int, error) {
	return s.repo.JoinByInvite(ctx, token, userID)
}
