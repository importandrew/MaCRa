package player

import (
	"github.com/importandrew/MaCRa/repo/database"
)

type Repo struct {
	db database.DB
}

func NewRepo(db database.DB) *Repo {
	return &Repo{db: db}
}
