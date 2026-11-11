package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_trimsTrailingSlash(t *testing.T) {
	assert.Equal(t, "https://macra.pro", New("noreply@macra.pro", "https://macra.pro/").siteURL)
	assert.Equal(t, "http://localhost:3000", New("noreply@macra.pro", "http://localhost:3000").siteURL)
}
