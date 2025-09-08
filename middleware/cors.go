package middleware

import (
    "github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        allowedOrigins := []string{
            "https://reva-baju.vercel.app",
            "https://dashboard-revabajuanak.vercel.app",
            "http://localhost:3000", // optional, buat dev
        }

        for _, o := range allowedOrigins {
            if origin == o {
                c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers",
            "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-CSRF-Token")
        c.Writer.Header().Set("Access-Control-Allow-Methods",
            "POST, OPTIONS, GET, PUT, DELETE, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
