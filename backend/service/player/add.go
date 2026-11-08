package player

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Add(ctx context.Context, email, name, password string) (models.Player, error) {
	if err := checkEmail(email); err != nil {
		return models.Player{}, err
	}

	if err := checkName(name); err != nil {
		return models.Player{}, err
	}

	if err := checkNewPassword(password); err != nil {
		return models.Player{}, err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return models.Player{}, err
	}

	id, err := s.repo.Add(ctx, email, name, hash)
	if err != nil {
		return models.Player{}, err
	}

	player := models.Player{
		ID:   id,
		Name: name,
	}

	return player, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
