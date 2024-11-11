package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type APIHandler struct {
	todoService TodoService
}

var ErrTodoNotFound = errors.New("todo not found")

func NewAPIHandler(todoService TodoService) *APIHandler {
	return &APIHandler{
		todoService: todoService,
	}
}

// @Summary Lấy tất cả các Todo
// @Description Lấy danh sách tất cả các Todo
// @Tags Todos
// @Produce  json
// @Success 200 {array} Todo
// @Router /todo [get]
func (h *APIHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	todos, err := h.todoService.GetAllTodo(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching todos: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// @Summary Lấy một Todo theo ID
// @Description Lấy thông tin chi tiết của một Todo theo ID
// @Tags Todos
// @Produce  json
// @Param id path string true "ID của Todo"
// @Success 200 {object} Todo
// @Failure 400 {string} string "ID không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy Todo"
// @Router /todo/getuser/{id} [get]
func (h *APIHandler) GetTodo(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/todo/getuser/")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetTodo(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
	}
}

// @Summary Tạo một Todo mới
// @Description Tạo một Todo mới với thông tin được cung cấp
// @Tags Todos
// @Accept  json
// @Produce  json
// @Param todo body Todo true "Thông tin của Todo"
// @Success 201 {object} Todo
// @Failure 400 {string} string "Request body không hợp lệ"
// @Router /todo/create [post]
func (h *APIHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo.CreatedAt = time.Now()
	newTodo, err := h.todoService.CreateTodo(ctx, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// @Summary Cập nhật một Todo
// @Description Cập nhật thông tin của một Todo theo ID
// @Tags Todos
// @Accept  json
// @Produce  json
// @Param id path string true "ID của Todo"
// @Param todo body Todo true "Thông tin cập nhật của Todo"
// @Success 200 {object} Todo
// @Failure 400 {string} string "Request body không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy Todo"
// @Router /todo/update/{id} [patch]
func (h *APIHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/todo/update/")
	if id == "" {
		http.Error(w, "ID is required", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTodo, err := h.todoService.UpdateTodo(ctx, id, todo)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTodo)
}

// @Summary Cập nhật trạng thái của một Todo
// @Description Cập nhật trạng thái của một Todo theo ID
// @Tags Todos
// @Produce  json
// @Param id path string true "ID của Todo"
// @Success 200 {object} Todo
// @Failure 404 {string} string "Không tìm thấy Todo"
// @Router /todo/update-status/{id} [patch]
func (h *APIHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/todo/update-status/"):]
	if id == "" {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	err := h.todoService.UpdateTodoStatus(ctx, id)
	if err != nil {
		http.Error(w, "Error updating todo status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todo, err := h.todoService.GetTodo(ctx, id)
	if err != nil {
		http.Error(w, "Error retrieving updated todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Xóa một Todo
// @Description Xóa một Todo theo ID
// @Tags Todos
// @Param id path string true "ID của Todo"
// @Success 204 {string} string "Xóa Todo thành công"
// @Failure 400 {string} string "ID không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy Todo"
// @Router /todo/delete/{id} [delete]
func (h *APIHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/todo/delete/"):]
	if id == "" {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	err := h.todoService.DeleteTodo(ctx, id)
	if err != nil {
		log.Println("Error deleting todo:", err)

		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
