// Command devtools drives a test stand's tournament from the command line,
// for checking it by hand: who sits where, random hands at the tables one
// doesn't sit at, the tour's time, test players. It works on the stand's
// database through the backend's services (the engine checks every hand),
// so it needs no player's password. Never point it at production: it
// refuses to run unless MACRA_DEVTOOLS=1.
//
//	MACRA_DEVTOOLS=1 go run ./cmd/devtools -c backend_config.json <command> [flags]
//
// Commands:
//
//	status  -e EVENT                   the current tour's tables: seats, hands, time left
//	play    -e EVENT [-n 3] [-t 2,3] [--except 12]
//	                                   enter N random hands at live tables (all, or -t), skipping
//	                                   the tables where --except (a user ID) sits
//	finish  -e EVENT [--except 12]     play the live tables to their end (and finish a game
//	                                   stopped at its time)
//	time    -e EVENT (--left 2m | --up) set the live tour's time left (--up: the time is up)
//	players -e EVENT -n 8 [--substitutes 2] [--password test1234]
//	                                   create verified test players (and substitutes) in the event
package main

import (
	"context"
	"errors"
	"fmt"
	mrand "math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/importandrew/MaCRa/config"
	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
	erepo "github.com/importandrew/MaCRa/repo/event"
	grepo "github.com/importandrew/MaCRa/repo/game"
	prepo "github.com/importandrew/MaCRa/repo/player"
	gsvc "github.com/importandrew/MaCRa/service/game"
	psvc "github.com/importandrew/MaCRa/service/player"
)

// Random hands: a hand's value is the minimum hand plus up to handSpread;
// one hand in penaltyOdds has a penalty of penaltyValue before it; of ten
// hands, dealInShare are deal-ins, then self-draws up to selfShare, the
// rest draws.
const (
	defaultHands = 3
	handSpread   = 25
	penaltyOdds  = 10
	penaltyValue = 10
	dealInShare  = 6
	selfShare    = 9
	outOfTen     = 10
)

// tools is what the commands work with.
type tools struct {
	db      database.DB
	events  *erepo.Repo
	game    *gsvc.Service
	players *psvc.Service
}

func main() {
	if os.Getenv("MACRA_DEVTOOLS") != "1" {
		failf("refusing to run: set MACRA_DEVTOOLS=1, and only on a test stand, never on production")
	}

	global := pflag.NewFlagSet("devtools", pflag.ExitOnError)
	configPath := global.StringP("config", "c", "/opt/macra/backend_config.json", "path to the backend's config file")
	global.SetInterspersed(false)
	_ = global.Parse(os.Args[1:])
	args := global.Args()
	if len(args) == 0 {
		failf("usage: devtools -c CONFIG status|play|finish|time|players [flags] (see the source's doc comment)")
	}

	ctx := context.Background()
	cfg, err := config.Read(*configPath)
	if err != nil {
		failf("config: %v", err)
	}
	// database.NewDB migrates the database up, as the backend does on start:
	// run the devtools of the same version as the stand's backend.
	db, err := database.NewDB(ctx, cfg.DBURL, false)
	if err != nil {
		failf("database: %v", err)
	}
	defer db.Close()
	playerRepo := prepo.NewRepo(db)
	t := tools{
		db:      db,
		events:  erepo.NewRepo(db),
		game:    gsvc.NewService(grepo.NewRepo(db), playerRepo),
		players: psvc.NewService(playerRepo, nil),
	}

	command, flags := args[0], pflag.NewFlagSet(args[0], pflag.ExitOnError)
	eventID := flags.IntP("event", "e", 0, "the tournament's event ID")
	hands := flags.IntP("hands", "n", defaultHands, "play: hands per table; players: how many players")
	tablesFlag := flags.StringP("tables", "t", "", "play: tables, e.g. 2,3 (default: all live ones)")
	except := flags.Int("except", 0, "play, finish: skip the table where this user ID sits")
	left := flags.Duration("left", 0, "time: time left of the tour, e.g. 2m or 30s")
	up := flags.Bool("up", false, "time: the time is up now")
	substitutes := flags.Int("substitutes", 0, "players: how many substitutes")
	password := flags.String("password", "test1234", "players: the test players' password")
	_ = flags.Parse(args[1:])
	if *eventID == 0 {
		failf("-e EVENT is required")
	}

	switch command {
	case "status":
		err = t.status(ctx, *eventID)
	case "play":
		err = t.play(ctx, *eventID, *hands, parseTables(*tablesFlag), *except, false)
	case "finish":
		err = t.play(ctx, *eventID, 0, nil, *except, true)
	case "time":
		err = t.setTime(ctx, *eventID, *left, *up)
	case "players":
		err = t.addPlayers(ctx, *eventID, *hands, *substitutes, *password)
	default:
		err = fmt.Errorf("unknown command %q", command)
	}
	if err != nil {
		failf("%s: %v", command, err)
	}
}

