package database

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimeColumns checks that every time column is timestamptz without a
// default: the backend sets each value itself, in UTC.
func TestTimeColumns(t *testing.T) {
	db, err := NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "schema", func(t *testing.T) {
		type column struct {
			Table   string  `db:"table_name"`
			Column  string  `db:"column_name"`
			Type    string  `db:"data_type"`
			Default *string `db:"column_default"`
		}

		rows, err := db.Querier().Query(context.Background(), `
			SELECT table_name::text, column_name::text, data_type::text, column_default::text
			FROM information_schema.columns
			WHERE table_schema = current_schema() AND data_type LIKE 'timestamp%'`)
		require.NoError(t, err)
		columns, err := pgx.CollectRows(rows, pgx.RowToStructByName[column])
		require.NoError(t, err)
		require.NotEmpty(t, columns)

		for _, c := range columns {
			assert.Equal(t, "timestamp with time zone", c.Type, "%s.%s", c.Table, c.Column)
			assert.Nil(t, c.Default, "%s.%s", c.Table, c.Column)
		}
	})

	db.Test(t, "read in UTC", func(t *testing.T) {
		rows, err := db.Querier().Query(context.Background(),
			`SELECT '2026-01-02 03:04:05+03'::timestamptz`)
		require.NoError(t, err)
		got, err := pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
		require.NoError(t, err)

		assert.Equal(t, time.UTC, got.Location())
		assert.True(t, got.Equal(time.Date(2026, 1, 2, 0, 4, 5, 0, time.UTC)))
	})
}
