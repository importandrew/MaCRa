package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	pgxlogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	logger "github.com/importandrew/MaCRa/logger"
)

const defaultDBURLEnv = "PSQL_MACRA"
const testDBURLEnv = defaultDBURLEnv + "_TEST"

type DB struct {
	pool *pgxpool.Pool
	mig  *migrator
}

func NewDB(ctx context.Context, dbURL string, dbLog bool) (DB, error) {
	if dbURL == "" {
		dbURL = os.Getenv(defaultDBURLEnv)
	}

	mig, err := NewMigrator(dbURL)
	if err != nil {
		return DB{}, fmt.Errorf("migrator: %w", err)
	}

	if err := mig.Up(); err != nil {
		return DB{}, fmt.Errorf("migrate: %w", err)
	}

	dbConf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return DB{}, fmt.Errorf("cannot parse db config url: %w", err)
	}

	if dbLog {
		pgxlogger := pgxlogrus.NewLogger(logger.StandardLogger())
		dbConf.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxlogger,
			LogLevel: tracelog.LogLevelTrace,
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		return DB{}, fmt.Errorf("pgxpool new: %w", err)
	}

	return DB{pool: pool, mig: mig}, nil
}

func NewTestDB() (DB, error) {
	return NewDB(context.Background(), os.Getenv(testDBURLEnv), false)
}

func (db *DB) Close() {
	db.pool.Close()
	db.mig.mig.Close()
}

func (db *DB) Acquire(ctx context.Context) (*pgxpool.Conn, error) { return db.pool.Acquire(ctx) }

// Querier runs SQL on the pool or inside a transaction: *pgxpool.Pool and
// pgx.Tx both implement it. On the pool, each statement takes a connection
// only until its rows are read. Begin inside a transaction starts a savepoint.
type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Querier returns the pool as a Querier, for repos to run their queries on.
func (db *DB) Querier() Querier { return db.pool }

// InTx begins a transaction on q (a savepoint if q is a transaction), calls
// fn with it and commits it if fn returns nil. Otherwise the transaction is
// rolled back and fn's error returned.
func InTx(ctx context.Context, q Querier, fn func(tx pgx.Tx) error) error {
	tx, err := q.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Test is wrapper for t.Run() which migrates db up before running the test
// and down after it finishes.
func (db *DB) Test(t *testing.T, name string, f func(t *testing.T)) {
	if err := db.mig.Up(); err != nil {
		t.Error(err)
	}

	defer func() {
		if err := db.mig.Down(); err != nil {
			t.Error(err)
		}
	}()

	t.Run(name, f)
}
