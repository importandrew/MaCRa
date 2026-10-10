package player

import (
	"context"

	"github.com/importandrew/MaCRa/config"
	"github.com/importandrew/MaCRa/models"
)

type PlayerService interface {
	SignUp(ctx context.Context, email, name, password, origin string) error
	Login(ctx context.Context, email, password string) (models.Player, error)
	ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error)
	Get(ctx context.Context, id int) (models.Player, error)
	Update(ctx context.Context, id int, name string) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerification(ctx context.Context, email, origin string) error
	ForgotPassword(ctx context.Context, email, origin string) error
	ResetPassword(ctx context.Context, token, password string) error
}

type Handler struct {
	player       PlayerService
	domain       string
	jwt          []byte
	allowOrigins []string
}

func NewHandler(player PlayerService, cfg config.Config) Handler {
	return Handler{
		player:       player,
		domain:       cfg.CookieDomain,
		jwt:          []byte(cfg.JWTKey),
		allowOrigins: cfg.AllowOrigins,
	}
}
