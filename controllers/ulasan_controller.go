package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UlasanController struct{}

func NewUlasanController() *UlasanController {
	return &UlasanController{}
}

type CreateUlasanRequest struct {
	Ulasan string `json:"ulasan" binding:"required"`
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
}
func (uc *UlasanController) CreateUlasan(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID Produk tidak valid"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terotentikasi"})
		return
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format ID Pengguna tidak valid di token"})
		return
	}

	var request CreateUlasanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var produk models.Produk
	if err := config.DB.First(&produk, uint(productID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi produk: " + err.Error()})
		return
	}

	var existingUlasan models.Ulasan
	err = config.DB.Where("produk_id = ? AND user_id = ?", productID, userID).First(&existingUlasan).Error
	if err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda sudah pernah memberikan ulasan untuk produk ini"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi ulasan sebelumnya: " + err.Error()})
		return
	}


	// TODO (Opsional Lanjutan): Cek apakah pengguna pernah membeli produk ini.
	// Ini memerlukan query ke tabel Order dan OrderDetail.
	// Misal:
	var orderDetail models.OrderDetail
	err = config.DB.Joins("JOIN orders ON orders.id = order_details.order_id").
		Where("orders.user_id = ? AND order_details.produk_id = ?", userID, productID).
		First(&orderDetail).Error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda harus membeli produk ini terlebih dahulu untuk memberikan ulasan."})
		return
	}


	ulasan := models.Ulasan{
		ProdukID: uint(productID),
		UserID:   userID,
		Ulasan:   request.Ulasan,
		Rating:   request.Rating,
	}

	if err := config.DB.Create(&ulasan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ulasan: " + err.Error()})
		return
	}

	if err := config.DB.Preload("User").Preload("Produk").First(&ulasan, ulasan.ID).Error; err != nil {
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ulasan berhasil dibuat",
		"ulasan":  ulasan,
	})
}


func (uc *UlasanController) GetUlasanByProductID(c *gin.Context) {
    productIDStr := c.Param("id")
    productID, err := strconv.ParseUint(productIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID Produk tidak valid"})
        return
    }


    var ulasan []models.Ulasan
    if err := config.DB.Where("produk_id = ?", productID).
        Preload("User").
        Preload("User.UserDetail").
        Find(&ulasan).Error; err != nil {
        // Debug 2: Log error database
        fmt.Printf("Error database: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan ulasan: " + err.Error()})
        return
    }

    // Jika TIDAK ADA ulasan
    if len(ulasan) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "data": nil, // Pastikan ini null
            "meta": gin.H{
                "status":  "success",
                "message": "Tidak ada ulasan untuk produk ini",
            },
        })
        return
    }

    // Jika ADA ulasan
    var response []gin.H
    for _, u := range ulasan {
        userData := gin.H{"id": u.User.ID, "email": u.User.Email}

        if u.User.UserDetail.UserID != 0 {
            userData["nama"] = u.User.UserDetail.Nama
            userData["telepon"] = u.User.UserDetail.Telepon
        }

        response = append(response, gin.H{
            "id":         u.ID,
            "produk_id":  u.ProdukID,
            "user":       userData,
            "rating":     u.Rating,
            "komentar":   u.Ulasan,
            "created_at": u.CreatedAt,
            "updated_at": u.UpdatedAt,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "data": response,
        "meta": gin.H{
            "status":  "success",
            "message": "Data ulasan berhasil didapatkan",
        },
    })
}