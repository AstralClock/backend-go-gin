package handlers

import (
	"backend-go-gin/controllers"
	"backend-go-gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// RegisterHandler handles the user registration
func Register(c *gin.Context) {
	var request RegisterRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// Create user input model
	userInput := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// Call the controller function
	user, err := controllers.RegisterUser(userInput, request.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "User registered successfully",
		"data": gin.H{
			"email": user.Email,
			"id":    user.ID, // Assuming User model has an ID field
		},
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	token, err := controllers.LoginUser(input.Email, input.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Email atau password salah"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
