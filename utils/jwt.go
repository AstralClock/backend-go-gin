package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("dkfjoiqhfnkdsj;fjds")

func GenerateJWT(userID uint) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(JWT_SECRET)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}