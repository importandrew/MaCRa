package models

type GamePlayers struct {
	ID      int      `json:"id"`
	Players []Player `json:"players"`
}

type Player struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type StartingPlayers struct {
	East  int `db:"east" json:"east"`
	South int `db:"south" json:"south"`
	West  int `db:"west" json:"west"`
	North int `db:"north" json:"north"`
}
