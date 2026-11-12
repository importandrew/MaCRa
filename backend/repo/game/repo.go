package game

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/repo/database"
)

// Repo runs its queries on db: the pool, or one transaction for a Repo that
// InTx passes to its callback.
type Repo struct {
	db database.Querier
}

func NewRepo(db database.DB) *Repo {
	return &Repo{db: db.Querier()}
}

// InTx calls fn with a Repo whose methods all run in one transaction, and
// commits it if fn returns nil. Inside fn, call methods on tx, not on r: a
// method of r runs outside the transaction.
func (r *Repo) InTx(ctx context.Context, fn func(tx *Repo) error) error {
	return database.InTx(ctx, r.db, func(tx pgx.Tx) error {
		return fn(&Repo{db: tx})
	})
}
