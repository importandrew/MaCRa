package player

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieAttrs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		configured     string
		requestHost    string
		forwardedProto string
		tls            bool
		wantDomain     string
		wantSecure     bool
	}{
		"http lan mismatched localhost config": {
			configured:  "localhost",
			requestHost: "mj.home:3001",
			wantDomain:  "",
			wantSecure:  false,
		},
		"http lan matching host with http prefix": {
			configured:  "http://mj.home",
			requestHost: "mj.home:3001",
			wantDomain:  "mj.home",
			wantSecure:  false,
		},
		"https matching domain": {
			configured:  "https://macra.pro",
			requestHost: "macra.pro",
			tls:         true,
			wantDomain:  "macra.pro",
			wantSecure:  true,
		},
		"empty config on https uses host-only": {
			configured:     "",
			requestHost:    "dev.macra.pro",
			forwardedProto: "https",
			wantDomain:     "",
			wantSecure:     true,
		},
		"empty config on http uses host-only": {
			configured:  "",
			requestHost: "mj.home:3001",
			wantDomain:  "",
			wantSecure:  false,
		},
		"bare domain on https request": {
			configured:  "example.com",
			requestHost: "example.com",
			tls:         true,
			wantDomain:  "example.com",
			wantSecure:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			domain, secure := cookieAttrs(tc.configured, tc.requestHost, tc.forwardedProto, tc.tls)
			assert.Equal(t, tc.wantDomain, domain)
			assert.Equal(t, tc.wantSecure, secure)
		})
	}
}
