package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Add(ctx context.Context, email, name, password string) (models.Player, error) {
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
