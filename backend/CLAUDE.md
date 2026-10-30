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
- `models/` holds shared types and **all public sentinel errors**. `models.Errors` maps each error to an HTTP status; `handler/httperr.Error` uses `errors.Is` against that map (default 500). To return a new client-facing status, add an error there and wrap it with `%w`.
- All routes are in `handler/handler.go`. `/auth/*` is public; `/api/*` runs `authmw.SetUserID` (JWT cookie, optional). `middleware.Require*` and `reqctx` parse and fetch path params (`eventId`, `gameId`, `playerId`) and the user ID.
- **Authorization is capability-based per event.** `service/event` builds an auth context (owner, member roles, game participation, event visibility/reach settings) and checks `models.Capability` values via `HasCapability`. Handlers call it before acting.
- **Game state is event-sourced.** No in-memory game state survives across requests. `service/game.load` rebuilds an `engine/game.Game` from the starting seating plus the stored round log (`RoundLog` rows replayed through `EnterRoundResult` / `Penalty` / `FalseMahjong`). Mutations apply to the engine, then persist the new round row. When the engine reports `GameEnded`, the service finalizes results. Undo = delete the last round row.
- DB migrations are `repo/database/migrations/NNNNNN_*.{up,down}.sql` (golang-migrate), embedded in the binary and applied automatically on startup by `database.NewDB`. Add new ones with the next sequence number and always provide a `down`.
