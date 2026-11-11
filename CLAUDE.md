# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

MaCRa (Mahjong Competitive Rules assistant) tracks MCR mahjong games: live scores viewable from any device, game logs, ratings, and events that users host. Monorepo with a Go backend (`backend/`, module `github.com/importandrew/MaCRa`) and a React/Vite frontend (`frontend/`). Each part has its own notes in `backend/CLAUDE.md` and `frontend/CLAUDE.md`.

## Working rules

- Start every task on a new branch cut from a freshly fetched master: `git fetch origin && git checkout -b feature/<short-name> origin/master`. Don't let the branch track `origin/master`; push with `git push -u origin feature/<short-name>`.
- Never push to `master` and never merge PRs. Only the repository owner merges.
- No force pushes (`--force`, `-f`, `--force-with-lease`, `+refspec`). Don't commit to, push to or rewrite branches you didn't create.
- Before every commit, run the backend tests and `golangci-lint run`, plus `npm run lint` and `npm run build` if `frontend/` changed. Master is not fully green yet: don't add new failures or lint issues compared with master. When unsure, run the same checks on master.
- Don't read or commit `.env*` files, keys, passwords or tokens. That includes `backend_config.json` with real values (`db_url`, `jwt_key`).
- Ask the owner before changing DB migrations (`backend/repo/database/migrations/`), CI (`.github/`), or dependencies (`go.mod`/`go.sum`, `package.json`/`package-lock.json`).
- Commit messages follow Conventional Commits: `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, `docs:`, … (as in the existing history).
- Write everything that lands in the repository or on GitHub in English: code, comments, commit messages, branch names, issues, PR titles and descriptions. This applies even when the conversation is in another language.
- When the code you work with does the same thing in different ways (e.g. two query or error-handling styles), or uses a pattern, coupling or package worth improving, mention it in a line or two at the end of your reply. Don't fix it without being asked; keep the task's diff focused. In new code, follow the prevailing pattern and say which one you followed.

`.claude/settings.json` enforces some of these rules:
- **Permission rules** deny force pushes, `gh pr merge` and reading `.env*`, and ask before `git push` and `gh pr create`.
- **`.claude/hooks/git-guard.js`** (PreToolUse on Bash) parses each command and blocks force pushes, pushes to `master`/`main` and PR merges, however the command is spelled.
- **`.claude/hooks/format.js`** (PostToolUse on Edit/Write) runs `gofmt -w` on edited `.go` files and the local prettier on edited `frontend/` files. Editing a file therefore also reformats any old unformatted code in it.

The rules above still apply in full. Put personal overrides in `.claude/settings.local.json`, which is gitignored.

## Product behavior

`docs/user-scenarios.md` is the agreed, desired behavior of the app, area by area. Read the relevant area before changing what users see or may do; if the code and the file disagree, raise it with the owner instead of picking one silently. Principles that apply everywhere:

- Mobile first (most users are on phones); supported from 320 px wide. Game screens avoid the keyboard: steppers and toggles, players picked by seat.
- At the table, scores change only when the player taps Refresh or enters a round; elsewhere auto-refresh is an opt-in switch.
- The event owner always has staff rights; a player always sees their own games while they can open the event; without access to an event, nothing in it is visible.
- An event's game rules freeze after its first game, and each game is replayed with the rules it started with; permissions can change at any time.
- Players' privacy: emails are never shown to other users; non-members get staff names and a player count, not the member list (ratings and games show names only as their reach allows).

## Backend ↔ frontend contracts

- The frontend and the API share one origin: the frontend calls relative `/api/...` and `/auth/...`, which nginx forwards to the backend on prod and vite's dev server does locally (`frontend/vite.config.ts`, target `MACRA_API_PROXY`, default `http://127.0.0.1:3001`). So there is no CORS, and the auth cookie is host-only (each site its own session). The backend's `site_url` in `backend_config.json` is the site's own address; links in auth emails point to it.
- `frontend/src/models/` mirrors the backend JSON shapes. Change both sides together.
- The UI is gated by per-event capabilities from `GET /api/event/:id/me` (`useEventCapabilities`). Adding a capability on the backend means handling it on the frontend.
- `deploy*` targets in the Makefiles copy builds to the production/dev server over SSH (`MACRA_SSHKEY`/`MACRA_HOST`) and restart systemd services. Never run them unprompted.
