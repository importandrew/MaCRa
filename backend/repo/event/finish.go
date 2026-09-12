package event

import (
	"context"
	"time"
)

func (r *Repo) Finish(ctx context.Context, eventID int) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		UPDATE events
		SET finished_at = $1
		WHERE id = $2`,
		time.Now().UTC(), eventID,
	)
	if err != nil {
		return err
	}

	return nil
}
