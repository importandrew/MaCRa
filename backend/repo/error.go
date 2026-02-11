package repo

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsTargetErr(targetSQLState int, err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == strconv.Itoa(targetSQLState)
	}

	return false
}
