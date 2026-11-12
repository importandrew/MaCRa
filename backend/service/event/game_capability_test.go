package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestGameCapabilities(t *testing.T) {
	now := time.Now()
	ago := func(d time.Duration) *time.Time { return new(now.Add(-d)) }
	live := func(logSize int) *models.GameSummary {
		return &models.GameSummary{Status: models.GameLive, LogSize: logSize}
	}
	ended := func(finishedAgo time.Duration) *models.GameSummary {
		return &models.GameSummary{Status: models.GameEnded, FinishedAt: ago(finishedAgo), LogSize: 16}
	}
	referee := []models.EventRole{models.EventReferee}
	player := []models.EventRole{models.EventPlayer}

	testData := map[string]struct {
		event   models.EventType
		roles   []models.EventRole
		owner   bool
		seated  bool
		game    *models.GameSummary
		want    []models.Capability
		wantNot []models.Capability
	}{
		"seated player, live game without rounds: abort, finish, enter": {
			event: models.EventRegular, roles: player, seated: true, game: live(0),
			want:    []models.Capability{models.CapViewGame, models.CapEnterRound, models.CapAbortGame, models.CapFinishGame},
			wantNot: []models.Capability{models.CapUndoRound},
		},
		"seated player, after the first round: no abort": {
			event: models.EventRegular, roles: player, seated: true, game: live(1),
			want:    []models.Capability{models.CapEnterRound, models.CapFinishGame},
			wantNot: []models.Capability{models.CapAbortGame, models.CapUndoRound},
		},
		"referee, not seated: abort, finish and undo, no entering": {
			event: models.EventRegular, roles: referee, game: live(5),
			want:    []models.Capability{models.CapViewGame, models.CapAbortGame, models.CapFinishGame, models.CapUndoRound},
			wantNot: []models.Capability{models.CapEnterRound},
		},
		"referee, game ended 10 minutes ago: undo only": {
			event: models.EventRegular, roles: referee, game: ended(10 * time.Minute),
			want:    []models.Capability{models.CapUndoRound},
			wantNot: []models.Capability{models.CapAbortGame, models.CapFinishGame, models.CapEnterRound},
		},
		"referee, game ended two hours ago: nothing to do": {
			event: models.EventRegular, roles: referee, game: ended(2 * time.Hour),
			wantNot: []models.Capability{models.CapUndoRound, models.CapAbortGame, models.CapFinishGame},
		},
		"owner without the referee role sees the game but can't act as a referee": {
			event: models.EventRegular, owner: true, game: live(3),
			want:    []models.Capability{models.CapViewGame},
			wantNot: []models.Capability{models.CapUndoRound, models.CapAbortGame, models.CapFinishGame},
		},
		"tournament: no abort or finish even for a referee": {
			event: models.EventTournament, roles: referee, game: live(0),
			want:    []models.Capability{models.CapUndoRound},
			wantNot: []models.Capability{models.CapAbortGame, models.CapFinishGame},
		},
		"no game: none of the game capabilities": {
			event: models.EventRegular, roles: referee,
			wantNot: []models.Capability{models.CapViewGame, models.CapUndoRound, models.CapAbortGame, models.CapFinishGame, models.CapEnterRound},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			event := testEvent(testCase.event)
			authCtx := EventAuthContext{
				Event:          event,
				CallerID:       5,
				Roles:          testCase.roles,
				IsOwner:        testCase.owner,
				IsMember:       true,
				IsSeatedInGame: testCase.seated,
				Game:           testCase.game,
				Now:            now,
			}
			for _, c := range testCase.want {
				assert.True(t, authCtx.hasCapability(c), c)
			}
			for _, c := range testCase.wantNot {
				assert.False(t, authCtx.hasCapability(c), c)
			}
		})
	}
}
