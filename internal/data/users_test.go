package data

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserJSONDoesNotExposePassword(t *testing.T) {
	user := User{
		ID:       "00000000-0000-4000-8000-000000000001",
		Email:    "rider@example.com",
		Password: "stored-password-hash",
		Role:     "rider",
	}

	encoded, err := json.Marshal(user)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id":"00000000-0000-4000-8000-000000000001",
		"email":"rider@example.com",
		"role":"rider"
	}`, string(encoded))
}
