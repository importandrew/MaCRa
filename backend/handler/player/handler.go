package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type PlayerService interface {
	Add(ctx context.Context, email, name, password string) (models.Player, error)
	Login(ctx context.Context, email, password string) (models.Player, error)
	ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error)
	Get(ctx context.Context, id int) (models.Player, error)
	Update(ctx context.Context, id int, name string) error
}

type Handler struct {
	player PlayerService
	domain string
	jwt    []byte
}

func NewHandler(player PlayerService, domain string, jwt []byte) Handler {
	return Handler{
		player: player,
		domain: domain,
		jwt:    jwt,
	}
}
