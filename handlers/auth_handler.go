package handlers

import (
	"backend-go-gin/config"
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


func Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	userInput := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user, token, err := controllers.RegisterUser(userInput, request.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "User registered successfully",
		"data": gin.H{
			"email": user.Email,
			"id":    user.ID,
			"token": token,
		},
	})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	var user models.User
	if err := config.DB.Where("verify_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Token tidak valid",
		})
		return
	}

	// Update status verifikasi
	user.IsVerified = true
	user.VerifyToken = "" // Hapus token setelah verifikasi
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal memverifikasi email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Email berhasil diverifikasi",
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

func LoginAdmin(c *gin.Context) {
    var input struct {
        Username    string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "Invalid data"})
        return
    }

    token, err := controllers.LoginAdmin(input.Username, input.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(200, gin.H{
        "token": token,
    })
}