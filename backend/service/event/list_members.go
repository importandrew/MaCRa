package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// ListMembers returns the event's members with their roles. Members (and the
// owner) get the full list; anyone else who can open the event gets only the
// staff, see staffOnly: players may not want to be listed to outsiders.
func (s *Service) ListMembers(ctx context.Context, eventID, callerID int) ([]models.EventMember, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return nil, err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, 0)
	if err != nil {
		return nil, err
	}

	members, err := s.repo.ListMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if authCtx.isParticipating() {
		return members, nil
	}

	return staffOnly(members, event.OwnerUserID), nil
}

// staffOnly keeps the owner, referees and observers, each with only their
// staff roles, so the list doesn't tell who else plays.
func staffOnly(members []models.EventMember, ownerID int) []models.EventMember {
	staff := make([]models.EventMember, 0, len(members))
	for _, member := range members {
		roles := make([]models.EventRole, 0, len(member.Roles))
		for _, role := range member.Roles {
			if role == models.EventReferee || role == models.EventObserver {
				roles = append(roles, role)
			}
		}

		if len(roles) == 0 && member.UserID != ownerID {
			continue
		}

		member.Roles = roles
		staff = append(staff, member)
	}

	return staff
}
