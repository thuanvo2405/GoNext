package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.POST("/auth/register/rider", app.registerRiderProfile)

	return router
}
