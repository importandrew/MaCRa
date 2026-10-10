package player

import (
	"context"
	"time"

	"github.com/importandrew/MaCRa/models"
)

type PlayerRepo interface {
	Add(ctx context.Context, email, name, hash string) (int, error)
	GetByEmail(ctx context.Context, email string) (models.Player, error)
	Get(ctx context.Context, id int) (models.Player, error)
	GetPassByEmail(ctx context.Context, email string) (string, error)
	ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error)
	Update(ctx context.Context, id int, name string) error
	IsEmailVerified(ctx context.Context, userID int) (bool, error)
	SetEmailVerified(ctx context.Context, userID int) error
	UpdatePassword(ctx context.Context, userID int, hash string) error
	InvalidateAuthTokens(ctx context.Context, userID int, purpose models.AuthPurpose) error
	CreateAuthToken(ctx context.Context, userID int, purpose models.AuthPurpose, tokenHash string, expiresAt time.Time) error
	GetAuthTokenByHash(ctx context.Context, tokenHash string) (models.AuthToken, error)
	MarkAuthTokenUsed(ctx context.Context, id int) error
}

type Mailer interface {
	SendVerification(to, origin, token string) error
	SendPasswordReset(to, origin, token string) error
}

type Service struct {
	repo PlayerRepo
	mail Mailer
}

func NewService(playerRepo PlayerRepo, mail Mailer) *Service {
	return &Service{repo: playerRepo, mail: mail}
}
