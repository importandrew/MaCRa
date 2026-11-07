package game

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

type PlayerWind struct {
	PlayerID     int    `db:"player_id"`
	StartingWind string `db:"starting_wind"`
}

// querier runs a query on a pool connection or inside a transaction.
type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (r *Repo) Players(ctx context.Context, gameID int) (models.StartingPlayers, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.StartingPlayers{}, err
	}
	defer conn.Release()

	return gamePlayers(ctx, conn, gameID)
}

// gamePlayers reads the game's starting seating through q. Code that already
// holds a connection or a transaction passes it here instead of calling
// Players, which would take a second pool connection: with enough concurrent
// requests every connection is held by a request waiting for another one, and
// they all hang.
func gamePlayers(ctx context.Context, q querier, gameID int) (models.StartingPlayers, error) {
	rows, err := q.Query(ctx, `
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
