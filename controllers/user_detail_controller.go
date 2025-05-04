package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserDetailController struct{}

func (uc *UserDetailController) SaveUserDetail(userID uint, nama, telepon, alamat, kodepos, provinsi, imgPath string) (*models.UserDetail, error) {
	userDetail := models.UserDetail{
		UserID:   userID,
		Nama:     nama,
		Telepon:  telepon,
		Alamat:   alamat,
		Kodepos:  kodepos,
		Provinsi: provinsi,
		Img:      imgPath,
	}

	if err := config.DB.Create(&userDetail).Error; err != nil {
		log.Printf("Gagal menyimpan data user detail: %v", err)
		return nil, err
	}

	return &userDetail, nil
}

func (uc *UserDetailController) UpdateUserDetail(userID uint, nama, telepon, alamat, kodepos, provinsi *string, imgPath string) (*models.UserDetail, error) {
    var userDetail models.UserDetail
    if err := config.DB.Where("user_id = ?", userID).First(&userDetail).Error; err != nil {
        return nil, err
    }

    // Update field yang diubah
    if nama != nil {
        userDetail.Nama = *nama
    }
    if telepon != nil {
        userDetail.Telepon = *telepon
    }
    if alamat != nil {
        userDetail.Alamat = *alamat
    }
    if kodepos != nil {
        userDetail.Kodepos = *kodepos
    }
    if provinsi != nil {
        userDetail.Provinsi = *provinsi
    }

    if imgPath != "" {
        if userDetail.Img != "" {
            if err := os.Remove(userDetail.Img); err != nil {
                log.Printf("Gagal menghapus gambar lama: %v", err)
            }
        }
        userDetail.Img = imgPath
    }

    if err := config.DB.Save(&userDetail).Error; err != nil {
        return nil, err
    }

    return &userDetail, nil
}

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) GetAllUserDetails(c *gin.Context) {
	var userDetails []struct {
		models.UserDetail
		Email string `json:"email"`
	}

	// Query dengan join ke tabel users
	if err := config.DB.
		Table("user_details").
		Select("user_details.*, users.email").
		Joins("LEFT JOIN users ON users.id = user_details.user_id").
		Find(&userDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details retrieved successfully",
		"data":    userDetails,
	})
}

func (uc *UserController) GetUserDetailByID(c *gin.Context) {
	userID := c.Param("id")

	var result struct {
		models.UserDetail
		Email string `json:"email"`
	}

	// Query dengan join
	if err := config.DB.
		Table("user_details").
		Select("user_details.*, users.email").
		Joins("LEFT JOIN users ON users.id = user_details.user_id").
		Where("user_details.user_id = ?", userID).
		First(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User detail not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User detail retrieved successfully",
		"data":    result,
	})
}