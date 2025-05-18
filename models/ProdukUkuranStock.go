package models

import _ "gorm.io/gorm"

type ProdukUkuran struct {
    ProdukID  uint `gorm:"primaryKey"`
    UkuranID  uint `gorm:"primaryKey"`
    Stok      int  `gorm:"not null"`
	Produk   Produk `gorm:"foreignKey:ProdukID"`
    Ukuran   Ukuran `gorm:"foreignKey:UkuranID"`
}