package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thuanvo2405/GoNext/internal/data"
)

func (app *application) updateDriverAvailability(c *gin.Context) {
	value, exists := c.Get("user")
	user, ok := value.(*data.User)
	if !exists || !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	var input struct {
		IsAvailable *bool `json:"is_available" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid availability input"})
		return
	}

	driver, err := app.models.DriverProfile.UpdateAvailability(c.Request.Context(), user.ID, *input.IsAvailable)
	if errors.Is(err, data.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver_profile": driver})
}
