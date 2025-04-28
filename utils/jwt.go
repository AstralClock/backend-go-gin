package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

var JWT_SECRET = []byte("dkfjoiqhfnkdsj;fjds")

func GenerateJWT(userID uint) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JWT_SECRET)
}

func VerifyToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return JWT_SECRET, nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }
    return claims, nil
}