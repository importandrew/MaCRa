package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/importandrew/MaCRa/config"
	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/repo/database"
	"github.com/importandrew/MaCRa/repo/player"
	"github.com/importandrew/MaCRa/test_scripts/user"
)

const (
	defaultConfigPath = "/opt/macra/backend_config.json"
	defaultDBURLEnv   = "PSQL_MACRA"
)

func main() {
	cfg, err := config.Read(defaultConfigPath)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	dbURL := cfg.DBURL
	if dbURL == "" {
		dbURL = os.Getenv(defaultDBURLEnv)
	}

	ctx := context.Background()
	db, err := database.NewDB(ctx, dbURL, false)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	prepo := player.NewRepo(db)

	start := time.Now()

	for i, user := range user.TestUsers() {
		err = prepo.Delete(ctx, user.Email)
		if err != nil {
			logger.Errorf("#%d: %v, %s", i, user, err.Error())
		}
	}

	elapsed := time.Since(start)
	fmt.Println("time taken:", elapsed)
}
