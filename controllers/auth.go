package controllers

import (
	"encoding/json"
	"fmt"
	"go-web-server/models"
	"go-web-server/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.User

	//decode json dari req
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//validasi input
	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	//cek apakah user sudah ada
	var existingUser models.User
	if err := utils.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Tambahkan log untuk error lainnya
		fmt.Println("Error checking username:", err)
		http.Error(w, "Failed to check username", http.StatusInternalServerError)
		return
	}

	// Cek apakah email sudah ada
	var exisitingEmail models.User
	if err := utils.DB.Where("email = ?", input.Email).First(&exisitingEmail).Error; err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Tambahkan log untuk error lainnya
		fmt.Println("Error checking email:", err)
		http.Error(w, "Failed to check email", http.StatusInternalServerError)
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	input.Password = string(hashedPassword)
	// input.CreatedAt = time.Now() // Set waktu pembuatan user
	// input.UpdatedAt = time.Now() // Set waktu pembaruan user
	//simpan user ke db

	//buat user baru
	if err := utils.DB.Create(&input).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := utils.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	//bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	//Generate JWT
	token, err := utils.GenerateJWT(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
