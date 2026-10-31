package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

	context.Set("userId", userId)
	context.Next()
}
