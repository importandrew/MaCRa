package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/game"
	"github.com/importandrew/MaCRa/handler/player"
)

type Handler struct {
	game   game.Handler
	player player.Handler
	authmw authmw.Auth
}

func New(game game.Handler, player player.Handler, jwt []byte) Handler {
	return Handler{
		game:   game,
		player: player,
		authmw: authmw.NewAuth(jwt),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// user auth management
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.player.SignUp)
		auth.POST("/login", h.player.Login)
		auth.POST("/logout", h.player.Logout)
	}

	// public handlers
	publicAPI := router.Group("/api")
	{
		publicAPI.GET("/game/:id/players", h.game.Players)
		publicAPI.GET("/game/:id/logs", h.game.Logs)
		publicAPI.GET("/game/:id/hud", h.game.HUD)
		publicAPI.GET("/game/:id/live", h.game.ViewLive)
		publicAPI.GET("/game/:id/result", h.game.FinalResult)
		publicAPI.GET("/game/:id/status", h.game.Status)

		publicAPI.GET("/players/list", h.player.List)
	}

	// authorized only handlers
	authAPI := publicAPI.Group("/", h.authmw.IdentifyUser)
	{
		authAPI.DELETE("/game/:id/abort", h.game.Abort)
		authAPI.POST("/game/:id/finish", h.game.Finish)

		// auto-detects user who called this api
		{
			authAPI.POST("/new", h.game.New)
			authAPI.GET("/now", h.game.Now)
			authAPI.POST("/round", h.game.EnterRoundResult)
			authAPI.POST("/false-mahjong", h.game.FalseMahjong)
			authAPI.POST("/penalty", h.game.Penalty)
			authAPI.GET("/last", h.game.Last)

			authAPI.GET("/profile", h.player.GetProfile)
			authAPI.PUT("/profile", h.player.UpdateProfile)

			// won't abort on 401
			authAPI.GET("/main", h.game.MainInfo)
		}
	}
}
