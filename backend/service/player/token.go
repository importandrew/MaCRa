package player

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/importandrew/MaCRa/models"
)

func generateAuthToken() string {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func hashAuthToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) issueToken(ctx context.Context, userID int, purpose models.AuthPurpose, ttl time.Duration) (string, error) {
	token := generateAuthToken()
	hash := hashAuthToken(token)

	if err := s.repo.InvalidateAuthTokens(ctx, userID, purpose); err != nil {
		return "", err
	}

	if err := s.repo.CreateAuthToken(ctx, userID, purpose, hash, time.Now().Add(ttl)); err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) lookupToken(ctx context.Context, raw string, purpose models.AuthPurpose) (models.AuthToken, error) {
	token, err := s.repo.GetAuthTokenByHash(ctx, hashAuthToken(raw))
	if err != nil {
		return models.AuthToken{}, err
	}

	if token.Purpose != purpose || token.UsedAt.Valid || time.Now().After(token.ExpiresAt) {
		return models.AuthToken{}, models.ErrInvalidAuthToken
	}

	return token, nil
}
