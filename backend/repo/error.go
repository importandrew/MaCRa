package repo

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL SQLSTATE codes checked with IsTargetErr.
const (
	SQLStateStringTooLong   = 22001 // string_data_right_truncation
	SQLStateUniqueViolation = 23505 // unique_violation
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
