package game

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoadLiveGames(t *testing.T) {
	testData := map[string]struct {
		mockSetup func(*MockGameRepo, *MockPlayerRepo)

		expectedErr assert.ErrorAssertionFunc
		verifyGames func(*testing.T, *Service)
	}{
		"ok single game": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{100}, nil)

				// Setup loadGame mocks
				mg.EXPECT().Players(mock.Anything, 100).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mg.EXPECT().GameLog(mock.Anything, 100).Return([]models.RoundLog{}, nil)
			},
			expectedErr: assert.NoError,
			verifyGames: func(t *testing.T, s *Service) {
				g, err := s.GetLiveGame(100)
				assert.NoError(t, err)
				assert.NotNil(t, g)
				assert.Equal(t, 100, g.ID())
			},
		},
		"ok multiple games": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{100, 200}, nil)

				// Game 100
				mg.EXPECT().Players(mock.Anything, 100).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mg.EXPECT().GameLog(mock.Anything, 100).Return([]models.RoundLog{}, nil)

				// Game 200
				mg.EXPECT().Players(mock.Anything, 200).Return(
					models.StartingPlayers{
						East:  5,
						South: 6,
						West:  7,
						North: 8,
					},
					nil,
				)
				mg.EXPECT().GameLog(mock.Anything, 200).Return([]models.RoundLog{}, nil)
			},
			expectedErr: assert.NoError,
			verifyGames: func(t *testing.T, s *Service) {
				g1, err := s.GetLiveGame(100)
				assert.NoError(t, err)
				assert.NotNil(t, g1)

				g2, err := s.GetLiveGame(200)
				assert.NoError(t, err)
				assert.NotNil(t, g2)
			},
		},
		"ok game with round logs": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{100}, nil)

				mg.EXPECT().Players(mock.Anything, 100).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mg.EXPECT().GameLog(mock.Anything, 100).Return(
					[]models.RoundLog{
						{
							ID:            1,
							Round:         1,
							EventSequence: 1,
							Outcome:       models.Deal,
							Winner:        1,
							Loser:         2,
							Value:         8,
							CreatedAt:     time.Now(),
						},
					},
					nil,
				)
			},
			expectedErr: assert.NoError,
			verifyGames: func(t *testing.T, s *Service) {
				g, err := s.GetLiveGame(100)
				assert.NoError(t, err)
				assert.NotNil(t, g)
				// Game should have advanced to round 2 after processing the log
				assert.Equal(t, 2, g.Round())
			},
		},
		"live games error": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"players error": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{100}, nil)
				mg.EXPECT().Players(mock.Anything, 100).Return(
					models.StartingPlayers{},
					errors.New("database error"),
				)
			},
			expectedErr: assert.Error,
		},
		"game log error": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{100}, nil)
				mg.EXPECT().Players(mock.Anything, 100).Return(
					models.StartingPlayers{
						East:  1,
						South: 2,
						West:  3,
						North: 4,
					},
					nil,
				)
				mg.EXPECT().GameLog(mock.Anything, 100).Return(nil, errors.New("database error"))
			},
			expectedErr: assert.Error,
		},
		"empty live games": {
			mockSetup: func(mg *MockGameRepo, mp *MockPlayerRepo) {
				mg.EXPECT().LiveGames(mock.Anything).Return([]int{}, nil)
			},
			expectedErr: assert.NoError,
			verifyGames: func(t *testing.T, s *Service) {
				// No games should be loaded
				_, err := s.GetLiveGame(100)
				assert.Error(t, err)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameRepo := NewMockGameRepo(t)
			mockPlayerRepo := NewMockPlayerRepo(t)
			service := NewService(mockGameRepo, mockPlayerRepo)

			testCase.mockSetup(mockGameRepo, mockPlayerRepo)

			ctx := context.Background()
			err := service.LoadLiveGames(ctx)

			testCase.expectedErr(t, err)
			if err != nil {
				return
			}

			if testCase.verifyGames != nil {
				testCase.verifyGames(t, service)
			}
		})
	}
}

func TestLoadGame_WithDifferentOutcomes(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	mockGameRepo.EXPECT().Players(mock.Anything, 100).Return(
		models.StartingPlayers{
			East:  1,
			South: 2,
			West:  3,
			North: 4,
		},
		nil,
	)

	// Test loading game with different outcome types
	mockGameRepo.EXPECT().GameLog(mock.Anything, 100).Return(
		[]models.RoundLog{
			{
				ID:             1,
				Round:          1,
				EventSequence:  1,
				Outcome:        models.FalseMahjong,
				Loser:          1,
				Value:          30,
				ForceNextRound: false,
				CreatedAt:      time.Now(),
			},
			{
				ID:             2,
				Round:          2,
				EventSequence:  1,
				Outcome:        models.Penalty,
				Loser:          2,
				Value:          10,
				ForceNextRound: true,
				CreatedAt:      time.Now(),
			},
			{
				ID:            3,
				Round:         3,
				EventSequence: 1,
				Outcome:       models.Deal,
				Winner:        3,
				Loser:         4,
				Value:         8,
				CreatedAt:     time.Now(),
			},
		},
		nil,
	)

	ctx := context.Background()
	err := service.loadGame(ctx, 100)

	assert.NoError(t, err)
	g, err := service.GetLiveGame(100)
	assert.NoError(t, err)
	assert.NotNil(t, g)
	// Game starts at round 1, "false mahjong" + "penalty" with ForceNextRound=true + "deal-in"
	// means it's 3rd round now
	assert.Equal(t, 3, g.Round())
}

func TestLoadGame_InvalidLogError(t *testing.T) {
	mockGameRepo := NewMockGameRepo(t)
	mockPlayerRepo := NewMockPlayerRepo(t)
	service := NewService(mockGameRepo, mockPlayerRepo)

	mockGameRepo.EXPECT().Players(mock.Anything, 100).Return(
		models.StartingPlayers{
			East:  1,
			South: 2,
			West:  3,
			North: 4,
		},
		nil,
	)

	// Invalid log with hand value < 8
	mockGameRepo.EXPECT().GameLog(mock.Anything, 100).Return(
		[]models.RoundLog{
			{
				ID:            1,
				Round:         1,
				EventSequence: 1,
				Outcome:       models.Deal,
				Winner:        1,
				Loser:         2,
				Value:         5, // Invalid, should be >= 8
				CreatedAt:     time.Now(),
			},
		},
		nil,
	)

	ctx := context.Background()
	err := service.loadGame(ctx, 100)

	assert.Error(t, err)
}
