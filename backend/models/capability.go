package models

type Capability string

const (
	CapManageEvent   Capability = "manage_event"
	CapManageMembers Capability = "manage_members"
	CapCreateGame    Capability = "create_game"
	CapAbortGame     Capability = "abort_game"
	CapFinishGame    Capability = "finish_game"
	CapUndoRound     Capability = "undo_round"
	CapWatchLive     Capability = "watch_live"
	CapWatchResults  Capability = "watch_results"
	CapViewRatings   Capability = "view_ratings"

	// Game level (GET /event/:id/game/:gameId/me): what the caller may do with
	// one game.
	CapViewGame   Capability = "view_game"
	CapEnterRound Capability = "enter_round"
	// CapEnterPenalty: enter a penalty or a false mahjong at one's own table.
	CapEnterPenalty Capability = "enter_penalty"
	// CapRefereeEntry: enter into the game as a referee, from its view (by
	// game id), whether or not one plays in it: penalties and false mahjongs.
	CapRefereeEntry Capability = "referee_entry"
)

// EventCapabilities are the capabilities that don't depend on a game, the ones
// GET /event/:id/me lists.
var EventCapabilities = []Capability{
	CapManageEvent, CapManageMembers, CapCreateGame,
	CapWatchLive, CapWatchResults, CapViewRatings,
}

// GameCapabilities are the capabilities about one game, the ones
// GET /event/:id/game/:gameId/me lists.
var GameCapabilities = []Capability{
	CapViewGame, CapEnterRound, CapEnterPenalty, CapRefereeEntry, CapUndoRound, CapAbortGame, CapFinishGame,
}
