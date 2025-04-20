package main

import (
	"fmt"
	"go-web-server/routes"
	"go-web-server/utils"
	"net/http"
)

func main() {
	utils.UpdateTaskID() // Inisialisasi ID Task

	r := routes.SetupRoutes()

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
