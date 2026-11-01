package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thuanvo2405/GoNext/internal/data"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		user       any
		setUser    bool
		wantStatus int
		wantBody   string
	}{
		{"driver allowed", "driver", &data.User{Role: "driver"}, true, http.StatusNoContent, ""},
		{"rider allowed", "rider", &data.User{Role: "rider"}, true, http.StatusNoContent, ""},
		{"rider denied driver route", "driver", &data.User{Role: "rider"}, true, http.StatusForbidden, `{"error":"you do not have permission"}`},
		{"driver denied rider route", "rider", &data.User{Role: "driver"}, true, http.StatusForbidden, `{"error":"you do not have permission"}`},
		{"missing user", "driver", nil, false, http.StatusInternalServerError, `{"error":"invalid user context"}`},
		{"wrong user type", "driver", "user-id", true, http.StatusInternalServerError, `{"error":"invalid user context"}`},
		{"nil user", "driver", (*data.User)(nil), true, http.StatusInternalServerError, `{"error":"invalid user context"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			router := gin.New()
			if tt.setUser {
				router.Use(func(c *gin.Context) {
					c.Set("user", tt.user)
					c.Next()
				})
			}
			handlerCalled := false
			router.GET("/protected", app.RequireRole(tt.role), func(c *gin.Context) {
				handlerCalled = true
				c.Status(http.StatusNoContent)
			})

			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/protected", nil))

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantStatus == http.StatusNoContent, handlerCalled)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			}
		})
	}
}
