package routes

import (
	"go-web-server/controllers"
	custommiddleware "go-web-server/middleware" // rename biar gak tabrakan dengan chi/middleware

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Logger global
	r.Use(middleware.Logger)

	// Public routes
	r.Post("/register", controllers.RegisterUser)
	r.Post("/login", controllers.LoginUser)

	// Protected routes dengan JWT buatan sendiri
	r.Group(func(protected chi.Router) {
		protected.Use(custommiddleware.JWTMiddleware)

		protected.Get("/tasks", controllers.GetTasks)
		protected.Get("/tasks/{id}", controllers.GetTaskByID)
		protected.Post("/tasks", controllers.CreateTask)
		protected.Put("/tasks/{id}", controllers.UpdateTask)
		protected.Delete("/tasks/{id}", controllers.DeleteTask)
	})

	return r
}
