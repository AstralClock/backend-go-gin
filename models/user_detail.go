package models

import "gorm.io/gorm"

type UserDetail struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
    UserID   uint   `gorm:"not null;unique" json:"user_id"` // Foreign key ke tabel users
    Nama     string `gorm:"not null" json:"nama"`
    Telepon  string `gorm:"not null" json:"telepon"`
    Alamat   string `gorm:"not null" json:"alamat"`
    Kodepos  string `gorm:"not null" json:"kodepos"`
    Provinsi string `gorm:"not null" json:"provinsi"`
    Img      string `gorm:"null" json:"img"` // Path atau URL gambar
}