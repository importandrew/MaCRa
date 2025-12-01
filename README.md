# MaCRa

***Ma**hjong **C**ompetitive **R**ules **a**ssistant* - application that helps you manage your mahjong MCR games. Everyone can access current game score from their own mobile device, statistics, game logs, host their own event.

## Development requirements
### Backend: 
- PostgreSQL - 15+
- Go - 1.22+
- `/opt/macra/backend_config.json` present. Config structure:
```
type Config struct {
	Addr         string   `json:"addr"`          // Backend API listen address.
	CertFile     string   `json:"cert_file"`     // Path to cert file for https. If not empty, server will run in TLS mode.
	KeyFile      string   `json:"key_file"`      // Path to key file for https.
	DBURL        string   `json:"db_url"`        // Database URL.
	CookieDomain string   `json:"cookie_domain"` // Domain for cookie. Protocol defines secure or not. Default: macra.pro.
	AllowOrigins []string `json:"allow_origins"` // List of origins allowed to API.
	DBLog        bool     `json:"db_log"`        // Whether to log DB actions for debugging.
	JWTKey       string   `json:"jwt_key"`       // Key to generate/parse JWT tokens.
}
```
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
Dev: `$ cd frontend && npm run dev` (*backend_config.json* should have relevant settings for `cookie_domain` and `allow_origins`)

Prod: `$ cd frontend && npm run build`

Backend API address should match `VITE_API_ADDR` parameter in corresponding mode `.env` file.