package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thuanvo2405/GoNext/internal/data"
)

func (app *application) Authenticate(context *gin.Context) {
	parts := strings.Fields(context.Request.Header.Get("Authorization"))

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	userId, err := app.tokens.Verify(parts[1])
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	user, err := app.models.Users.GetByID(context.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "could not authenticate user",
			})
		}
		return
	}

	context.Set("user", user)
	context.Next()
}

func (app *application) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user")
		user, ok := value.(*data.User)
		if !exists || !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "invalid user context",
			})
			return
		}

		if user.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission",
			})
			return
		}

		c.Next()
	}
}
