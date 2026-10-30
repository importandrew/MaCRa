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

`.claude/settings.json` enforces some of these rules:
- **Permission rules** deny force pushes, `gh pr merge` and reading `.env*`, and ask before `git push` and `gh pr create`.
- **`.claude/hooks/git-guard.js`** (PreToolUse on Bash) parses each command and blocks force pushes, pushes to `master`/`main` and PR merges, however the command is spelled.
- **`.claude/hooks/format.js`** (PostToolUse on Edit/Write) runs `gofmt -w` on edited `.go` files and the local prettier on edited `frontend/` files. Editing a file therefore also reformats any old unformatted code in it.

The rules above still apply in full. Put personal overrides in `.claude/settings.local.json`, which is gitignored.

## Commands

### Backend (run from `backend/`)
- Run: `go run ./cmd/main.go [-c path/to/backend_config.json]` (default config path `/opt/macra/backend_config.json`; see `config/backend_config.json.example`)
- Regenerate mocks: `mockery` (config in `.mockery.yml`; writes `mocks_test.go` next to each listed interface — add new interfaces there)
- Use `go test -count=1 ./...` to bypass the test cache when DB state matters.

Repo tests (`repo/...`) need a real PostgreSQL in `PSQL_MACRA_TEST`. `database.DB.Test` runs migrations **up before and down after** each test, so the test DB is wiped — never point it at real data. Service and handler tests use mockery mocks and need no DB.

Config: `jwt_key` and `mail_from` are required (`Validate` fails otherwise). If `db_url` is empty, `PSQL_MACRA` env is used. Mail is sent via SMTP on `127.0.0.1:25` (no auth).

### Frontend (run from `frontend/`)
- Install with `npm ci`, which installs from the lockfile without modifying it.
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
- DB migrations are `repo/database/migrations/NNNNNN_*.{up,down}.sql` (golang-migrate), embedded in the binary and applied automatically on startup by `database.NewDB`. Add new ones with the next sequence number and always provide a `down`.

## Frontend architecture

- API access lives in `src/providers/*/…Api.ts` as hooks/functions built on `components/Request.tsx` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: toast on non-OK (except 403) and logout + redirect to `/login` on 401.
- `src/models/` mirrors backend JSON shapes; keep them in sync when changing backend models.