// failf prints the message and exits with status 1.
func failf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// parseTables reads "2,3" into table numbers.
func parseTables(list string) []int {
	var tables []int
	for part := range strings.SplitSeq(list, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			tables = append(tables, n)
		}
	}
	return tables
}

// sitter returns who sits at a seat now.
func sitter(seat models.Seat) (int, string) {
	if seat.Sitter != nil {
		return seat.Sitter.ID, seat.Sitter.Name
	}
	return seat.UserID, seat.Name
}

// status prints the current tour and its tables.
func (t tools) status(ctx context.Context, eventID int) error {
	state, err := t.events.TourState(ctx, eventID)
	if err != nil {
		return err
	}
	if state.Current == nil {
		fmt.Printf("No tour yet; next: %d, %d players, %d missing.\n", state.Next, state.Active, state.Missing)
		return nil
	}
	tour := state.Current
	fmt.Printf("Tour %d: %s (%s)", tour.Number, tour.Status, tour.Seating)
	if tour.Status == models.TourLive && tour.StartedAt != nil && tour.Minutes != nil {
		left := time.Until(tour.StartedAt.Add(time.Duration(*tour.Minutes) * time.Minute)).Round(time.Second)
		if left > 0 {
			fmt.Printf(", %s left of %d min", left, *tour.Minutes)
		} else {
			fmt.Printf(", time is up (%s ago)", -left)
		}
	}
	fmt.Println()
	for _, table := range state.Tables {
		printTable(table)
	}
	if len(state.Substitutes) > 0 {
		fmt.Print("Substitutes:")
		for _, sub := range state.Substitutes {
			fmt.Printf(" %s [user %d]", sub.Name, sub.ID)
			if sub.Table > 0 {
				fmt.Printf(" at table %d", sub.Table)
			}
			fmt.Print(";")
		}
		fmt.Println()
	}
	return nil
}

// play enters random hands at the live tour's tables: hands per table, or,
// with toEnd, until each game ends (finishing one stopped at its time).
// Tables are all live ones, or those listed; the table where except sits
// is skipped.
func (t tools) play(ctx context.Context, eventID, hands int, tables []int, except int, toEnd bool) error {
	state, err := t.events.TourState(ctx, eventID)
	if err != nil {
		return err
	}
	if state.Current == nil || state.Current.Status != models.TourLive {
		return errors.New("no tour is live: start it on the desk first")
	}
	rng := mrand.New(mrand.NewPCG(uint64(time.Now().UnixNano()), 1)) //nolint:gosec // test data
	for _, table := range state.Tables {
		if table.Status != models.GameLive || (len(tables) > 0 && !slices.Contains(tables, table.Table)) {
			continue
		}
		if slices.ContainsFunc(table.Seats, func(s models.Seat) bool { id, _ := sitter(s); return id == except }) {
			fmt.Printf("table %d: skipped (user %d sits there)\n", table.Table, except)
			continue
		}
		entered, err := t.playTable(ctx, eventID, table, hands, toEnd, rng)
		fmt.Printf("table %d: %d hands entered", table.Table, entered)
		if err != nil {
			fmt.Printf(", stopped: %v", err)
		}
		fmt.Println()
	}
	return t.status(ctx, eventID)
}

// maxHands bounds the hands play enters at a table, so a bug can't loop.
const maxHands = 100

// playTable enters random hands at one table as whoever sits at its east
// seat now (any sitter may enter the table's hands).
func (t tools) playTable(ctx context.Context, eventID int, table models.TourTable, hands int, toEnd bool, rng *mrand.Rand) (int, error) {
	entrant, _ := sitter(table.Seats[0])
	entered := 0
	for (toEnd || entered < hands) && entered < maxHands {
		now, err := t.game.Now(ctx, entrant, eventID)
		switch {
		case errors.Is(err, models.ErrGameNotFound):
			return entered, nil // the game has ended
		case err != nil:
			return entered, err
		case now.Status != models.GameLive:
			return entered, nil
		case now.Timer != nil && now.Timer.Stopped:
			// The time is up and the rule ends the game at once.
			return entered, t.game.Finish(ctx, now.ID)
		}
		if _, err := t.randomHand(ctx, eventID, entrant, now, rng); err != nil {
			if errors.Is(err, models.ErrTimeUp) {
				return entered, t.game.Finish(ctx, now.ID)
			}
			return entered, err
		}
		entered++
	}
	return entered, nil
}

