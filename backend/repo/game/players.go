package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

type PlayerWind struct {
	PlayerID     int    `db:"player_id"`
	StartingWind string `db:"starting_wind"`
}

func (r *Repo) Players(ctx context.Context, gameID int) (models.StartingPlayers, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.StartingPlayers{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT player_id, starting_wind 
		FROM game_players 
		WHERE game_id = $1`,
		gameID)
	if err != nil {
		return models.StartingPlayers{}, err
	}

	players, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerWind])
	if err != nil {
		return models.StartingPlayers{}, err
	}

	if len(players) != 4 {
		return models.StartingPlayers{}, fmt.Errorf("expected 4 players, got %d", len(players))
	}

	var startingPlayers models.StartingPlayers
	for _, p := range players {
		switch p.StartingWind {
		case "east":
			startingPlayers.East = p.PlayerID
		case "south":
			startingPlayers.South = p.PlayerID
		case "west":
			startingPlayers.West = p.PlayerID
		case "north":
			startingPlayers.North = p.PlayerID
		}
	}

	return startingPlayers, nil
}
