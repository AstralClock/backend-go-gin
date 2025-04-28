package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"log"
	"os"
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