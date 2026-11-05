package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testKey = []byte("test-key")

func signed(t *testing.T, claims *UserClaims, key []byte) string {
	t.Helper()

	tkn, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	require.NoError(t, err)

	return tkn
}

func expiringIn(d time.Duration) *UserClaims {
	return &UserClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
		UserID:    7,
	}
}

func TestIdentifyUser(t *testing.T) {
	t.Parallel()

	fresh, err := GenerateJWT(7, testKey)
	require.NoError(t, err)

	tests := map[string]struct {
		token       string
		wantErr     bool
		wantRefresh bool
	}{
		"fresh token": {
			token: fresh,
		},
		"token older than a day needs refresh": {
			token:       signed(t, expiringIn(Lifetime-RefreshAfter-time.Minute), testKey),
			wantRefresh: true,
		},
		"legacy 90-day token needs refresh": {
			token:       signed(t, expiringIn(90*24*time.Hour), testKey),
			wantRefresh: true,
		},
		"expired token": {
			token:   signed(t, expiringIn(-time.Minute), testKey),
			wantErr: true,
		},
		"token without expiration": {
			token:   signed(t, &UserClaims{UserID: 7}, testKey),
			wantErr: true,
		},
		"token signed with another key": {
			token:   signed(t, expiringIn(Lifetime), []byte("other-key")),
			wantErr: true,
		},
		"garbage": {
			token:   "garbage",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			identity, err := IdentifyUser(tc.token, testKey)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, 7, identity.UserID)
			assert.Equal(t, tc.wantRefresh, identity.NeedsRefresh)
		})
	}
}
