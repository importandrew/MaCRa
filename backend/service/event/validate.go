package event

import (
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/importandrew/MaCRa/models"
)

// eventName returns the trimmed name, or an error if it is empty or longer
// than models.MaxEventNameLength.
func eventName(name string) (string, error) {
	name = strings.TrimSpace(name)
	switch {
	case name == "":
		return "", models.ErrEmptyName
	case utf8.RuneCountInString(name) > models.MaxEventNameLength:
		return "", models.ErrLongEventName
	}
	return name, nil
}

// eventDescription returns the trimmed description, or an error if it is
// longer than models.MaxEventDescriptionLength.
func eventDescription(description string) (string, error) {
	description = strings.TrimSpace(description)
	if utf8.RuneCountInString(description) > models.MaxEventDescriptionLength {
		return "", models.ErrLongDescription
	}
	return description, nil
}

func validateVisibility(visibility models.EventVisibility) error {
	if !slices.Contains(models.ValidEventVisibilities, visibility) {
		return models.ErrInvalidVisibility
	}
	return nil
}

func validateReach(reach models.ContentReach) error {
	if !slices.Contains(models.ValidContentReach, reach) {
		return models.ErrInvalidContentReach
	}
	return nil
}

// checkNotFinished returns models.ErrEventAlreadyFinished for a finished
// event: its members and settings stay as they are (only the visibility can
// change, see Update).
func checkNotFinished(event models.Event) error {
	if event.FinishedAt != nil {
		return models.ErrEventAlreadyFinished
	}
	return nil
}
