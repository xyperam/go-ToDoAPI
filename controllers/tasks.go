package controllers

import (
	"encoding/json"
	"fmt"
	"go-web-server/middleware"
	"go-web-server/models" // Mengimpor models untuk tipe Task
	"go-web-server/utils"  // Mengimpor utils untuk akses Tasks dan fungsi terkait
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context
	userIDCtx := r.Context().Value(middleware.UserIDKey)
	if userIDCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDCtx.(int)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusUnauthorized)
		return
	}

	// Ambil parameter query "completed" jika ada
	completedTasks := r.URL.Query().Get("completed")
	query := utils.DB

	var tasks []models.Task

	// Filter berdasarkan parameter "completed" jika ada
	if completedTasks != "" {
		// Konversi parameter "completed" ke bool
		completed, err := strconv.ParseBool(completedTasks)
		if err != nil {
			http.Error(w, "Invalid query parameter 'completed', must be true or false", http.StatusBadRequest)
			return
		}
		query = query.Where("completed = ?", completed)
	}

	// Filter berdasarkan userID
	query = query.Where("user_id = ?", userID)

	// Ambil data dari database
	if err := query.Find(&tasks).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve tasks: %v", err), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type untuk respons JSON
	w.Header().Set("Content-Type", "application/json")

	// Kirim data tasks sebagai JSON
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode tasks into JSON", http.StatusInternalServerError)
		return
	}
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Panggil LoadTasksFromFile untuk membaca data dari file
	var task models.Task
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}

	if err := utils.DB.First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		}
		return
	}
	// Cari task berdasarkan ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
		return
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	userIDCtx := r.Context().Value(middleware.UserIDKey)
	if userIDCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDCtx.(int)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusUnauthorized)
		return
	}

	// Decode body
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi title
	if task.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	fmt.Println("userID from context:", userID)
	// Set user ID ke task
	task.UserID = userID

	// Set timestamps
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Cek user valid atau tidak
	var user models.User
	if err := utils.DB.First(&user, "id = ?", task.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to validate user", http.StatusInternalServerError)
		return
	}

	// Simpan task
	if err := utils.DB.Create(&task).Error; err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Update Task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Convert id ke int
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	// Decode body request
	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Cari dan update task berdasarkan ID
	var task models.Task
	if err := utils.DB.First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		}
		return
	}
	//update task
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.Completed = updatedTask.Completed
	if err := utils.DB.Save(&task).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Kirim response JSON
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Panggil LoadTasksFromFile untuk membaca data dari file
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	// Cari task berdasarkan ID dan hapus
	var task models.Task
	if err := utils.DB.First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		}
		return
	}
	if err := utils.DB.Delete(&task).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // Set status 204 No Content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"}) // Kirim response JSON
}
