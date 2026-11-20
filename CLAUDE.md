# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

MaCRa (Mahjong Competitive Rules assistant) tracks MCR mahjong games: live scores viewable from any device, game logs, ratings, and events that users host. A game ends with TP (table points by place: 4/2/1/0, shared on ties) and GP (the score) per player; ratings sort by TP, then GP. Monorepo with a Go backend (`backend/`, module `github.com/importandrew/MaCRa`) and a React/Vite frontend (`frontend/`). Each part has its own notes in `backend/CLAUDE.md` and `frontend/CLAUDE.md`: read the one for the part you change before starting. If `AGENTS.local.md` exists in the repository root, read it too: the owner's private, git-ignored working notes for this machine; never commit or quote it.

## Working rules

- Start every task on a new branch cut from a freshly fetched master: `git fetch origin && git checkout -b feature/<short-name> origin/master --no-track` (prefixes in use: `feature/`, `fix/`, `chore/`, `refactor/`, `docs/`). Push with `git push -u origin feature/<short-name>`.
- Tournament work goes to the integration branch `feature/tournament`, not to master: cut the task branch from `origin/feature/tournament` and open the PR with `--base feature/tournament`. Never push to `feature/tournament` itself, and ask before merging master into it.
- Commit, push or open a PR only when asked.
- Never push to `master` and never merge PRs. Only the repository owner merges.
- No force pushes (`--force`, `-f`, `--force-with-lease`, `+refspec`). Don't commit to, push to or rewrite branches you didn't create.
- Before every commit, run the backend tests and `golangci-lint run`, plus `npm run lint` and `npm run build` if `frontend/` changed: `(cd backend && go test -count=1 -p 1 ./... && golangci-lint run)` and `(cd frontend && npm run lint && npm run build)`. Master is not fully green yet: don't add new failures or lint issues compared with master. When unsure, run the same checks on master.
- Don't read or commit `.env*` files, keys, passwords or tokens. That includes `backend_config.json` with real values (`db_url`, `jwt_key`).
- Ask the owner before changing DB migrations (`backend/repo/database/migrations/`), CI (`.github/`), or dependencies (`go.mod`/`go.sum`, `package.json`/`package-lock.json`).
- Commit messages follow Conventional Commits: `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, `docs:`, … (as in the existing history).
- Write everything that lands in the repository or on GitHub in English: code, comments, commit messages, branch names, issues, PR titles and descriptions. This applies even when the conversation is in another language.
- Git-ignored local notes (`*.local.md` in the root) may exist in the checkout. Never refer to them or their item numbers ("KI n", "TODO n", a plan item N) in code, comments, commits, PRs or issues: describe the thing itself. Citing `docs/user-scenarios.md` ("area 4") is fine. Before a commit or PR, grep the diff and the message for `KI [0-9]|TODO [0-9]|item [0-9]|of the pla[n]`.
- When the code you work with does the same thing in different ways (e.g. two query or error-handling styles), or uses a pattern, coupling or package worth improving, mention it in a line or two at the end of your reply. Don't fix it without being asked; keep the task's diff focused. In new code, follow the prevailing pattern and say which one you followed.
- When a fix runs into the existing structure, don't bend around it with helper layers (extra interfaces, "own vs external state" props, reload keys, retry guards): stop and propose one or two clean redesigns with trade-offs.
- Fixing a reported UI regression means restoring the previous behavior exactly; offer further layout or UX changes as a question instead of shipping them.
- When you describe how something works now, also propose how it should work, and why; mark what is purely a product decision.
- The owner deploys only when no game is live, so new request fields for game screens can be required at once, with no fallback for an old client. Mention in the PR if a tab left open across a deploy could break.
- PR descriptions: what changes for users, how it was verified, and anything the owner must do (migration, server, config).

`.claude/settings.json` enforces some of these rules:
- **Permission rules** deny force pushes, `gh pr merge` and reading `.env*`, and ask before `git push` and `gh pr create`.
- **`.claude/hooks/git-guard.js`** (PreToolUse on Bash) parses each command and blocks force pushes, pushes to `master`/`main` and PR merges, however the command is spelled.
- **`.claude/hooks/format.js`** (PostToolUse on Edit/Write) runs `gofmt -w` on edited `.go` files and the local prettier on edited `frontend/` files. Editing a file therefore also reformats any old unformatted code in it.

The rules above still apply in full. Put personal overrides in `.claude/settings.local.json`, which is gitignored. Other coding agents (e.g. opencode) don't run these hooks and permission rules: follow the rules yourself, and format every file you edit (`gofmt -w` for `.go`, `frontend/node_modules/.bin/prettier --write` for files under `frontend/`).

## Product behavior

`docs/user-scenarios.md` is the agreed, desired behavior of the app, area by area. Read the relevant area before changing what users see or may do; if the code and the file disagree, raise it with the owner instead of picking one silently. Principles that apply everywhere:

- Mobile first (most users are on phones); supported from 320 px wide. Game screens avoid the keyboard: steppers and toggles, players picked by seat; touch targets at least 44 px.
- At the table, scores change only when the player taps Refresh or enters a round; elsewhere auto-refresh is an opt-in switch.
- The event owner sees everything in their event and manages it; acting as a referee (undo, abort) takes the referee role, also for the owner. A player always sees their own games while they can open the event; without access to an event, nothing in it is visible.
- An event's game rules freeze after its first game, and each game is replayed with the rules it started with; permissions can change at any time.
- Players' privacy: emails are never shown to other users; non-members get staff names and a player count, not the member list (ratings and games show names only as their reach allows).

## Backend ↔ frontend contracts

- The frontend and the API share one origin: the frontend calls relative `/api/...` and `/auth/...`, which nginx forwards to the backend on prod and vite's dev server does locally (`frontend/vite.config.ts`, target `MACRA_API_PROXY`, default `http://127.0.0.1:3001`). So there is no CORS, and the auth cookie is host-only (each site its own session). The backend's `site_url` in `backend_config.json` is the site's own address; links in auth emails point to it.
- `frontend/src/models/` mirrors the backend JSON shapes. Change both sides together.
- API errors are stable codes (`{"error": "<code>"}` plus data); the frontend words each code (`apiErrors` in `frontend/src/i18n/locales/`). A new client-facing error needs a backend sentinel and code, and its text in every locale.
- The UI is gated by capabilities: per event from `GET /api/event/:id/me` (`useEventCapabilities`), per game from `GET /api/event/:id/game/:gameId/me` (`useGameCapabilities`: enter rounds, undo, abort, finish). Adding a capability on the backend means handling it on the frontend.
- `deploy*` targets in the Makefiles copy builds to the production/dev server over SSH (`MACRA_SSHKEY`/`MACRA_HOST`) and restart systemd services. Never run them unprompted.
