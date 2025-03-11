package controllers

import (
    "errors"
    "backend-go-gin/utils"
)

// Static admin user for testing
var adminUser = struct {
    Username string
    Password string
}{
    Username: "fufufafa",
    Password: "admin123",
}

// LoginUser authenticates the admin user and returns a JWT token
func LoginUser(Username, password string) (string, error) {
    if Username != adminUser.Username || password != adminUser.Password {
        return "", errors.New("invalid Username or password")
    }

    // Create JWT token
    token, err := utils.GenerateJWT(1) // Assuming user ID is 1 for the static admin
    if err != nil {
        return "", errors.New("could not create token")
    }

    return token, nil
}