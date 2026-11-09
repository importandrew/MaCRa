# Backend

## Running and testing

- Run: `go run ./cmd/main.go [-c path/to/backend_config.json]`. The default config path is `/opt/macra/backend_config.json`; see `config/backend_config.json.example`.
- Config: `jwt_key` and `mail_from` are required (`Validate` fails otherwise). If `db_url` is empty, the `PSQL_MACRA` env var is used. Mail is sent via SMTP on `127.0.0.1:25` without auth.
- Regenerate mocks with `mockery`. It reads `.mockery.yml` and writes `mocks_test.go` next to each listed interface, so add new interfaces there.
- Repo tests (`repo/...`) need a real PostgreSQL in `PSQL_MACRA_TEST`. `database.DB.Test` runs migrations **up before and down after** each test, so the test DB is wiped. Never point it at real data. Service and handler tests use mocks and need no DB.
- Use `go test -count=1 ./...` to bypass the test cache when DB state matters.

## Architecture

Layered, wired manually in `cmd/main.go`: `repo` (pgx SQL) → `service` (business logic) → `handler` (gin) → `server` (gin engine + CORS).

- Each layer declares the narrow interfaces it consumes (e.g. `service/game.GameRepo`, `handler/event.GameService`) in the consumer package. Those interfaces are what mockery mocks.
- Repo code reads results with `Query` + `pgx.CollectOneRow` / `pgx.CollectRows`, never `QueryRow(...).Scan(...)`. Use `pgx.RowTo[T]` for a single column; for several columns, use `pgx.RowToStructByName` with a `models` type or a local struct with `db` tags (e.g. `roundInsertResult` in `repo/game/save_round.go`). `CollectOneRow` returns `pgx.ErrNoRows` for no rows; map it to a `models` sentinel error where the caller needs one (e.g. `ErrGameNotFound` in `repo/game/current.go`).
- A repo method that holds a connection or a transaction must not call another `Repo` method that acquires its own (e.g. `Players`, `SetStatus`): run the query on the connection or `tx` it holds (e.g. `gamePlayers(ctx, tx, …)`). Otherwise the work isn't in the transaction, and under load every pool connection is held by a request waiting for a second one, so they all hang.
- `models/` holds shared types and **all public sentinel errors**. `models.Errors` maps each error to an HTTP status; `handler/httperr.Error` uses `errors.Is` against that map (default 500). To return a new client-facing status, add an error there and wrap it with `%w`.
- All routes are in `handler/handler.go`. `/auth/*` is public; `/api/*` runs `authmw.SetUserID` (JWT cookie: none → anonymous; invalid or expired → cookie cleared and 401 `models.ErrSessionExpired`; older than a day → reissued). `middleware.Require*` and `reqctx` parse and fetch path params (`eventId`, `gameId`, `playerId`) and the user ID.
- **Authorization is capability-based per event.** `service/event` builds an auth context (owner, member roles, game participation, event visibility/reach settings) and checks `models.Capability` values via `HasCapability`. Handlers call it before acting.
- **Game state is event-sourced.** No in-memory game state survives across requests. `service/game.load` rebuilds an `engine/game.Game` from the starting seating plus the stored round log (`RoundLog` rows replayed through `EnterRoundResult` / `Penalty` / `FalseMahjong`). Mutations apply to the engine, then persist the new round row. When the engine reports `GameEnded`, the service finalizes results. Undo = delete the last round row (`repo/game.UndoLastRound`, one transaction; for an ended game it also deletes the results and makes the game live again).
- **A round holds one main outcome and any number of penalties.** Draw, self-draw or deal-in ends the round; penalties and false mahjongs can come before it, any number of times, and one with `force_next_round` ends the round without a main outcome. So a `rounds` row is a log entry, not a round: count entries (`LogSize`), not round numbers. New entries and undos carry the log size the client saw, and `repo/game.SaveRoundOutcome` / `UndoLastRound` refuse them (`ErrLogChanged`) unless the game still has that log size, checked under a lock on the game row (`lockUnchangedLog`).
- DB migrations are `repo/database/migrations/NNNNNN_*.{up,down}.sql` (golang-migrate), embedded in the binary and applied automatically on startup by `database.NewDB`. Add new ones with the next sequence number and always provide a `down`.

## One live game per player per event

- Many endpoints take no game ID. They resolve the game from the requesting user and the event: `repo/game.CurrentGame` returns that player's newest `live` game in the event. These are `POST /new` (refuses a player who already has one), `GET /now`, `POST /round`, `/penalty` and `/false-mahjong`, plus the member-removal check (`ErrMemberInCurrentGame`). `GET /last` is resolved the same way but returns the player's most recently finished (`ended`) game. The frontend pages built on them (`/now`, `/last`, the draft) assume one live game per player.
- So "at most one live game per player per event" is an invariant that the code relies on, not just a validation. The schema can't enforce it (status and event live in `games`, players in `game_players`). `repo/game.NewGame` enforces it under `pg_advisory_xact_lock(event_id, player_id)` and checks again inside the transaction; the service-level check alone lost a race to a double tap. Anything else that creates or revives a `live` game must take the same locks: undoing the final round does (`reviveGame` in `repo/game/undo_last_round.go`, `ErrUndoFinalRound` if a player has started another game).
- A finished event has no live games. `repo/event.Finish` locks the event row `FOR UPDATE` before checking for live games; whatever creates or revives a live game locks it `FOR SHARE` first and refuses a finished event (`lockOpenEvent` in `repo/game/new.go`, `ErrEventAlreadyFinished`).
- Endpoints under `/game/:gameId` (abort, undo, finish, live, logs, …) work on an explicit game and don't depend on this.
