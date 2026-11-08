package player

import (
	"context"
)

func (s *Service) Update(ctx context.Context, id int, name string) error {
	if err := checkName(name); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, id, name); err != nil {
		return err
	}

	return nil
}
