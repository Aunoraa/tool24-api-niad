// @title Todo API
// @version 1.0
// @description This is a sample API for managing todos.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /todo

package main

import (
	_ "api/docs"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

// main function initializes the server and routes.
func main() {
	db, err := NewDb()
	if err != nil {
		fmt.Printf("Lỗi khi khởi tạo cơ sở dữ liệu: %v\n", err)
		return
	}
	defer db.Close()

	todoService := NewDbTodoService(db)
	apiHandler := NewAPIHandler(todoService)

	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/todo", apiHandler.GetAllTodos).Methods(http.MethodGet)
	router.HandleFunc("/todo/getuser/{id}", apiHandler.GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/todo/create", apiHandler.CreateTodo).Methods(http.MethodPost)
	router.HandleFunc("/todo/update/{id}", apiHandler.UpdateTodo).Methods(http.MethodPatch)
	router.HandleFunc("/todo/update-status/{id}", apiHandler.UpdateTodoStatus).Methods(http.MethodPatch)
	router.HandleFunc("/todo/delete/{id}", apiHandler.DeleteTodo).Methods(http.MethodDelete)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Thay bằng domain của frontend
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
