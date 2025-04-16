package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
var tasks = []Task{}
var taskID int = 1




func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	// r.Put("/tasks/{id}", updateTask)
	// r.Delete("/tasks/{id}", deleteTask)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

// Get Tasks


func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if err:= json.NewEncoder(w).Encode(tasks); err != nil{
		http.Error(w,"Error encoding response",http.StatusInternalServerError)
		return
	}
}

//Create Task
func createTask(w http.ResponseWriter, r *http.Request){
	var task Task

	if err:= json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w,"Invalid request payload",http.StatusBadRequest)
		return 
	}
	// json.NewDecoder(r.Body).Decode(&task)
	task.ID = taskID
	taskID++
	tasks = append(tasks,task)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task);err!=nil{
		http.Error(w,"Error encoding response",http.StatusInternalServerError)
		return
	}
}
