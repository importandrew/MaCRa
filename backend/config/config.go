package config

import (
	"encoding/json"
	"os"

	"github.com/importandrew/MaCRa/models"
)

type Config struct {
	Addr         string   `json:"addr"`          // Backend API listen address.
	CertFile     string   `json:"cert_file"`     // Path to cert file for https. If not empty, server will run in TLS mode.
	KeyFile      string   `json:"key_file"`      // Path to key file for https.
	DBURL        string   `json:"db_url"`        // Database URL.
	CookieDomain string   `json:"cookie_domain"` // Exact hostname for this instance (optional http:// or https:// prefix). http:// forces Secure=false; https:// forces Secure=true. Empty → host-only cookie. Host mismatch → host-only. Do not use a parent domain shared by sibling envs.
	AllowOrigins []string `json:"allow_origins"` // List of origins allowed to API.
	DBLog        bool     `json:"db_log"`        // Whether to log DB actions for debugging.
	JWTKey       string   `json:"jwt_key"`       // Key to generate/parse JWT tokens.
	MailFrom     string   `json:"mail_from"`     // Envelope and From address for outbound auth mail (must match Postfix/DKIM).
}

func Read(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.JWTKey == "" {
		return models.ErrJWTKeyEmpty
	}

	if c.MailFrom == "" {
		return models.ErrMailFromEmpty
	}

	return nil
}
