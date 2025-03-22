package handlers

import (
	"backend-go-gin/controllers"
	"backend-go-gin/utils"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UserDetailHandler struct {
	userDetailController *controllers.UserDetailController
}

func NewUserDetailHandler() *UserDetailHandler {
	return &UserDetailHandler{
		userDetailController: &controllers.UserDetailController{},
	}
}

func (uh *UserDetailHandler) SaveUserDetail(c *gin.Context) {
	// Ambil user ID dari context (setelah autentikasi JWT)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	// Bind form data
	var request struct {
		Nama     string `form:"nama" binding:"required"`
		Telepon  string `form:"telepon" binding:"required"`
		Alamat   string `form:"alamat" binding:"required"`
		Kodepos  string `form:"kodepos" binding:"required"`
		Provinsi string `form:"provinsi" binding:"required"`
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, err := c.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengunggah gambar"})
		return
	}

	var imgPath string
	if file != nil {
		// Batasi ukuran file maksimal 3MB
		if file.Size > 3<<20 { // 3MB
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar (maksimal 3MB)"})
			return
		}

		// Generate nama file random
		randomName := utils.RandomString(10) // 10 karakter random
		ext := filepath.Ext(file.Filename)   // Ambil ekstensi file asli
		filename := filepath.Join("uploads", "avatar", randomName+ext)

		// Simpan file
		if err := c.SaveUploadedFile(file, filename); err != nil {
			log.Printf("Gagal menyimpan gambar: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
			return
		}
		imgPath = filename
	}

	// Panggil controller untuk menyimpan data
	userDetail, err := uh.userDetailController.SaveUserDetail(
		userID.(uint), // Konversi userID ke uint
		request.Nama,
		request.Telepon,
		request.Alamat,
		request.Kodepos,
		request.Provinsi,
		imgPath,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil disimpan",
		"data":    userDetail,
	})
}

func (uh *UserDetailHandler) UpdateUserDetail(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
        return
    }

    // Bind form data
    var request struct {
        Nama     *string `form:"nama"`
        Telepon  *string `form:"telepon"`
        Alamat   *string `form:"alamat"`
        Kodepos  *string `form:"kodepos"`
        Provinsi *string `form:"provinsi"`
    }
    if err := c.ShouldBind(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Handle file upload (jika ada gambar baru)
    file, err := c.FormFile("img")
    if err != nil && err != http.ErrMissingFile {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengunggah gambar"})
        return
    }

    var imgPath string
    if file != nil {
        // Batasi ukuran file maksimal 3MB
        if file.Size > 3<<20 { // 3MB
            c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar (maksimal 3MB)"})
            return
        }

        // Generate nama file random
        randomName := utils.RandomString(10) 
        ext := filepath.Ext(file.Filename)   
        filename := filepath.Join("uploads", "avatar", randomName+ext)

        // Simpan file baru
        if err := c.SaveUploadedFile(file, filename); err != nil {
            log.Printf("Gagal menyimpan gambar: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
            return
        }
        imgPath = filename
    }

    // Panggil controller untuk update data
    userDetail, err := uh.userDetailController.UpdateUserDetail(
        userID.(uint), // Konversi userID ke uint
        request.Nama,
        request.Telepon,
        request.Alamat,
        request.Kodepos,
        request.Provinsi,
        imgPath,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data berhasil diupdate",
        "data":    userDetail,
    })
}