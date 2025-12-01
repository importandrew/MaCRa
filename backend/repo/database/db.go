package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	logger "github.com/importandrew/MaCRa/logger"
	pgxlogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
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
