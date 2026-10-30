# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

MaCRa (Mahjong Competitive Rules assistant) tracks MCR mahjong games: live scores viewable from any device, game logs, ratings, and events that users host. Monorepo with a Go backend (`backend/`, module `github.com/importandrew/MaCRa`) and a React/Vite frontend (`frontend/`).

## Working rules

- Start every task on a new branch cut from a freshly fetched master: `git fetch origin && git checkout -b feature/<short-name> origin/master`. Don't let the branch track `origin/master`; push with `git push -u origin feature/<short-name>`.
- Never push to `master` and never merge PRs. Only the repository owner merges.
- No force pushes (`--force`, `-f`, `--force-with-lease`, `+refspec`). Don't commit to, push to or rewrite branches you didn't create.
- Before every commit, run the backend tests and `golangci-lint run`, plus `npm run lint` and `npm run build` if `frontend/` changed. Don't add new failures or lint issues compared with master (see "Known baseline" below).
- Don't read or commit `.env*` files, keys, passwords or tokens. That includes `backend_config.json` with real values (`db_url`, `jwt_key`).
- Ask the owner before changing DB migrations (`backend/repo/database/migrations/`), CI (`.github/`), or dependencies (`go.mod`/`go.sum`, `package.json`/`package-lock.json`).
- Commit messages follow Conventional Commits: `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, `docs:`, … (as in the existing history).

`.claude/settings.json` enforces some of these rules with permission rules. It denies force pushes, `gh pr merge` and reading `.env*`, and it asks before `git push` and `gh pr create`. Pattern rules can be bypassed by unusual command spellings, so the rules above still apply in full. Put personal overrides in `.claude/settings.local.json`, which is gitignored.

## Commands

### Backend (run from `backend/`)
- Run: `go run ./cmd/main.go [-c path/to/backend_config.json]` (default config path `/opt/macra/backend_config.json`; see `config/backend_config.json.example`)
- Build: `make build` → `dist/macra-backend`
- Lint: `golangci-lint run` (strict v2 config in `.golangci.yml`)
- Regenerate mocks: `mockery` (config in `.mockery.yml`; writes `mocks_test.go` next to each listed interface — add new interfaces there)
- All tests: `go test ./...`
- Single test: `go test ./service/game -run TestEnterRound`

Use `go test -count=1 ./...` to bypass the test cache when DB state matters.

Repo tests (`repo/...`) need a real PostgreSQL in `PSQL_MACRA_TEST`. `database.DB.Test` runs migrations **up before and down after** each test, so the test DB is wiped — never point it at real data. Service and handler tests use mockery mocks and need no DB.

Config: `jwt_key` and `mail_from` are required (`Validate` fails otherwise). If `db_url` is empty, `PSQL_MACRA` env is used. Mail is sent via SMTP on `127.0.0.1:25` (no auth).

### Frontend (run from `frontend/`)
- `npm ci` (install from the lockfile without modifying it), `npm run dev` (port 3000), `npm run build` (`tsc -b && vite build`), `npm run build-dev` (`--mode devprod`)
- `npm run lint` / `npm run lint:fix` (ESLint + Prettier)
- No frontend test suite. CI (`.github/workflows/frontend.yml`) runs only frontend lint + build on Node 18.
- API base URL comes from `VITE_API_ADDR` in `.env.<mode>`; for local dev the backend's `allow_origins`/`cookie_domain` must match the frontend origin.

### Known baseline (master @ 7082fc5, 2026-10-30)
All of the commands above work in the dev container (Go 1.24, golangci-lint 2.6.2, Node 22). Results on master:
- `go test ./...`: everything passes except `handler/player` `TestNormalizeOrigin` and `TestHandler_SignUp`. `normalizeOrigin` accepts only configured `allow_origins`, but the tests expect `https://macra.pro` to be allowed implicitly, as CORS does in `server.NewServer`.
- `golangci-lint run`: 27 existing issues (gofmt, godot, mnd, nilerr, predeclared, testifylint, unparam, …).
- `npm run lint`: 0 errors, 240 warnings. `npm run build` passes.

Update this section when the baseline changes.

Root `Makefile` builds both. `deploy*` targets scp to a server via `MACRA_SSHKEY`/`MACRA_HOST` and restart systemd services — production actions, don't run them unprompted.

## Backend architecture

Layered, wired manually in `cmd/main.go`: `repo` (pgx SQL) → `service` (business logic) → `handler` (gin) → `server` (gin engine + CORS).

- Each layer declares the narrow interfaces it consumes (e.g. `service/game.GameRepo`, `handler/event.GameService`) in the consumer package; those interfaces are what mockery mocks.
- `models/` holds shared types and **all public sentinel errors**. `models.Errors` maps each error to an HTTP status; `handler/httperr.Error` uses `errors.Is` against that map (default 500). To return a new client-facing status, add an error there and wrap it with `%w`.
- Routes are all in `handler/handler.go`. `/auth/*` is public; `/api/*` runs `authmw.SetUserID` (JWT cookie, optional). `middleware.Require*` and `reqctx` parse/fetch path params (`eventId`, `gameId`, `playerId`) and user ID.
- **Authorization is capability-based per event.** `service/event` builds an auth context (owner, member roles, game participation, event visibility/reach settings) and checks `models.Capability` values via `HasCapability`; handlers call it before acting. `GET /api/event/:id/me` returns the caller's capabilities, which the frontend uses to gate UI (`useEventCapabilities`).
- **Game state is event-sourced.** No in-memory game state survives across requests: `service/game.load` rebuilds a `engine/game.Game` from the starting seating plus the stored round log (`RoundLog` rows replayed through `EnterRoundResult` / `Penalty` / `FalseMahjong`). Mutations apply to the engine, then persist the new round row; when the engine reports `GameEnded`, the service finalizes results. Undo = delete the last round row.
- `engine/game` is the pure MCR rules/scoring engine (seating, winds, rounds, penalties, score) with no DB dependencies; it is unit-tested directly.
- DB migrations are `repo/database/migrations/NNNNNN_*.{up,down}.sql` (golang-migrate), embedded in the binary and applied automatically on startup by `database.NewDB`. Add new ones with the next sequence number and always provide a `down`.
- Logging goes through `logger/` (logrus wrapper); lint forbids std `log` outside `main.go`.
- `test_scripts/user/` holds helper programs for adding/deleting test users.

Lint rules that commonly bite: comments must end with a period (`godot`), no named returns, no `init()`, `math/rand/v2` only, exhaustive switches, `nolint` must name the linter, function length ≤100 lines / 50 statements.

## Frontend architecture

React 19 + TypeScript, Vite, Tailwind v4, HeroUI components, react-router v7, recharts.

- `main.tsx` wraps the app in `AuthProvider` → `EventProvider` → `BrowserRouter`; routes are in `App.tsx`, with layouts in `src/layouts/` and route guards in `components/Auth/` (`RequireAuth`, `RequireEvent`).
- API access lives in `src/providers/*/…Api.ts` as hooks/functions built on `components/Request.tsx` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: toast on non-OK (except 403) and logout + redirect to `/login` on 401.
- `src/models/` mirrors backend JSON shapes; keep them in sync when changing backend models.
- Error tracking via Hawk is enabled in non-development builds when `VITE_HAWK_TOKEN` is set.
