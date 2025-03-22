package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"backend-go-gin/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(input models.User, confirmPassword string) (models.User, error) {
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("email sudah terdaftar")
	}

	if input.Password != confirmPassword {
		return models.User{}, errors.New("password dan konfirmasi password tidak cocok")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New("gagal mengenkripsi password")
	}

	verifyToken := utils.GenerateRandomToken() // Buat fungsi ini di package utils

	// Buat user baru
	user := models.User{
		Email:       input.Email,
		Password:    string(hashedPassword),
		IsVerified:  false,
		VerifyToken: verifyToken,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return models.User{}, errors.New("gagal menyimpan user")
	}

	if err := utils.SendVerificationEmail(user.Email, verifyToken); err != nil {
		return models.User{}, errors.New("gagal mengirim email verifikasi")
	}

	return user, nil
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
