package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo"
)

func (r *Repo) Update(ctx context.Context, id int, name string) error {
	res, err := r.db.Exec(ctx, "UPDATE users SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		// Same mapping as Add: a name the column can't hold is a client error.
		if repo.IsTargetErr(repo.SQLStateStringTooLong, err) {
			return models.ErrLongName
		}
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrPlayerNotFound
	}

	return nil
}
