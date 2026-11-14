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
		auth.POST("/verify-email", h.player.VerifyEmail)
		auth.POST("/resend-verification", h.player.ResendVerification)
		auth.POST("/forgot-password", h.player.ForgotPassword)
		auth.POST("/reset-password", h.player.ResetPassword)
	}

	// main api
	api := router.Group("/api", h.authmw.SetUserID)
	{
		api.GET("/event", h.event.ListVisible)
		api.GET("/invite/:token", h.event.Invite)

		eventAccess := api.Group("/event/:eventId", middleware.RequireEvent, h.event.HasAccess)
		{
			eventAccess.GET("", h.event.Get)
			eventAccess.GET("/main", h.event.MainInfo)
			eventAccess.GET("/me", h.event.GetMe)
			eventAccess.GET("/ratings", h.event.Ratings)
			eventAccess.GET("/games", h.event.ListGames)
			eventAccess.GET("/members", h.event.ListMembers)
			eventAccess.GET("/game-history", h.event.GameHistory)
			eventAccess.GET("/player/:playerId/game-history", middleware.RequirePlayer, h.event.PlayerGameHistory)
			eventAccess.GET("/player/:playerId/stats", middleware.RequirePlayer, h.event.PlayerStats)

			eventGame := eventAccess.Group("/game/:gameId", middleware.RequireGame, h.event.GameBelongsToEvent)
			{
				eventGame.GET("/live", h.event.GameViewLive)
				eventGame.GET("/result", h.event.GameResult)
				eventGame.GET("/players", h.event.GamePlayers)
				eventGame.GET("/logs", h.event.GameLogs)
				eventGame.GET("/hud", h.event.GameHUD)
				eventGame.GET("/status", h.event.GameStatus)
				eventGame.GET("/me", h.event.GameMe)
			}
		}

		// required auth
		required := api.Group("", middleware.RequireUser)
		{
			required.GET("/players", h.player.List)
			required.GET("/profile", h.player.GetProfile)
			required.PUT("/profile", h.player.UpdateProfile)

			required.POST("/event", h.event.Create)
			required.POST("/invite/:token/join", h.event.JoinByInvite)

			h.registerEventActions(required.Group("/event/:eventId", middleware.RequireEvent, h.event.HasAccess))
		}
	}
}

// registerEventActions registers what a signed-in user does in an event: its
// settings and members, their own game (resolved from the user) and games
// by id.
func (h *Handler) registerEventActions(event *gin.RouterGroup) {
	event.PATCH("", h.event.Update)
	event.PUT("/rules", h.event.UpdateRules)
	event.PUT("/settings", h.event.UpdateSettings)
	event.PUT("/tournament", h.event.UpdateTournament)
	event.POST("/finish", h.event.FinishEvent)

	event.POST("/members", h.event.AddMembers)
	event.DELETE("/members", h.event.RemoveMember)
	event.PUT("/members/:playerId", middleware.RequirePlayer, h.event.SetMemberRoles)
	event.GET("/invite", h.event.InviteToken)
	event.POST("/invite", h.event.NewInviteToken)
	event.DELETE("/invite", h.event.TurnOffInvite)
	event.PUT("/joining", h.event.UpdateJoining)
	event.GET("/requests", h.event.JoinRequests)
	event.POST("/requests/:playerId/accept", middleware.RequirePlayer, h.event.AcceptJoinRequest)
	event.POST("/requests/:playerId/decline", middleware.RequirePlayer, h.event.DeclineJoinRequest)
	event.POST("/request", h.event.RequestToJoin)
	event.DELETE("/request", h.event.WithdrawJoinRequest)

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
		eventGame.POST("/penalty", h.event.GamePenalty)
		eventGame.POST("/false-mahjong", h.event.GameFalseMahjong)
		eventGame.POST("/finish", h.event.FinishGame)
	}
}
