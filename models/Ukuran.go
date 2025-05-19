package models

import "gorm.io/gorm"

type Ukuran struct {
    gorm.Model
    Ukuran   string    `gorm:"not null" json:"ukuran"`
    Produks  []Produk  `gorm:"many2many:produk_ukurans;"` // Tambahkan relasi balik
}