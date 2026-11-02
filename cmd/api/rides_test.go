package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thuanvo2405/GoNext/internal/data"
)

func TestEstimateFareValidation(t *testing.T) {
	for _, tt := range []struct {
		name, body string
	}{
		{"empty", ""},
		{"malformed", `{"pickup":`},
		{"missing points", `{}`},
		{"missing destination", `{"pickup":{"lat":0,"lng":0}}`},
		{"null pickup", `{"pickup":null,"destination":{"lat":1,"lng":1}}`},
		{"missing latitude", `{"pickup":{"lng":0},"destination":{"lat":1,"lng":1}}`},
		{"null longitude", `{"pickup":{"lat":0,"lng":null},"destination":{"lat":1,"lng":1}}`},
		{"wrong type", `{"pickup":{"lat":"zero","lng":0},"destination":{"lat":1,"lng":1}}`},
		{"latitude too high", `{"pickup":{"lat":91,"lng":0},"destination":{"lat":1,"lng":1}}`},
		{"latitude too low", `{"pickup":{"lat":0,"lng":0},"destination":{"lat":-91,"lng":1}}`},
		{"longitude too high", `{"pickup":{"lat":0,"lng":181},"destination":{"lat":1,"lng":1}}`},
		{"longitude too low", `{"pickup":{"lat":0,"lng":0},"destination":{"lat":1,"lng":-181}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			router := gin.New()
			router.POST("/rides/estimate", app.estimateFareHandler)
			req := httptest.NewRequest(http.MethodPost, "/rides/estimate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, `{"error":"invalid estimate input"}`, w.Body.String())
		})
	}
}

func TestEstimateFareResponse(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		status     int
	}{
		{"same points", `{"pickup":{"lat":0,"lng":0},"destination":{"lat":0,"lng":0}}`, http.StatusBadRequest},
		{"same nonzero points", `{"pickup":{"lat":10,"lng":20},"destination":{"lat":10,"lng":20}}`, http.StatusBadRequest},
		{"same latitude", `{"pickup":{"lat":0,"lng":0},"destination":{"lat":0,"lng":1}}`, http.StatusOK},
		{"same longitude", `{"pickup":{"lat":0,"lng":0},"destination":{"lat":1,"lng":0}}`, http.StatusOK},
		{"boundary values", `{"pickup":{"lat":90,"lng":180},"destination":{"lat":-90,"lng":-180}}`, http.StatusOK},
	} {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			router := gin.New()
			router.POST("/rides/estimate", app.estimateFareHandler)
			req := httptest.NewRequest(http.MethodPost, "/rides/estimate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			require.Equal(t, tt.status, w.Code)
			if tt.status == http.StatusBadRequest {
				assert.JSONEq(t, `{"error":"pickup and destination must be different"}`, w.Body.String())
				return
			}
			var response struct {
				DistanceKm     float64 `json:"distance_km"`
				DistanceSource string  `json:"distance_source"`
				Estimates      []struct {
					VehicleType string `json:"vehicle_type"`
					FareVND     int    `json:"fare_vnd"`
				} `json:"estimates"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
			assert.Greater(t, response.DistanceKm, 0.0)
			assert.Equal(t, "straight_line", response.DistanceSource)
			require.Len(t, response.Estimates, 3)
			for i, vehicleType := range []string{"motorcycle", "car", "electric_car"} {
				assert.Equal(t, vehicleType, response.Estimates[i].VehicleType)
				assert.Greater(t, response.Estimates[i].FareVND, 0)
				assert.Zero(t, response.Estimates[i].FareVND%1000)
			}
			if tt.name == "same latitude" {
				assert.InDelta(t, 111.194927, response.DistanceKm, 0.001)
				assert.Equal(t, 556_000, response.Estimates[0].FareVND)
				assert.Equal(t, 1_335_000, response.Estimates[1].FareVND)
				assert.Equal(t, 1_557_000, response.Estimates[2].FareVND)
			}
		})
	}
}

// Inject the user context to test authorization without database access.
func TestEstimateFareRole(t *testing.T) {
	for _, role := range []string{"rider", "driver"} {
		t.Run(role, func(t *testing.T) {
			app := &application{}
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user", &data.User{ID: "user-id", Role: role})
				c.Next()
			})
			router.POST("/rides/estimate", app.RequireRole("rider"), app.estimateFareHandler)
			req := httptest.NewRequest(http.MethodPost, "/rides/estimate", strings.NewReader(`{"pickup":{"lat":0,"lng":0},"destination":{"lat":0,"lng":1}}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if role == "driver" {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}
