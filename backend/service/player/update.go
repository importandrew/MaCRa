package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// maxNameLength matches the users.name VARCHAR(50) column.
const maxNameLength = 50

func (s *Service) Update(ctx context.Context, id int, name string) error {
	if name == "" {
		return models.ErrEmptyName
	}

	if len(name) > maxNameLength {
		return models.ErrLongName
	}

	if err := s.repo.Update(ctx, id, name); err != nil {
		return err
	}

	return nil
}
