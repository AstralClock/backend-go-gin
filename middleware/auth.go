package middleware

import (
    "backend-go-gin/utils"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "Token tidak ditemukan"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(utils.JWT_SECRET), nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Token tidak valid"})
            c.Abort()
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}