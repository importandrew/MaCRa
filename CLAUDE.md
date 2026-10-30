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

`.claude/settings.json` enforces some of these rules:
- **Permission rules** deny force pushes, `gh pr merge` and reading `.env*`, and ask before `git push` and `gh pr create`.
- **`.claude/hooks/git-guard.js`** (PreToolUse on Bash) parses each command and blocks force pushes, pushes to `master`/`main` and PR merges, however the command is spelled.
- **`.claude/hooks/format.js`** (PostToolUse on Edit/Write) runs `gofmt -w` on edited `.go` files and the local prettier on edited `frontend/` files. Editing a file therefore also reformats any old unformatted code in it.

The rules above still apply in full. Put personal overrides in `.claude/settings.local.json`, which is gitignored.

## Backend ↔ frontend contracts

- The frontend reaches the API at `VITE_API_ADDR` from `frontend/.env.<mode>`. For local dev, the backend's `allow_origins` and `cookie_domain` in `backend_config.json` must match the frontend origin, or auth cookies and CORS break.
- `frontend/src/models/` mirrors the backend JSON shapes. Change both sides together.
- The UI is gated by per-event capabilities from `GET /api/event/:id/me` (`useEventCapabilities`). Adding a capability on the backend means handling it on the frontend.
- `deploy*` targets in the Makefiles copy builds to the production/dev server over SSH (`MACRA_SSHKEY`/`MACRA_HOST`) and restart systemd services. Never run them unprompted.
