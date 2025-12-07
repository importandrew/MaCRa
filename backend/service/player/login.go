package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, email, password string) (models.Player, error) {
	storedPass, err := s.repo.GetPassByEmail(ctx, email)
	if err != nil {
		return models.Player{}, err
	}

	if err = validatePassword(storedPass, password); err != nil {
		return models.Player{}, models.ErrBadCredentials
	}

	player, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}

func validatePassword(storedPass, incomingPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(incomingPass))
	if err != nil {
		return err
	}

	return nil
}
