package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thuanvo2405/GoNext/internal/data"
	"github.com/thuanvo2405/GoNext/internal/password"
)

func (app *application) registerRiderProfile(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid registration input",
		})
		return
	}

	if len(input.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "your password too long",
		})
		return
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not process registration",
		})
		return
	}

	user := &data.User{
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword,
		Role:     "rider",
	}

	err = app.models.Users.Insert(c.Request.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"message": "Email is already, please try another"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong!"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Create Profile Successful",
		"user":    user,
	})
}

func (app *application) registerDriverProfile(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		VehicleType string `json:"vehicle_type" binding:"required,oneof=motorcycle car electric_car"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid registration input",
		})
		return
	}

	if len(input.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "your password too long",
		})
		return
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not process registration",
		})
		return
	}

	user := &data.User{
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword,
		Role:     "driver",
	}

	driver := &data.DriverProfile{
		VehicleType: input.VehicleType,
	}

	if err := app.models.DriverProfile.Insert(c.Request.Context(), user, driver); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"message": "Email is already, please try another"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong!"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Create Profile Successful",
		"user":    user,
	})
}

func (app *application) loginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login input",
		})
		return
	}

	if len(input.Password) > 72 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	user, err := app.models.Users.GetByEmail(c.Request.Context(), strings.ToLower(input.Email))
	if errors.Is(err, data.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not process login",
		})
		return
	}

	if !password.Verify(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	accessToken, err := app.tokens.Sign(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "credentials verified",
		"user":    user,
		"token":   accessToken,
	})
}

func (app *application) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString("userId"),
	})
}
