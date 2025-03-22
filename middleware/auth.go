package middleware

import (
    "backend-go-gin/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil token dari header Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header diperlukan"})
            c.Abort()
            return
        }

        // Format: "Bearer <token>"
        tokenString := strings.Split(authHeader, " ")[1]
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
            c.Abort()
            return
        }

        // Verifikasi token
        claims, err := utils.VerifyToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
            c.Abort()
            return
        }

        // Simpan userID di context
        c.Set("userID", claims.UserID)
        c.Next()
    }
}