package middleware

import (
    "backend-go-gin/utils"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "strings"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header diperlukan"})
            c.Abort()
            return
        }

        tokenString := strings.Split(authHeader, " ")[1]
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
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