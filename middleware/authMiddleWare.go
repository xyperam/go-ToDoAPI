package middleware

import (
	"context"
	"fmt"
	"go-web-server/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "userID" // Key untuk menyimpan userID dalam context

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil Authorization header dari request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Ambil token setelah "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		// Dekode token dan ambil klaim
		claims := jwt.MapClaims{} // Gunakan MapClaims untuk mengakses klaim secara dinamis
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verifikasi metode penandatanganan
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return utils.JWTKey, nil
		})

		// Cek validitas token
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		log.Println("Token Claims: ", claims)

		// Ambil userID dari klaim (sub) dan konversi sesuai format yang ada
		sub, ok := claims["sub"]
		if !ok {
			http.Error(w, "Missing user ID in token", http.StatusUnauthorized)
			return
		}

		var userID int
		switch v := sub.(type) {
		case float64:
			// Jika sub berupa float64, cast ke int
			userID = int(v)
		case string:
			// Jika sub berupa string, coba konversi ke int
			parsedID, err := strconv.Atoi(v)
			if err != nil {
				http.Error(w, "Invalid user ID format in token", http.StatusUnauthorized)
				return
			}
			userID = parsedID
		default:
			// Jika sub bukan float64 atau string
			http.Error(w, "Invalid user ID type in token", http.StatusUnauthorized)
			return
		}

		// Simpan userID ke dalam context untuk digunakan di handler berikutnya
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Panggil handler berikutnya dengan context yang sudah diperbarui
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
