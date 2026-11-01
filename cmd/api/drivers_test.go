package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thuanvo2405/GoNext/internal/data"
)

// Invalid input is rejected before database access.
func TestDriverAvailabilityValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty body", ""},
		{"malformed JSON", `{"is_available":`},
		{"missing field", `{}`},
		{"null field", `{"is_available":null}`},
		{"string field", `{"is_available":"true"}`},
		{"numeric field", `{"is_available":1}`},
		{"array field", `{"is_available":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user", &data.User{ID: "driver-id", Role: "driver"})
				c.Next()
			})
			router.PATCH("/drivers/availability", app.RequireRole("driver"), app.updateDriverAvailability)
			req := httptest.NewRequest(http.MethodPatch, "/drivers/availability", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, `{"error":"invalid availability input"}`, w.Body.String())
		})
	}
}

func TestRiderCannotUpdateDriverAvailability(t *testing.T) {
	app := &application{}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user", &data.User{ID: "rider-id", Role: "rider"})
		c.Next()
	})
	router.PATCH("/drivers/availability", app.RequireRole("driver"), app.updateDriverAvailability)
	req := httptest.NewRequest(http.MethodPatch, "/drivers/availability", strings.NewReader(`{"is_available":true}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"error":"you do not have permission"}`, w.Body.String())
}
