package game

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/models"
)

const testEventID = 1

var defaultStartingPlayers = models.StartingPlayers{
	East:  1,
	South: 2,
	West:  3,
	North: 4,
}

func expectLoadGame(m *MockGameRepo, gameID int, logs []models.RoundLog) {
	m.EXPECT().Players(mock.Anything, gameID).Return(defaultStartingPlayers, nil).Once()
	m.EXPECT().Rules(mock.Anything, gameID).Return(models.StandardGameRules(), nil).Once()
	if logs == nil {
		logs = []models.RoundLog{}
	}
	m.EXPECT().GameLog(mock.Anything, gameID).Return(logs, nil).Once()
}

func expectCurrentGame(m *MockGameRepo, playerID, gameID int, logs []models.RoundLog, times int) {
	for range times {
		expectCurrentGameLoad(m, playerID, gameID, logs)
	}
}

func expectCurrentGameLoad(m *MockGameRepo, playerID, gameID int, logs []models.RoundLog) {
	m.EXPECT().CurrentGame(mock.Anything, playerID, testEventID).Return(gameID, nil)
	expectLoadGame(m, gameID, logs)
}

// expectCurrentGameLoads expects the entry's load of the live game, then the
// table's loads (Now) for the rest.
func expectCurrentGameLoads(m *MockGameRepo, playerID, gameID int, logsPerLoad [][]models.RoundLog) {
	for i, logs := range logsPerLoad {
		if i == 0 {
			expectCurrentGameLoad(m, playerID, gameID, logs)
		} else {
			expectTableLoad(m, playerID, gameID, logs)
		}
	}
}

// expectTableLoad expects Now's load of the player's live game: not a
// tournament game.
func expectTableLoad(m *MockGameRepo, playerID, gameID int, logs []models.RoundLog) {
	m.EXPECT().SeatedGame(mock.Anything, playerID, testEventID).Return(gameID, nil)
	m.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil)
	expectLoadGame(m, gameID, logs)
	m.EXPECT().TourSeat(mock.Anything, gameID).Return(nil, nil)
}

func expectLiveGame(m *MockGameRepo, gameID int, logs []models.RoundLog) {
	m.EXPECT().Status(mock.Anything, gameID).Return(models.GameLive, nil)
	expectLoadGame(m, gameID, logs)
}

func expectPlayerGets(mp *MockPlayerRepo) {
	names := map[int]string{1: "Alice", 2: "Bob", 3: "Charlie", 4: "David"}
	for id, name := range names {
		mp.On("Get", mock.Anything, id).Return(models.Player{ID: id, Name: name}, nil).Maybe()
	}
}

func drawRoundLogs(count int) []models.RoundLog {
	logs := make([]models.RoundLog, count)
	now := time.Now()
	for i := range count {
		logs[i] = models.RoundLog{
			ID:            i + 1,
			Round:         i + 1,
			EventSequence: 1,
			Outcome:       models.Draw,
			CreatedAt:     now,
		}
	}
	return logs
}
