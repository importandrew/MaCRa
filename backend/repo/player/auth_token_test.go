package player

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

// countTokens returns the number of the user's tokens of the purpose: all of
// them, and those not yet used.
func countTokens(t *testing.T, conn *pgxpool.Conn, userID int, purpose models.AuthPurpose) (int, int) {
	t.Helper()

	rows, err := conn.Query(context.Background(), `
		SELECT count(*) AS all_tokens, count(*) FILTER (WHERE used_at IS NULL) AS unused
		FROM auth_tokens WHERE user_id = $1 AND purpose = $2`,
		userID, purpose,
	)
	require.NoError(t, err)

	type counts struct {
		All    int `db:"all_tokens"`
		Unused int `db:"unused"`
	}
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[counts])
	require.NoError(t, err)

	return c.All, c.Unused
}

// ageTokens moves the user's tokens back in time, as if issued d ago.
func ageTokens(t *testing.T, conn *pgxpool.Conn, userID int, d time.Duration) {
	t.Helper()

	_, err := conn.Exec(context.Background(),
		`UPDATE auth_tokens SET created_at = created_at - $1::interval WHERE user_id = $2`,
		strconv.Itoa(int(d.Seconds()))+" seconds", userID)
	require.NoError(t, err)
}

func TestIssueAuthToken(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	expires := time.Now().Add(time.Hour)
	verify, reset := models.AuthPurposeEmailVerify, models.AuthPurposePasswordReset

	db.Test(t, "limits and invalidation", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		defer conn.Release()
		createTestUser(t, conn, 1)

		require.NoError(t, repo.IssueAuthToken(ctx, 1, verify, "h1", expires, models.AuthMailLimits))

		// A second email within the minute is refused, whatever its purpose.
		err = repo.IssueAuthToken(ctx, 1, reset, "h2", expires, models.AuthMailLimits)
		assert.ErrorIs(t, err, models.ErrAuthMailLimit)
		all, _ := countTokens(t, conn, 1, reset)
		assert.Zero(t, all, "a refused token is not stored")

		// A minute later the next one goes out and replaces the unused one of
		// its purpose.
		ageTokens(t, conn, 1, 2*time.Minute)
		require.NoError(t, repo.IssueAuthToken(ctx, 1, verify, "h3", expires, models.AuthMailLimits))
		all, unused := countTokens(t, conn, 1, verify)
		assert.Equal(t, 2, all)
		assert.Equal(t, 1, unused)

		// Five within the hour, one per minute; the sixth is refused.
		for i := range 3 {
			ageTokens(t, conn, 1, 2*time.Minute)
			require.NoError(t, repo.IssueAuthToken(ctx, 1, reset, "r"+strconv.Itoa(i), expires, models.AuthMailLimits))
		}
		ageTokens(t, conn, 1, 2*time.Minute)
		err = repo.IssueAuthToken(ctx, 1, reset, "r-sixth", expires, models.AuthMailLimits)
		assert.ErrorIs(t, err, models.ErrAuthMailLimit)

		// After the hour, the address may get mail again.
		ageTokens(t, conn, 1, time.Hour)
		require.NoError(t, repo.IssueAuthToken(ctx, 1, reset, "r-later", expires, models.AuthMailLimits))
	})

	db.Test(t, "unknown user", func(t *testing.T) {
		err := NewRepo(db).IssueAuthToken(context.Background(), 42, verify, "h", expires, models.AuthMailLimits)
		assert.ErrorIs(t, err, models.ErrPlayerNotFound)
	})

	db.Test(t, "concurrent requests", func(t *testing.T) {
		repo := NewRepo(db)
		ctx := context.Background()
		conn, err := db.Acquire(ctx)
		require.NoError(t, err)
		createTestUser(t, conn, 1)
		conn.Release()

		const attempts = 5
		errs := make(chan error, attempts)
		var wg sync.WaitGroup
		for i := range attempts {
			wg.Go(func() {
				errs <- repo.IssueAuthToken(ctx, 1, reset, "c"+strconv.Itoa(i), expires, models.AuthMailLimits)
			})
		}
		wg.Wait()
		close(errs)

		issued := 0
		for err := range errs {
			if err == nil {
				issued++
				continue
			}
			assert.ErrorIs(t, err, models.ErrAuthMailLimit)
		}
		assert.Equal(t, 1, issued, "one email per minute, also for requests at the same time")
	})
}
