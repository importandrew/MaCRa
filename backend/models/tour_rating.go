package models

// TourRating is a tournament's rating by tour (docs/user-scenarios.md, area
// 8, "Rating"): the finished tours in order, the tour in progress (live; 0
// when none), and the players: the active ones by place, then those who
// withdrew.
type TourRating struct {
	Tours      []int           `json:"tours"`
	InProgress int             `json:"in_progress"`
	Rows       []TourRatingRow `json:"rows"`
}

// TourRatingRow is a player in the rating by tour: their place (nil for a
// player who withdrew or hasn't played a finished tour; players equal in TP
// and GP share it), totals, and their result in each finished tour (aligned
// with TourRating.Tours; nil when they didn't play it).
type TourRatingRow struct {
	PlayerID int    `json:"player_id"`
	Name     string `json:"name"`
	// Number is the participant number; nil while not drawn.
	Number    *int          `json:"number"`
	Withdrawn bool          `json:"withdrawn"`
	Rank      *int          `json:"rank"`
	TP        float32       `json:"tp"`
	GP        int           `json:"gp"`
	Tours     []*TourResult `json:"tours"`
}

// TourResult is a player's result in one tour, and their table.
type TourResult struct {
	TP    float32 `json:"tp"`
	GP    int     `json:"gp"`
	Table int     `json:"table"`
}

// TourResultRow is a player's result in a finished tour, as the repo reads
// it.
type TourResultRow struct {
	Tour     int     `db:"tour"`
	PlayerID int     `db:"player_id"`
	Name     string  `db:"name"`
	TP       float32 `db:"tp"`
	GP       int     `db:"gp"`
	Table    int     `db:"table_number"`
}
