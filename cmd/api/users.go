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
