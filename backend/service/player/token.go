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

// authTokenBytes is the amount of random data in an auth token (256 bits).
const authTokenBytes = 32

func generateAuthToken() string {
	buf := make([]byte, authTokenBytes)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func hashAuthToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// issueToken creates a token for the user and purpose, replacing the unused
// ones of that purpose. It returns models.ErrAuthMailLimit if the user's
// address has had too many auth emails lately (models.AuthMailLimits).
func (s *Service) issueToken(ctx context.Context, userID int, purpose models.AuthPurpose, ttl time.Duration) (string, error) {
	token := generateAuthToken()
	hash := hashAuthToken(token)

	err := s.repo.IssueAuthToken(ctx, userID, purpose, hash, time.Now().Add(ttl), models.AuthMailLimits)
	if err != nil {
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
