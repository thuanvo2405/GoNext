package token

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test-only key, not an application default or production secret.
const testSecret = "test-only-signing-secret-at-least-32-bytes"

func TestManagerSignAndVerify(t *testing.T) {
	manager, err := NewManager(testSecret, 30*time.Minute)
	require.NoError(t, err)
	userID := "00000000-0000-4000-8000-000000000001"
	before := time.Now().Unix()
	signed, err := manager.Sign(userID)
	require.NoError(t, err)
	after := time.Now().Unix()

	subject, err := manager.Verify(signed)
	require.NoError(t, err)
	assert.Equal(t, userID, subject)

	parsed, err := jwt.Parse(signed, func(_ *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	require.NoError(t, err)
	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims["userId"])
	expiration, err := claims.GetExpirationTime()
	require.NoError(t, err)
	require.NotNil(t, expiration)
	assert.GreaterOrEqual(t, expiration.Unix(), before+1800)
	assert.LessOrEqual(t, expiration.Unix(), after+1800)
}

func TestManagerVerifyInvalidTokens(t *testing.T) {
	manager, err := NewManager(testSecret, 30*time.Minute)
	require.NoError(t, err)
	now := time.Now()
	validClaims := func() jwt.MapClaims {
		return jwt.MapClaims{
			"userId": "user-id",
			"exp":    now.Add(30 * time.Minute).Unix(),
		}
	}
	signClaims := func(method jwt.SigningMethod, claims jwt.MapClaims, key any) string {
		t.Helper()
		signed, err := jwt.NewWithClaims(method, claims).SignedString(key)
		require.NoError(t, err)
		return signed
	}

	expired := validClaims()
	expired["exp"] = now.Add(-time.Minute).Unix()
	invalidExpiry := validClaims()
	invalidExpiry["exp"] = "not-a-timestamp"
	missingUserID := validClaims()
	delete(missingUserID, "userId")
	emptyUserID := validClaims()
	emptyUserID["userId"] = ""
	blankUserID := validClaims()
	blankUserID["userId"] = "   "
	numericUserID := validClaims()
	numericUserID["userId"] = 123
	nullUserID := validClaims()
	nullUserID["userId"] = nil
	futureNotBefore := validClaims()
	futureNotBefore["nbf"] = now.Add(time.Hour).Unix()

	valid := signClaims(jwt.SigningMethodHS256, validClaims(), []byte(testSecret))
	parts := strings.Split(valid, ".")
	parts[1] = base64.RawURLEncoding.EncodeToString([]byte(`{"userId":"another-user","exp":4102444800}`))
	tampered := strings.Join(parts, ".")

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"malformed token", "not-a-jwt"},
		{"wrong key", signClaims(jwt.SigningMethodHS256, validClaims(), []byte("different-test-secret"))},
		{"expired", signClaims(jwt.SigningMethodHS256, expired, []byte(testSecret))},
		{"invalid expiration type", signClaims(jwt.SigningMethodHS256, invalidExpiry, []byte(testSecret))},
		{"missing user ID", signClaims(jwt.SigningMethodHS256, missingUserID, []byte(testSecret))},
		{"empty user ID", signClaims(jwt.SigningMethodHS256, emptyUserID, []byte(testSecret))},
		{"blank user ID", signClaims(jwt.SigningMethodHS256, blankUserID, []byte(testSecret))},
		{"numeric user ID", signClaims(jwt.SigningMethodHS256, numericUserID, []byte(testSecret))},
		{"null user ID", signClaims(jwt.SigningMethodHS256, nullUserID, []byte(testSecret))},
		{"future not-before", signClaims(jwt.SigningMethodHS256, futureNotBefore, []byte(testSecret))},
		{"unsigned token", signClaims(jwt.SigningMethodNone, validClaims(), jwt.UnsafeAllowNoneSignatureType)},
		{"tampered payload", tampered},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject, err := manager.Verify(tt.token)
			assert.Error(t, err)
			assert.Empty(t, subject)
		})
	}
}

// Document the current Verify policy: expiration is optional and all HMAC
// algorithms are accepted. These tests do not endorse this as a stricter policy.
func TestManagerVerifyCurrentPolicy(t *testing.T) {
	manager, err := NewManager(testSecret, 30*time.Minute)
	require.NoError(t, err)
	tests := []struct {
		name   string
		method jwt.SigningMethod
		claims jwt.MapClaims
	}{
		{"missing expiration accepted", jwt.SigningMethodHS256, jwt.MapClaims{"userId": "user-id"}},
		{"HS384 accepted", jwt.SigningMethodHS384, jwt.MapClaims{"userId": "user-id", "exp": time.Now().Add(30 * time.Minute).Unix()}},
		{"HS512 accepted", jwt.SigningMethodHS512, jwt.MapClaims{"userId": "user-id", "exp": time.Now().Add(30 * time.Minute).Unix()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := jwt.NewWithClaims(tt.method, tt.claims).SignedString([]byte(testSecret))
			require.NoError(t, err)
			subject, err := manager.Verify(signed)
			require.NoError(t, err)
			assert.Equal(t, "user-id", subject)
		})
	}
}

func TestNewManagerInvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		ttl    time.Duration
	}{
		{"empty secret", "", 30 * time.Minute},
		{"short secret", strings.Repeat("a", 31), 30 * time.Minute},
		{"blank secret", strings.Repeat(" ", 32), 30 * time.Minute},
		{"zero TTL", testSecret, 0},
		{"negative TTL", testSecret, -time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.secret, tt.ttl)
			require.Error(t, err)
			assert.Nil(t, manager)
		})
	}
}

func TestNewManagerMinimumSecretLength(t *testing.T) {
	manager, err := NewManager(strings.Repeat("a", 32), 30*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, manager)
}
