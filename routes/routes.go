package routes

import (
	"go-web-server/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger) // middleware untuk logging

	r.Get("/tasks", controllers.GetTasks)
	r.Get("/tasks/{id}", controllers.GetTaskByID)
	r.Post("/tasks", controllers.CreateTask)
	r.Put("/tasks/{id}", controllers.UpdateTask)
	r.Delete("/tasks/{id}", controllers.DeleteTask)

	return r
}
