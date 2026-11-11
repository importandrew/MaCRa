package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestValidate_siteURL(t *testing.T) {
	testData := map[string]struct {
		siteURL     string
		expectedErr error
	}{
		"prod":                {siteURL: "https://macra.pro"},
		"trailing slash":      {siteURL: "https://macra.pro/"},
		"local with a port":   {siteURL: "http://localhost:3000"},
		"empty":               {siteURL: "", expectedErr: models.ErrSiteURLInvalid},
		"no scheme":           {siteURL: "macra.pro", expectedErr: models.ErrSiteURLInvalid},
		"other scheme":        {siteURL: "ftp://macra.pro", expectedErr: models.ErrSiteURLInvalid},
		"with a path":         {siteURL: "https://macra.pro/app", expectedErr: models.ErrSiteURLInvalid},
		"with a query string": {siteURL: "https://macra.pro?x=1", expectedErr: models.ErrSiteURLInvalid},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			cfg := Config{JWTKey: "key", MailFrom: "noreply@macra.pro", SiteURL: testCase.siteURL}

			assert.ErrorIs(t, cfg.Validate(), testCase.expectedErr)
		})
	}
}
