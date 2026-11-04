package player

import (
	"context"
	"unicode/utf8"

	"github.com/importandrew/MaCRa/models"
)

// maxNameLength matches the users.name VARCHAR(50) column, which counts
// characters, not bytes.
const maxNameLength = 50

func (s *Service) Update(ctx context.Context, id int, name string) error {
	if name == "" {
		return models.ErrEmptyName
	}

	// RuneCountInString, not len: len counts bytes, and a Cyrillic letter is
	// two bytes in UTF-8, so a 26-letter Russian name was rejected.
	if utf8.RuneCountInString(name) > maxNameLength {
		return models.ErrLongName
	}

	if err := s.repo.Update(ctx, id, name); err != nil {
		return err
	}

	return nil
}
