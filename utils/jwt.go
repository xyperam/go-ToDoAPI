package utils

import (
	"fmt"
	"go-web-server/models"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte("secret_key") // Ganti dengan kunci rahasia yang lebih kuat
func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),         // Konversi userID ke string
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Token berlaku selama 24 jam
		IssuedAt:  jwt.NewNumericDate(time.Now()),     // Waktu pembuatan token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		log.Println("Failed to generate token:", err)
		return "", err
	}
	fmt.Println("Generated JWT:", tokenString)
	return tokenString, nil
}
