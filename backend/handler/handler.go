package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/event"
	"github.com/importandrew/MaCRa/handler/middleware"
	"github.com/importandrew/MaCRa/handler/player"
)

type Handler struct {
	event  event.Handler
	player player.Handler
	authmw middleware.Auth
}

func New(eventHandler event.Handler, player player.Handler, jwt []byte) Handler {
	return Handler{
		event:  eventHandler,
		player: player,
		authmw: middleware.NewAuth(jwt),
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

	// main api
	api := router.Group("/api", h.authmw.SetUserID)
	{
		api.GET("/event", h.event.ListVisible)

		eventAccess := api.Group("/event/:eventId", middleware.RequireEvent, h.event.HasAccess)
		{
			eventAccess.GET("", h.event.Get)
			eventAccess.GET("/main", h.event.MainInfo)
			eventAccess.GET("/me", h.event.GetMe)
			eventAccess.GET("/ratings", h.event.Ratings)
			eventAccess.GET("/games", h.event.ListGames)

			eventGame := eventAccess.Group("/game/:gameId", middleware.RequireGame, h.event.GameBelongsToEvent)
			{
				eventGame.GET("/live", h.event.GameViewLive)
				eventGame.GET("/result", h.event.GameResult)
				eventGame.GET("/players", h.event.GamePlayers)
				eventGame.GET("/logs", h.event.GameLogs)
				eventGame.GET("/hud", h.event.GameHUD)
				eventGame.GET("/status", h.event.GameStatus)
			}
		}

		// required auth
		required := api.Group("", middleware.RequireUser)
		{
			required.GET("/players", h.player.List)
			required.GET("/profile", h.player.GetProfile)
			required.PUT("/profile", h.player.UpdateProfile)

			required.POST("/event", h.event.Create)

			event := api.Group("/event/:eventId", middleware.RequireEvent, h.event.HasAccess)
			{
				event.PATCH("", h.event.Update)
				event.POST("/finish", h.event.FinishEvent)

				event.GET("/members", h.event.ListMembers)
				event.POST("/members", h.event.AddMembers)
				event.DELETE("/members", h.event.RemoveMember)

				event.GET("/game-history", h.event.GameHistory)
				event.GET("/player/:playerId/game-history", middleware.RequirePlayer, h.event.PlayerGameHistory)

				// user relative game block
				{
					event.POST("/new", h.event.NewGame)
					event.GET("/now", h.event.Now)
					event.GET("/last", h.event.Last)
					event.POST("/round", h.event.EnterRoundResult)
					event.POST("/penalty", h.event.Penalty)
					event.POST("/false-mahjong", h.event.FalseMahjong)
				}

				// explicit game block
				eventGame := event.Group("/game/:gameId", middleware.RequireGame, h.event.GameBelongsToEvent)
				{
					eventGame.DELETE("/abort", h.event.Abort)
					eventGame.DELETE("/undo-round", h.event.UndoLastRound)
					eventGame.POST("/finish", h.event.FinishGame)
				}
			}
		}
	}
}
