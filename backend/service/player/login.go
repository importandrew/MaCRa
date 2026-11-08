package player

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Login(ctx context.Context, email, password string) (models.Player, error) {
	storedPass, err := s.repo.GetPassByEmail(ctx, email)
	if errors.Is(err, models.ErrPlayerNotFound) {
		// Compare against a dummy hash, so an unknown email takes as long as a
		// wrong password and gets the same answer: sign-in must not reveal
		// which emails are registered.
		_ = validatePassword(dummyHash(), password)
		return models.Player{}, models.ErrBadCredentials
	}
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

	verified, err := s.repo.IsEmailVerified(ctx, player.ID)
	if err != nil {
		return models.Player{}, err
	}

	if !verified {
		return models.Player{}, models.ErrEmailNotVerified
	}

	return player, nil
}

// dummyHash is a bcrypt hash of no one's password, with the cost real hashes use.
var dummyHash = sync.OnceValue(func() string {
	hash, err := hashPassword("macra: no such user")
	if err != nil {
		panic(err)
	}
	return hash
})

func validatePassword(storedPass, incomingPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(incomingPass))
	if err != nil {
		return err
	}

	return nil
}
