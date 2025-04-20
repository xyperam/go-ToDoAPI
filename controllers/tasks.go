package controllers

import (
	"encoding/json"
	"go-web-server/models" // Mengimpor models untuk tipe Task
	"go-web-server/utils"  // Mengimpor utils untuk akses Tasks dan fungsi terkait
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Panggil LoadTasksFromFile dari utils untuk membaca data dari file
	if err := utils.LoadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	completedTasks := r.URL.Query().Get("completed")
	var filteredTasks []models.Task

	// Filter berdasarkan parameter "completed" jika ada
	if completedTasks != "" {
		completed, err := strconv.ParseBool(completedTasks)
		if err != nil {
			http.Error(w, "Invalid query parameter", http.StatusBadRequest)
			return
		}
		// Filter tasks berdasarkan status completed
		for _, task := range utils.Tasks {
			if task.Completed == completed {
				filteredTasks = append(filteredTasks, task)
			}
		}
	} else {
		// Jika tidak ada filter, ambil semua tasks
		filteredTasks = utils.Tasks
	}

	// Set header Content-Type untuk respons JSON
	w.Header().Set("Content-Type", "application/json")

	// Kirim response dalam format JSON
	if err := json.NewEncoder(w).Encode(filteredTasks); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Panggil LoadTasksFromFile untuk membaca data dari file
	if err := utils.LoadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	// Cari task berdasarkan ID
	for _, t := range utils.Tasks {
		if t.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(t); err != nil {
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
				return
			}
			return
		}
	}
}

// Create Task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task // Gunakan models.Task
	// Panggil LoadTasksFromFile untuk membaca data dari file
	if err := utils.LoadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}
	// Decode data dari body request
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Set ID baru untuk task
	task.ID = utils.TaskID
	utils.TaskID++
	// Tambahkan task ke dalam list
	utils.Tasks = append(utils.Tasks, task)

	// Simpan tasks ke file
	if err := utils.SaveTasksToFile(); err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Kirim response JSON
	json.NewEncoder(w).Encode(task)
}

// Update Task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Panggil LoadTasksFromFile untuk membaca data dari file
	if err := utils.LoadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

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
	for i, t := range utils.Tasks {
		if t.ID == taskID {
			utils.Tasks[i].Title = updatedTask.Title
			utils.Tasks[i].Description = updatedTask.Description
			utils.Tasks[i].Completed = updatedTask.Completed
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Simpan perubahan ke file
			if err := utils.SaveTasksToFile(); err != nil {
				http.Error(w, "Failed to Update Task", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Task Updated successfully"))
			return
		}
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Panggil LoadTasksFromFile untuk membaca data dari file
	if err := utils.LoadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	// Cari task berdasarkan ID dan hapus
	for i, t := range utils.Tasks {
		if t.ID == taskID {
			utils.Tasks = append(utils.Tasks[:i], utils.Tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Simpan perubahan ke file
			if err := utils.SaveTasksToFile(); err != nil {
				http.Error(w, "Failed to delete task", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Task deleted successfully"))
			return
		}
	}
}