// randomHand enters one random hand: mostly a deal-in or a self-draw worth
// 8 to 32 above the minimum, sometimes a draw, now and then a penalty
// before it.
func (t tools) randomHand(ctx context.Context, eventID, entrant int, now models.GameNow, rng *mrand.Rand) (models.GameNow, error) {
	seats := []int{now.Me.ID, now.Right.ID, now.Opp.ID, now.Left.ID}
	winner := seats[rng.IntN(len(seats))]
	loser := winner
	for loser == winner {
		loser = seats[rng.IntN(len(seats))]
	}
	value := now.Rules.MinHand + rng.IntN(handSpread)

	if rng.IntN(penaltyOdds) == 0 {
		var err error
		if now, err = t.game.Penalty(ctx, models.Penalty, entrant, loser, penaltyValue, false, eventID, now.LogSize); err != nil {
			return now, fmt.Errorf("penalty: %w", err)
		}
	}
	switch roll := rng.IntN(outOfTen); {
	case roll < dealInShare:
		return t.game.EnterRound(ctx, entrant, models.Deal, value, winner, loser, eventID, now.LogSize)
	case roll < selfShare:
		return t.game.EnterRound(ctx, entrant, models.Self, value, winner, 0, eventID, now.LogSize)
	default:
		return t.game.EnterRound(ctx, entrant, models.Draw, 0, 0, 0, eventID, now.LogSize)
	}
}

// setTime moves the live tour's start so that left remains of its time (at
// a table without added minutes), or so that it is up now.
func (t tools) setTime(ctx context.Context, eventID int, left time.Duration, up bool) error {
	if !up && left <= 0 {
		return errors.New("give --left 2m or --up")
	}
	if up {
		left = 0
	}
	res, err := t.db.Querier().Exec(ctx, `
		UPDATE tours SET started_at = $2::timestamptz - make_interval(mins => minutes) + make_interval(secs => $3::float8)
		WHERE event_id = $1 AND status = $4`,
		eventID, time.Now().UTC(), left.Seconds(), models.TourLive)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no tour is live")
	}
	fmt.Printf("The live tour's time: %s left (tables with added time get theirs on top).\n", left)
	if up {
		fmt.Println("With \"the game ends at once\", the backend finishes the stopped games within about 15 seconds.")
	}
	return t.status(ctx, eventID)
}

// addPlayers creates verified test players and substitutes and adds them
// to the event with their role.
func (t tools) addPlayers(ctx context.Context, eventID, players, substitutes int, password string) error {
	stamp := time.Now().Format("0102150405")
	add := func(i int, role models.EventRole, label string) error {
		email := fmt.Sprintf("t%d-%s-%s%d@test.local", eventID, stamp, role, i)
		name := fmt.Sprintf("%s %d", label, i)
		player, err := t.players.Add(ctx, email, name, password)
		if err != nil {
			return fmt.Errorf("%s: %w", email, err)
		}
		if _, err := t.db.Querier().Exec(ctx, `UPDATE users SET email_verified_at = $2 WHERE id = $1`,
			player.ID, time.Now().UTC()); err != nil {
			return err
		}
		if err := t.events.AddMemberRoles(ctx, eventID, player.ID, []models.EventRole{role}); err != nil {
			return fmt.Errorf("%s: %w", email, err)
		}
		fmt.Printf("%s [user %d] %s\n", name, player.ID, email)
		return nil
	}
	for i := 1; i <= players; i++ {
		if err := add(i, models.EventPlayer, "Test player"); err != nil {
			return err
		}
	}
	for i := 1; i <= substitutes; i++ {
		if err := add(i, models.EventSubstitute, "Test sub"); err != nil {
			return err
		}
	}
	fmt.Printf("Password: %s\n", password)
	return nil
}

// printTable prints a table's game and seats, who sits there now first.
func printTable(table models.TourTable) {
	fmt.Printf("  table %d (game %d, %s, %d log entries", table.Table, table.GameID, table.Status, table.LogSize)
	if table.ExtraMinutes > 0 {
		fmt.Printf(", +%d min", table.ExtraMinutes)
	}
	fmt.Println("):")
	for _, seat := range table.Seats {
		id, name := sitter(seat)
		note := ""
		if seat.Sitter != nil {
			note = fmt.Sprintf(" (for %s)", seat.Name)
		}
		if seat.Forfeit {
			note += " 0/0"
		}
		fmt.Printf("    %-5s #%-3d %s [user %d]%s\n", seat.Wind, seat.Number, name, id, note)
	}
}
