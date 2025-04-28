package migrations

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"fmt"
)

func Migrate() {
	// Daftar semua model yang akan dimigrasi
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.UserDetail{},
		&models.Produk{},
        &models.Order{},
        &models.OrderDetail{},
        &models.Cart{},
        &models.CartDetail{},
        &models.Ulasan{},
		&models.Payment{},
		// Tambahkan model lain di sini
	); err != nil {
		panic("Gagal melakukan migrasi database")
	}
	fmt.Println("Migrasi database berhasil!")
}
