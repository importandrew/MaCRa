package event

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// UpdateTournament sets a tournament's settings, as far as they may still
// change (docs/user-scenarios.md, area 8, "Settings and when they freeze"):
// what happens when time is up only until the first tour starts; the number
// of tours only until the last one starts, and never below the tours started.
// The tour length may change any time: a tour keeps the length it started
// with. It refuses a regular event (models.ErrNotTournament), a finished one
// (models.ErrEventAlreadyFinished) and a frozen setting
// (models.ErrSettingFrozen, with the setting's name).
func (r *Repo) UpdateTournament(ctx context.Context, eventID int, settings models.TournamentSettings) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, err := tx.lockOpenEvent(ctx, eventID); err != nil {
			return err
		}

		rows, err := tx.db.Query(ctx, `
			SELECT e.tournament, e.started_at,
				(SELECT count(*) FROM tours t WHERE t.event_id = e.id AND t.status IN ($2, $3)) AS started_tours
			FROM events e WHERE e.id = $1`,
			eventID, models.TourLive, models.TourFinished)
		if err != nil {
			return err
		}
		type current struct {
			Tournament   *models.TournamentSettings `db:"tournament"`
			StartedAt    *time.Time                 `db:"started_at"`
			StartedTours int                        `db:"started_tours"`
		}
		cur, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[current])
		if err != nil {
			return err
		}
		if cur.Tournament == nil {
			return models.ErrNotTournament
		}
		if err := checkFrozen(*cur.Tournament, settings, cur.StartedAt != nil, cur.StartedTours); err != nil {
			return err
		}

		_, err = tx.db.Exec(ctx, `UPDATE events SET tournament = $1 WHERE id = $2`, settings, eventID)
		return err
	})
}

// checkFrozen returns models.ErrSettingFrozen, with the setting's name, if
// settings change one that may no longer change.
func checkFrozen(old, settings models.TournamentSettings, started bool, startedTours int) error {
	frozen := func(name string) error {
		return fmt.Errorf("%s: %w", name, models.WithData(models.ErrSettingFrozen, models.ErrorData{"setting": name}))
	}
	if started && (settings.TimeUp != old.TimeUp || settings.ExtraHands != old.ExtraHands) {
		return frozen("time_up")
	}
	if settings.Tours != old.Tours && (startedTours >= old.Tours || settings.Tours < startedTours) {
		return frozen("tours")
	}
	return nil
}
