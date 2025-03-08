package handlers

import (
	"backend-go-gin/controllers"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var input struct {
        Username    string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "Invalid data"})
        return
    }

    token, err := controllers.LoginUser(input.Username, input.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(200, gin.H{
        "token": token,
    })
}
