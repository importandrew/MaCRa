# MaCRa

***Ma**hjong **C**ompetitive **R**ules **a**ssistant* - application that helps you manage your mahjong MCR games. Everyone can access current game score from their own mobile device, statistics, game logs, host their own event.

## Development requirements
### Backend: 
- PostgreSQL - 15+
- Go - 1.22+
- `/opt/macra/backend_config.json` present: see `backend/config/backend_config.json.example`; every field is described in `backend/config/config.go`.

If `db_url` is missing, server will look for DB URL in `PSQL_MACRA` env variable. Testing depends on `PSQL_MACRA_TEST` env variable set, that contains DB URL to testing DB.

### Frontend:
- NPM - 18+
- `cd frontend && npm install`

## How to start
### Backend
Run backend server:

`$ cd backend && go run ./cmd/main.go`

or

`$ cd backend && make build && ./dist/macra-backend`

### Frontend
Dev: `$ cd frontend && npm run dev`. The dev server forwards `/api` and `/auth` to the backend at `http://127.0.0.1:3001` (set `MACRA_API_PROXY` to change it), so open the app on the dev server's port only; the backend's `site_url` should be that address (e.g. `http://localhost:3000`).

Prod: `$ cd frontend && npm run build`. The app calls the API on its own origin, so the web server must forward `/api` and `/auth` to the backend.