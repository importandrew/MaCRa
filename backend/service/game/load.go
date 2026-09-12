package game

// func (s *Service) LoadLiveGames(ctx context.Context) error {
// 	games, err := s.gameRepo.LiveGames(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	// sort games by event
// 	eventGames := make(map[int][]int)
// 	for _, game := range games {
// 		_, ok := eventGames[game.EventID]
// 		if !ok {
// 			eventGames[game.EventID] = make([]int, 0)
// 		}
// 		eventGames[game.EventID] = append(eventGames[game.EventID], game.GameID)
// 	}

// 	for eventID, gameIDs := range eventGames {
// 		for _, gameID := range gameIDs {
// 			err = s.loadGame(ctx, eventID, gameID)
// 			if err != nil {
// 				return fmt.Errorf("load game (event ID %d, game ID %d): %w", eventID, gameID, err)
// 			}
// 		}
// 	}

// 	return nil
// }

// func (s *Service) loadGame(ctx context.Context, eventID, gameID int) error {
// 	players, err := s.gameRepo.Players(ctx, gameID)
// 	if err != nil {
// 		return err
// 	}

// 	seating := game.Seating([4]int{
// 		players.East,
// 		players.South,
// 		players.West,
// 		players.North,
// 	})

// 	if err := s.validatePlayers(seating[:], eventID, &gameID); err != nil {
// 		return err
// 	}

// 	gameCfg := game.Config{
// 		ID:      gameID,
// 		Players: seating,
// 		Rules:   game.Rules{},
// 	}

// 	g := game.NewGame(gameCfg)

// 	roundLogs, err := s.gameRepo.GameLog(ctx, gameID)
// 	if err != nil {
// 		return err
// 	}

// 	for _, log := range roundLogs {
// 		switch log.Outcome {
// 		case models.FalseMahjong:
// 			err = g.FalseMahjong(log.Loser, log.Value, log.ForceNextRound)
// 		case models.Penalty:
// 			err = g.Penalty(log.Loser, log.Value, log.ForceNextRound)
// 		default:
// 			err = g.EnterRoundResult(log.Outcome, log.Value, log.Loser, log.Winner)
// 		}
// 		if err != nil {
// 			return fmt.Errorf("log record %d game %d: %w", log.ID, gameID, err)
// 		}
// 	}

// 	_, ok := s.liveGames[eventID]
// 	if !ok {
// 		s.liveGames[eventID] = make(map[int]GameEngine)
// 	}
// 	s.liveGames[eventID][gameID] = g

// 	return nil
// }
