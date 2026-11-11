package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/importandrew/MaCRa/models"
)

type Config struct {
	Addr     string `json:"addr"`      // Backend API listen address.
	CertFile string `json:"cert_file"` // Path to cert file for https. If not empty, server will run in TLS mode.
	KeyFile  string `json:"key_file"`  // Path to key file for https.
	DBURL    string `json:"db_url"`    // Database URL.
	SiteURL  string `json:"site_url"`  // The site's address, e.g. https://macra.pro: the frontend and the API (under /api and /auth) live there. Links in auth emails point to it.
	DBLog    bool   `json:"db_log"`    // Whether to log DB actions for debugging.
	JWTKey   string `json:"jwt_key"`   // Key to generate/parse JWT tokens.
	MailFrom string `json:"mail_from"` // Envelope and From address for outbound auth mail (must match Postfix/DKIM).
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

	if err := validateSiteURL(c.SiteURL); err != nil {
		return err
	}

	return nil
}

// validateSiteURL reports whether siteURL is a bare http(s) address with no
// path, so links like siteURL+"/reset-password?…" come out right.
func validateSiteURL(siteURL string) error {
	u, err := url.Parse(siteURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" ||
		(u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("%w: %q", models.ErrSiteURLInvalid, siteURL)
	}

	return nil
}
