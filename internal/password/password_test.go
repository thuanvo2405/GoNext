package password

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordAndVerify(t *testing.T) {
	plaintext := "TestPassword123!"
	hash, err := HashPassword(plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	assert.NotEqual(t, plaintext, hash)
	assert.True(t, Verify(plaintext, hash))
	assert.False(t, Verify("wrong-password", hash))
}

func TestVerifyInvalidHash(t *testing.T) {
	for _, hash := range []string{"", "not-a-bcrypt-hash"} {
		t.Run(hash, func(t *testing.T) {
			assert.False(t, Verify("TestPassword123!", hash))
		})
	}
}

func TestHashPasswordByteLimit(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		tooLong   bool
	}{
		{"72 ASCII bytes accepted", strings.Repeat("a", 72), false},
		{"72 Unicode bytes accepted", strings.Repeat("界", 24), false},
		{"73 ASCII bytes rejected", strings.Repeat("a", 73), true},
		{"75 Unicode bytes rejected", strings.Repeat("界", 25), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.plaintext)
			if tt.tooLong {
				require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
				assert.Empty(t, hash)
				return
			}

			require.NoError(t, err)
			assert.True(t, Verify(tt.plaintext, hash))
		})
	}
}
