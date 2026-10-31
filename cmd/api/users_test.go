package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Invalid inputs must be rejected before the handler accesses the database.
func TestRegisterRiderValidation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "malformed JSON",
			body:       `{"email":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       `{"password":"password123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       `{"email":"rider@example.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty fields",
			body:       `{"email":"","password":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email",
			body:       `{"email":"not-an-email","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password below minimum length",
			body:       `{"email":"rider@example.com","password":"1234567"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "email has wrong type",
			body:       `{"email":123,"password":"password123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password has wrong type",
			body:       `{"email":"rider@example.com","password":12345678}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "null fields",
			body:       `{"email":null,"password":null}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password exceeds 72 bytes",
			body:       `{"email":"rider@example.com","password":"` + strings.Repeat("a", 73) + `"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "your password too long",
		},
		{
			name:       "Unicode password exceeds 72 bytes but not 72 characters",
			body:       `{"email":"rider@example.com","password":"` + strings.Repeat("界", 25) + `"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "your password too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			router := app.setupRouter()
			req := httptest.NewRequest(http.MethodPost, "/auth/register/rider", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
			if tt.wantError == "your password too long" {
				assert.JSONEq(t, `{"error":"your password too long"}`, w.Body.String())
			} else {
				assert.JSONEq(t, `{"error":"invalid registration input"}`, w.Body.String())
			}
		})
	}
}
