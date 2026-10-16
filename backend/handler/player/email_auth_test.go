package player

import (
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeOrigin(t *testing.T) {
	tests := map[string]struct {
		origin  string
		extra   []string
		want    string
		wantErr error
	}{
		"prod": {
			origin: "https://macra.pro",
			want:   "https://macra.pro",
		},
		"prod with path stripped": {
			origin: "https://macra.pro/login",
			want:   "https://macra.pro",
		},
		"dev": {
			origin: "https://dev.macra.pro",
			want:   "https://dev.macra.pro",
		},
		"allow list local": {
			origin: "http://localhost:5173",
			extra:  []string{"http://localhost:5173"},
			want:   "http://localhost:5173",
		},
		"rejected": {
			origin:  "https://evil.example",
			wantErr: models.ErrInvalidOrigin,
		},
		"empty": {
			origin:  "",
			wantErr: models.ErrInvalidOrigin,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := normalizeOrigin(tc.origin, tc.extra)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
