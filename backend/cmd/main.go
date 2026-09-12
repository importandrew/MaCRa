package main

import (
	"context"

	"github.com/importandrew/MaCRa/config"
	"github.com/importandrew/MaCRa/handler"
	ehndl "github.com/importandrew/MaCRa/handler/event"
	phndl "github.com/importandrew/MaCRa/handler/player"
	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/repo/database"
	erepo "github.com/importandrew/MaCRa/repo/event"
	grepo "github.com/importandrew/MaCRa/repo/game"
	prepo "github.com/importandrew/MaCRa/repo/player"
	"github.com/importandrew/MaCRa/server"
	esvc "github.com/importandrew/MaCRa/service/event"
	gsvc "github.com/importandrew/MaCRa/service/game"
	psvc "github.com/importandrew/MaCRa/service/player"
	"github.com/spf13/pflag"
)

const (
	defaultConfigPath = "/opt/macra/backend_config.json"
)

func main() {
	ctx := context.Background()

	var configPath string
	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "path to config file")
	pflag.Parse()

	cfg, err := config.Read(configPath)
	if err != nil {
		logger.Fatalf("cfg read: %v\n", err)
	}

	if err := cfg.Validate(); err != nil {
		logger.Fatalf("cfg validate: %v\n", err)
	}

	// db
	db, err := database.NewDB(ctx, cfg.DBURL, cfg.DBLog)
	if err != nil {
		logger.Fatalf("new db: %v\n", err)
	}
	defer db.Close()

	// repo
	gameRepo := grepo.NewRepo(db)
	playerRepo := prepo.NewRepo(db)
	eventRepo := erepo.NewRepo(db)

	// service
	gameService := gsvc.NewService(gameRepo, playerRepo)
	playerService := psvc.NewService(playerRepo)
	eventService := esvc.NewService(eventRepo, gameService, playerService)

	// handler
	eventHandler := ehndl.NewHandler(eventService, gameService)
	playerHandler := phndl.NewHandler(playerService, cfg.CookieDomain, []byte(cfg.JWTKey))
	h := handler.New(eventHandler, playerHandler, []byte(cfg.JWTKey))

	// server
	router := server.NewServer(cfg.AllowOrigins...)
	h.RegisterRoutes(router)

	if cfg.CertFile != "" {
		router.RunTLS(cfg.Addr, cfg.CertFile, cfg.KeyFile)
	}
	router.Run(cfg.Addr)
}
