package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	err := router.Run() // listens on 0.0.0.0:8080 by default
	if err != nil {
		log.Fatal(err)
	}
}
