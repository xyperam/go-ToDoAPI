package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

const dataFile = "tasks.json"

var tasks = []Task{}
var taskID int = 1

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTasks)
	r.Get("/tasks/{id}", getTaskByID)
	r.Post("/tasks", createTask)
	r.Put("/tasks/{id}", updateTask)
	r.Delete("/tasks/{id}", deleteTask)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

func loadTasksFromFile() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &tasks)
}

func saveTasksToFile() error {
	data, err := json.MarshalIndent(tasks, " ", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// Get Tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := loadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	if err := loadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	for _, t := range tasks {
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
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	//check baca data dari file
	if err := loadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	} //decode body
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// json.NewDecoder(r.Body).Decode(&task)
	task.ID = taskID
	taskID++
	tasks = append(tasks, task)

	if err := saveTasksToFile(); err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Update Task
func updateTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	//check baca data dari file
	if err := loadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	//convert id to int
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	//decode body
	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == taskID {
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].Completed = updatedTask.Completed
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := saveTasksToFile(); err != nil {
				http.Error(w, "Failed to Update Task", http.StatusInternalServerError)
				return
			}
			return
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	if err := loadTasksFromFile(); err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := saveTasksToFile(); err != nil {
				http.Error(w, "Failed to delete task", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Task deleted successfully"))
			return
		}
	}
}
