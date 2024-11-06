package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockTodoService struct {
	todos []Todo
}

func (m *MockTodoService) GetAllTodos() ([]Todo, error) {
	return m.todos, nil
}

func (m *MockTodoService) GetTodo(id string) (*Todo, error) {
	for _, todo := range m.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockTodoService) CreateTodo(todo Todo) (*Todo, error) {
	todo.ID = fmt.Sprintf("%d", len(m.todos)+1)
	m.todos = append(m.todos, todo)
	return &todo, nil
}
func (m *MockTodoService) UpdateTodo(id string, todo Todo) (*Todo, error) {
	for i, existing := range m.todos {
		if existing.ID == id {
			m.todos[i] = todo
			return &todo, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockTodoService) DeleteTodo(id string) error {
	for i, todo := range m.todos {
		if todo.ID == id {
			// Xóa todo
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo not found")
}

func (m *MockTodoService) UpdateTodoStatus(id string) error {
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos[i].Done = !todo.Done
			return nil
		}
	}
	return errors.New("not found")
}

func TestAPIHandler_UpdateTodoStatus(t *testing.T) {
	mockService := &MockTodoService{
		todos: []Todo{
			{ID: "1", Title: "Task 1", Desc: "Description 1", Done: false, CreatedAt: time.Now()},
			{ID: "2", Title: "Task 2", Desc: "Description 2", Done: true, CreatedAt: time.Now()},
		},
	}

	handler := NewAPIHandler(mockService)

	reqBody := map[string]bool{"done": true}
	reqBodyJSON, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("PATCH", "/todo/update-status/1", bytes.NewBuffer(reqBodyJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateTodoStatus(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	updatedTodo, err := mockService.GetTodo("1")
	assert.NoError(t, err)

	assert.True(t, updatedTodo.Done)
}

func TestAPIHandler_GetTodo(t *testing.T) {
	mockService := &MockTodoService{
		todos: []Todo{
			{ID: "1", Title: "Task 1", Desc: "Description 1", Done: false, CreatedAt: time.Now()},
		},
	}

	handler := NewAPIHandler(mockService)

	req, err := http.NewRequest("GET", "/todo/getuser/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetTodo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedTodo := mockService.todos[0] // Lấy đối tượng Todo mong đợi

	var actualTodo Todo
	err = json.Unmarshal(rr.Body.Bytes(), &actualTodo)
	assert.NoError(t, err)

	assert.Equal(t, expectedTodo.ID, actualTodo.ID)
	assert.Equal(t, expectedTodo.Title, actualTodo.Title)
	assert.Equal(t, expectedTodo.Desc, actualTodo.Desc)
	assert.Equal(t, expectedTodo.Done, actualTodo.Done)
	assert.True(t, expectedTodo.CreatedAt.Equal(actualTodo.CreatedAt), "CreatedAt times are not equal")
}

func TestAPIHandler_CreateTodo(t *testing.T) {
	mockService := &MockTodoService{
		todos: []Todo{},
	}

	handler := NewAPIHandler(mockService)

	newTodo := Todo{
		Title: "Test Todo",
		Desc:  "This is a test todo",
		Done:  false,
	}

	newTodoJSON, err := json.Marshal(newTodo)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/todo/create", bytes.NewBuffer(newTodoJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateTodo(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdTodo Todo
	err = json.Unmarshal(rr.Body.Bytes(), &createdTodo)
	assert.NoError(t, err)

	assert.Equal(t, newTodo.Title, createdTodo.Title)
	assert.Equal(t, newTodo.Desc, createdTodo.Desc)
	assert.Equal(t, newTodo.Done, createdTodo.Done)

	assert.NotEmpty(t, createdTodo.ID)

	currentTime := time.Now()

	assert.True(t, createdTodo.CreatedAt.After(currentTime.Add(-1*time.Second)), "CreatedAt should be after current time")
}

func TestAPIHandler_UpdateTodo(t *testing.T) {
	mockService := &MockTodoService{
		todos: []Todo{
			{ID: "1", Title: "Initial Todo", Desc: "This is the initial todo", Done: false, CreatedAt: time.Now()},
		},
	}

	handler := NewAPIHandler(mockService)

	// Todo đã cập nhật
	updatedTodo := Todo{
		ID:    "1",
		Title: "Updated Todo",
		Desc:  "This is an updated todo",
		Done:  true,
	}

	updatedTodoJSON, err := json.Marshal(updatedTodo)
	assert.NoError(t, err)

	req, err := http.NewRequest("PATCH", "/todo/update/1", bytes.NewBuffer(updatedTodoJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.UpdateTodo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resultTodo Todo
	err = json.Unmarshal(rr.Body.Bytes(), &resultTodo)
	assert.NoError(t, err)

	assert.Equal(t, updatedTodo.Title, resultTodo.Title)
	assert.Equal(t, updatedTodo.Desc, resultTodo.Desc)
	assert.Equal(t, updatedTodo.Done, resultTodo.Done)
}
func TestAPIHandler_DeleteTodo(t *testing.T) {
	mockService := &MockTodoService{
		todos: []Todo{
			{ID: "1", Title: "Test Todo", Desc: "This is a test todo", Done: false, CreatedAt: time.Now()},
		},
	}

	handler := NewAPIHandler(mockService)

	req, err := http.NewRequest("DELETE", "/todo/delete/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.DeleteTodo(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code) // Thay đổi từ http.StatusOK sang http.StatusNoContent

	req, err = http.NewRequest("GET", "/todo/getuser/1", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.GetTodo(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
