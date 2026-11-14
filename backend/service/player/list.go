package player

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/importandrew/MaCRa/models"
)

const (
	// minSearchLength is how many characters a name search needs.
	minSearchLength = 2
	// maxSearchResults is how many players a search returns at most.
	maxSearchResults = 20
)

// Search finds players for the player pickers (docs/user-scenarios.md,
// "Finding players"): a number is a user ID and matches exactly; otherwise
// names containing the query, ignoring case, exact matches first, then
// prefixes, then the rest, at most maxSearchResults. A name query shorter
// than minSearchLength finds nobody.
func (s *Service) Search(ctx context.Context, query string) ([]models.Player, error) {
	query = strings.TrimSpace(query)

	if id, err := strconv.Atoi(query); err == nil && id > 0 && !strings.HasPrefix(query, "+") {
		player, err := s.repo.Get(ctx, id)
		switch {
		case errors.Is(err, models.ErrPlayerNotFound):
			return []models.Player{}, nil
		case err != nil:
			return nil, err
		default:
			return []models.Player{player}, nil
		}
	}

	if utf8.RuneCountInString(query) < minSearchLength {
		return []models.Player{}, nil
	}

	return s.repo.SearchByName(ctx, query, maxSearchResults)
}
