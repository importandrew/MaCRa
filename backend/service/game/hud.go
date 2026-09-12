package game

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) HUD(ctx context.Context, gameID int, linger int) (models.GameHUD, error) {
	status, err := s.Status(ctx, gameID)
	if err != nil {
		return models.GameHUD{}, err
	}

	var hud models.GameHUD
	switch status {
	case models.GameEnded:
		hud, err = s.finalHUD(ctx, gameID)
	case models.GameLive:
		hud, err = s.liveHUD(ctx, gameID, linger)
	default:
		return models.GameHUD{}, fmt.Errorf("unexpected game status: %v", status)
	}
	if err != nil {
		return models.GameHUD{}, err
	}

	return hud, nil
}

func (s *Service) finalHUD(ctx context.Context, gameID int) (models.GameHUD, error) {
	result, err := s.GameResult(ctx, gameID)
	if err != nil {
		return models.GameHUD{}, err
	}

	hud := models.GameHUD{
		Final: result.PlayerResults,
	}

	return hud, nil
}

func (s *Service) liveHUD(ctx context.Context, gameID, linger int) (models.GameHUD, error) {
	gameLog, err := s.Logs(ctx, gameID)
	if err != nil {
		return models.GameHUD{}, err
	}

	players, err := s.Players(ctx, gameID)
	if err != nil {
		return models.GameHUD{}, err
	}

	game, err := s.LiveGame(ctx, gameID)
	if err != nil {
		return models.GameHUD{}, err
	}

	showDiffAndMessage := false
	if len(gameLog.RoundLogs) > 0 {
		lingerDuration := models.DefaultDiffLingerDuration
		if linger > 0 {
			lingerDuration = time.Duration(linger) * time.Second
		}

		// if right now is proper windows to show diff and message
		if !time.Now().After(gameLog.RoundLogs[len(gameLog.RoundLogs)-1].CreatedAt.Add(lingerDuration)) {
			showDiffAndMessage = true
		}
	}

	playerHUDs := make([]models.PlayerHUD, 0, len(players.Players))
	for _, player := range players.Players {
		playerHUD := models.PlayerHUD{
			Name:  player.Name,
			Wind:  game.CurrentWind(player.ID).Kanji(),
			Score: game.Score()[player.ID].GP,
		}

		if showDiffAndMessage {
			playerHUD.Diff = diff(gameLog.RoundLogs, game.StartWind(player.ID))
		}

		playerHUDs = append(playerHUDs, playerHUD)
	}

	// player HUDs are composed in starting order, reorder according to round seating
	reseatOrder := models.RoundSeats[game.Round()]

	round := game.Round()

	hud := models.GameHUD{
		Live: &models.GameLiveHUD{
			Round: models.Round{
				Num:   round,
				Kanji: models.RoundKanji(round),
			},
			PlayerHUDs: rearrangeSlice(playerHUDs, reseatOrder),
		},
	}

	if showDiffAndMessage {
		hud.Live.Message = message(gameLog, players)
	}

	return hud, nil
}

func diff(logs []models.RoundLog, wind models.Wind) string {
	count := len(logs)
	lastIdx := count - 1
	prevIdx := count - 2 //nolint: mnd

	switch count {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%+d", logs[lastIdx].Score(wind))
	default:
		return fmt.Sprintf("%+d", logs[lastIdx].Score(wind)-logs[prevIdx].Score(wind))
	}
}

func message(logs models.GameLog, players models.GamePlayers) string {
	if len(logs.RoundLogs) == 0 {
		return ""
	}

	log := logs.RoundLogs[len(logs.RoundLogs)-1]
	var result string
	switch log.Outcome {
	case models.Draw:
		result = "DRAW"
	case models.Self:
		result = fmt.Sprintf("SELF by %s [%d]", name(log.Winner, players), log.Value)
	case models.Deal:
		result = fmt.Sprintf("DEAL-IN from %s to %s [%d]", name(log.Loser, players), name(log.Winner, players), log.Value)
	case models.Penalty:
		result = fmt.Sprintf("PENALTY for %s [%d]", name(log.Loser, players), log.Value)
	case models.FalseMahjong:
		result = fmt.Sprintf("FALSE MAHJONG by %s [%d]", name(log.Loser, players), log.Value)
	}

	return fmt.Sprintf("Round %d: %s", log.Round, result)
}

func name(id int, players models.GamePlayers) string {
	idx := slices.IndexFunc(players.Players, func(p models.Player) bool {
		return p.ID == id
	})
	if idx == -1 {
		return "Unknown"
	}
	return players.Players[idx].Name
}

func rearrangeSlice[T any](original []T, indexOrder []models.Wind) []T {
	result := make([]T, 0, len(indexOrder))
	for _, index := range indexOrder {
		result = append(result, original[index])
	}
	return result
}
