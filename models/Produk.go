package models

import "gorm.io/gorm"

type Produk struct {
    gorm.Model
    NamaProduk string  `json:"nama_produk" gorm:"not null"`
    Deskripsi  string  `json:"deskripsi" gorm:"not null"`
    Kategori   string  `json:"kategori" gorm:"not null"`
    Tag        string  `json:"tag" gorm:"not null"`
    Harga      float64 `json:"harga" gorm:"not null"`
    Jumlah     int     `json:"jumlah" gorm:"not null"`
    Image      string  `json:"image" gorm:"not null"`
}