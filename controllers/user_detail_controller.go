package controllers

import (
    "backend-go-gin/config"
    "backend-go-gin/models"
    "log"
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