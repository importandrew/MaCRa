package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

type PlayerRepo interface {
	Add(ctx context.Context, email, name, hash string) (int, error)
	GetByEmail(ctx context.Context, email string) (models.Player, error)
	Get(ctx context.Context, id int) (models.Player, error)
	GetPassByEmail(ctx context.Context, email string) (string, error)
	ListByNameLimit(ctx context.Context, name string, limit int) ([]models.Player, error)
}

type Service struct {
	repo PlayerRepo
}

func NewService(playerRepo PlayerRepo) *Service {
	return &Service{repo: playerRepo}
}
