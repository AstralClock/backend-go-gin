package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"backend-go-gin/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(input models.User, confirmPassword string) (models.User, string, error) {
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return models.User{}, "", errors.New("email sudah terdaftar")
	}

	if input.Password != confirmPassword {
		return models.User{}, "", errors.New("password dan konfirmasi password tidak cocok")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", errors.New("gagal mengenkripsi password")
	}

	verifyToken := utils.GenerateRandomToken() 

	user := models.User{
		Email:       input.Email,
		Password:    string(hashedPassword),
		IsVerified:  false,
		VerifyToken: verifyToken,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return models.User{}, "", errors.New("gagal menyimpan user")
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return models.User{}, "", errors.New("gagal membuat token JWT")
	}

	if err := utils.SendVerificationEmail(user.Email, verifyToken); err != nil {
		return models.User{}, "", errors.New("gagal mengirim email verifikasi")
	}

	return user, jwtToken, nil
}

func LoginUser(email, password string) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("password salah")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New("gagal membuat token")
	}

	return token, nil
}

var adminUser = struct {
    Username string
    Password string
}{
    Username: "fufufafa",
    Password: "admin123",
}

func LoginAdmin(Username, password string) (string, error) {
    if Username != adminUser.Username || password != adminUser.Password {
        return "", errors.New("invalid Username or password")
    }

    // Create JWT token
    token, err := utils.GenerateJWT(1) // Assuming user ID is 1 for the static admin
    if err != nil {
        return "", errors.New("could not create token")
    }

    return token, nil
}
