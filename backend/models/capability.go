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
)
