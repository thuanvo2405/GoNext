package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thuanvo2405/GoNext/internal/data"
)

func TestGetMe(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		setUser    bool
		wantStatus int
		wantBody   string
	}{
		{
			name: "current user", setUser: true,
			value:      &data.User{ID: "user-id", Email: "rider@example.com", Role: "rider", Password: "hidden-hash"},
			wantStatus: http.StatusOK,
			wantBody:   `{"user":{"id":"user-id","email":"rider@example.com","role":"rider"}}`,
		},
		{"missing context", nil, false, http.StatusInternalServerError, `{"error":"invalid user context"}`},
		{"wrong type", "user-id", true, http.StatusInternalServerError, `{"error":"invalid user context"}`},
		{"nil user", (*data.User)(nil), true, http.StatusInternalServerError, `{"error":"invalid user context"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.setUser {
				c.Set("user", tt.value)
			}
			app.GetMe(c)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
