package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// checkLogSize refuses an entry made from an outdated table: seen is the
// number of round log entries the client saw, logSize the number the game has
// now. Without it a round entered twice would be stored as the next round, and
// a penalty entered twice would be fined twice. A nil seen (a client that
// doesn't send it) skips the check; the repo still refuses an entry if the log
// changes while it is being saved.
func checkLogSize(logSize int, seen *int) error {
	if seen != nil && *seen != logSize {
		return fmt.Errorf("%w (seen %d log entries, the game has %d)", models.ErrRoundAlreadyEntered, *seen, logSize)
	}

	return nil
}
