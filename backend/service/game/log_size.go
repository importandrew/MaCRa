package game

import (
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// checkLogSize refuses an entry made from an outdated table: seen is the
// number of round log entries the client saw, logSize the number the game has
// now. Without it a round entered twice would be stored as the next round, and
// a penalty entered twice would be fined twice. The repo checks again while
// saving, in case the log changes in between.
func checkLogSize(logSize, seen int) error {
	if seen != logSize {
		return fmt.Errorf("%w (seen %d log entries, the game has %d)", models.ErrRoundAlreadyEntered, seen, logSize)
	}

	return nil
}
