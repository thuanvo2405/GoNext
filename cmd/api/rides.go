package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thuanvo2405/GoNext/internal/geo"
	"github.com/thuanvo2405/GoNext/internal/pricing"
)

type coordinateInput struct {
	Lat *float64 `json:"lat" binding:"required,min=-90,max=90"`
	Lng *float64 `json:"lng" binding:"required,min=-180,max=180"`
}

func (app *application) estimateFareHandler(c *gin.Context) {
	var input struct {
		Pickup      coordinateInput `json:"pickup" binding:"required"`
		Destination coordinateInput `json:"destination" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid estimate input",
		})
		return
	}

	if *input.Pickup.Lat == *input.Destination.Lat && *input.Pickup.Lng == *input.Destination.Lng {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pickup and destination must be different",
		})
		return
	}

	distanceKm := geo.HaversineKm(geo.Coordinate{
		Lat: *input.Pickup.Lat,
		Lng: *input.Pickup.Lng,
	}, geo.Coordinate{
		Lat: *input.Destination.Lat,
		Lng: *input.Destination.Lng,
	})

	c.JSON(http.StatusOK, gin.H{
		"distance_km":     distanceKm,
		"distance_source": "straight_line",
		"estimates":       pricing.EstimateFares(distanceKm),
	})
}
