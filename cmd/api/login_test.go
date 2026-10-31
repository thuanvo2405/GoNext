package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These cases stop before database access; no test database is required.
func TestLoginInputValidation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"empty body", "", http.StatusBadRequest},
		{"malformed JSON", `{"email":`, http.StatusBadRequest},
		{"missing email", `{"password":"password123"}`, http.StatusBadRequest},
		{"missing password", `{"email":"rider@example.com"}`, http.StatusBadRequest},
		{"empty fields", `{"email":"","password":""}`, http.StatusBadRequest},
		{"invalid email", `{"email":"not-an-email","password":"password123"}`, http.StatusBadRequest},
		{"email has wrong type", `{"email":123,"password":"password123"}`, http.StatusBadRequest},
		{"password has wrong type", `{"email":"rider@example.com","password":123}`, http.StatusBadRequest},
		{"null fields", `{"email":null,"password":null}`, http.StatusBadRequest},
		{
			"password exceeds 72 bytes",
			`{"email":"rider@example.com","password":"` + strings.Repeat("a", 73) + `"}`,
			http.StatusUnauthorized,
		},
		{
			"Unicode password exceeds 72 bytes",
			`{"email":"rider@example.com","password":"` + strings.Repeat("界", 25) + `"}`,
			http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			app.setupRouter().ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusUnauthorized {
				assert.JSONEq(t, `{"error":"invalid email or password"}`, w.Body.String())
			} else {
				assert.JSONEq(t, `{"error":"invalid login input"}`, w.Body.String())
			}
		})
	}
}
