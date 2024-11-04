package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

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

	router.HandleFunc("/todo", apiHandler.GetAllTodos).Methods(http.MethodGet)
	router.HandleFunc("/todo/getuser/{id}", apiHandler.GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/todo/create", apiHandler.CreateTodo).Methods(http.MethodPost)
	router.HandleFunc("/todo/update/{id}", apiHandler.UpdateTodo).Methods(http.MethodPatch)
	router.HandleFunc("/todo/update-status/{id}", apiHandler.UpdateTodoStatus).Methods(http.MethodPatch)
	router.HandleFunc("/todo/delete/{id}", apiHandler.DeleteTodo).Methods(http.MethodDelete)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Thay bằng domain của frontend
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", corsHandler)
}
