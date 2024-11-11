package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type MockTodoStore struct {
	mock.Mock
}

func (m *MockTodoStore) GetAllTodos() ([]Todo, error) {
	args := m.Called()
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockTodoStore) GetTodo(id string) (*Todo, error) {
	args := m.Called(id)

	if todo := args.Get(0); todo != nil {
		return todo.(*Todo), args.Error(1)
	}

	// Trả về nil nếu không tìm thấy Todo
	return nil, args.Error(1)
}

func (m *MockTodoStore) CreateTodo(todo Todo) (*Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodoStore) UpdateTodo(id string, todo Todo) (*Todo, error) {
	args := m.Called(id, todo)
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodoStore) DeleteTodo(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockTodoStore) UpdateTodoStatus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/todo", h.GetAllTodos).Methods("GET")
	router.ServeHTTP(w, r)
}
func TestGetAllTodos(t *testing.T) {
	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}
	mockData := []Todo{
		{
			ID:        "1",
			Title:     "Test Todo",
			Desc:      "This is a test todo",
			Done:      false,
			CreatedAt: time.Now(),
			DoneAt:    nil,
		},
	}
	mockStore.On("GetAllTodos").Return(mockData, nil)
	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var actualTodos []Todo
	err = json.Unmarshal(rr.Body.Bytes(), &actualTodos)
	if err != nil {
		t.Fatal(err)
	}
	for i := range actualTodos {
		actualTodos[i].CreatedAt = time.Time{}
		actualTodos[i].DoneAt = nil
	}
	for i := range mockData {
		mockData[i].CreatedAt = time.Time{}
		mockData[i].DoneAt = nil
	}
	if !reflect.DeepEqual(actualTodos, mockData) {
		t.Errorf("expected body %+v, got %+v", mockData, actualTodos)
	}
	mockStore.AssertExpectations(t)
}
func TestGetTodo(t *testing.T) {

	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}

	t.Run("Test Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/todo/getuser/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.GetTodo(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Test Missing ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/todo/getuser/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.GetTodo(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Test Todo Not Found", func(t *testing.T) {
		mockStore.On("GetTodo", "not_found").Return(nil, fmt.Errorf("Todo not found"))
		req, err := http.NewRequest("GET", "/todo/getuser/not_found", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler.GetTodo(rr, req)

		// Kiểm tra mã lỗi HTTP trả về
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Kiểm tra nội dung phản hồi lỗi, loại bỏ newline
		expected := "Todo not found"
		actual := strings.TrimSpace(rr.Body.String()) // Loại bỏ newline hoặc khoảng trắng thừa
		assert.Equal(t, expected, actual)

		// Kiểm tra mock đã được gọi với ID "not_found"
		mockStore.AssertExpectations(t)
	})

	t.Run("Test Get Todo Successfully", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/todo/getuser/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Tạo một con trỏ đến đối tượng Todo
		mockData := &Todo{
			ID:        "1",
			Title:     "Test Todo",
			Desc:      "This is a test todo",
			Done:      false,
			CreatedAt: time.Date(2024, time.November, 8, 15, 45, 50, 681403600, time.Local),
			DoneAt:    nil,
		}

		mockStore.On("GetTodo", "1").Return(mockData, nil) // Trả về con trỏ mockData

		rr := httptest.NewRecorder()
		handler.GetTodo(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var todoResponse Todo
		err = json.NewDecoder(rr.Body).Decode(&todoResponse)
		if err != nil {
			t.Fatal(err)
		}

		// So sánh các trường thời gian
		if !todoResponse.CreatedAt.Equal(mockData.CreatedAt) {
			t.Errorf("Expected CreatedAt: %v, but got: %v", mockData.CreatedAt, todoResponse.CreatedAt)
		}

		// So sánh phần còn lại của đối tượng Todo
		assert.Equal(t, *mockData, todoResponse) // So sánh con trỏ *mockData với todoResponse
	})
}
func TestCreateTodo_Success(t *testing.T) {
	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}

	todoRequest := Todo{
		Title: "Test Todo",
		Desc:  "This is a test todo",
		Done:  false,
	}

	expectedTodo := &Todo{
		ID:        "1",
		Title:     "Test Todo",
		Desc:      "This is a test todo",
		Done:      false,
		CreatedAt: time.Now(),
	}

	mockStore.On("CreateTodo", mock.AnythingOfType("main.Todo")).Return(expectedTodo, nil)

	reqBody, _ := json.Marshal(todoRequest)
	req, err := http.NewRequest("POST", "/todo", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.CreateTodo(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var todoResponse Todo
	err = json.NewDecoder(rr.Body).Decode(&todoResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedTodo.Title, todoResponse.Title)
	assert.Equal(t, expectedTodo.Desc, todoResponse.Desc)
	assert.Equal(t, expectedTodo.Done, todoResponse.Done)

	assert.NotEqual(t, time.Time{}, todoResponse.CreatedAt)

	mockStore.AssertExpectations(t)
}
func TestUpdateTodo(t *testing.T) {
	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}

	t.Run("Test Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/todo/update/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		expected := "Method not allowed\n" // Kiểm tra cả ký tự xuống dòng
		assert.Equal(t, expected, rr.Body.String())
	})

	t.Run("Test Missing ID", func(t *testing.T) {
		req, err := http.NewRequest("PATCH", "/todo/update/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Test Invalid Request Body", func(t *testing.T) {
		req, err := http.NewRequest("PATCH", "/todo/update/1", bytes.NewBufferString("invalid json"))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request body\n", rr.Body.String())
	})

	t.Run("Test Todo Not Found", func(t *testing.T) {
		todoID := "not_found"
		mockStore.On("UpdateTodo", todoID, Todo{}).Return(&Todo{}, errors.New("not found"))

		req, err := http.NewRequest("PATCH", "/todo/update/"+todoID, bytes.NewBufferString("{}"))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Todo not found\n", rr.Body.String())
		mockStore.AssertExpectations(t)
	})

	t.Run("Test Failed to Update Todo", func(t *testing.T) {
		todoID := "1"
		mockStore.On("UpdateTodo", todoID, Todo{}).Return(&Todo{}, errors.New("todo not found"))

		req, err := http.NewRequest("PATCH", "/todo/update/"+todoID, bytes.NewBufferString("{}"))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Failed to update todo\n", rr.Body.String())
		mockStore.AssertExpectations(t)
	})

	t.Run("Test Update Todo Successfully", func(t *testing.T) {
		todoID := "1"
		fixedTime := time.Date(2024, time.November, 8, 16, 29, 3, 238919400, time.Local)

		updatedTodo := &Todo{
			ID:        todoID,
			Title:     "Updated Todo",
			Desc:      "This is an updated todo",
			Done:      true,
			CreatedAt: fixedTime,
			DoneAt:    nil,
		}

		mockStore.On("UpdateTodo", todoID, Todo{
			Title: "Updated Todo",
			Desc:  "This is an updated todo",
			Done:  true,
		}).Return(updatedTodo, nil)

		reqBody, err := json.Marshal(Todo{
			Title: "Updated Todo",
			Desc:  "This is an updated todo",
			Done:  true,
		})
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("PATCH", "/todo/update/"+todoID, bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.UpdateTodo(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var todoResponse Todo
		err = json.NewDecoder(rr.Body).Decode(&todoResponse)
		if err != nil {
			t.Fatal(err)
		}

		if !updatedTodo.CreatedAt.Equal(todoResponse.CreatedAt) {
			t.Errorf("Expected CreatedAt: %v, but got: %v", updatedTodo.CreatedAt, todoResponse.CreatedAt)
		}

		assert.Equal(t, updatedTodo.ID, todoResponse.ID)
		assert.Equal(t, updatedTodo.Title, todoResponse.Title)
		assert.Equal(t, updatedTodo.Desc, todoResponse.Desc)
		assert.Equal(t, updatedTodo.Done, todoResponse.Done)
		assert.Nil(t, todoResponse.DoneAt)

		mockStore.AssertExpectations(t)
	})

}
func TestUpdateTodoStatus(t *testing.T) {
	// Khởi tạo mockStore trong mỗi test case
	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}

	todoID := "1"
	updatedTodo := &Todo{
		ID:        todoID,
		Title:     "Updated Todo",
		Desc:      "This is an updated todo",
		Done:      true,
		CreatedAt: time.Now(),
		DoneAt:    nil,
	}

	t.Run("Success", func(t *testing.T) {
		// Reset mock expectations
		mockStore.ExpectedCalls = nil

		mockStore.On("UpdateTodoStatus", todoID).Return(nil)
		mockStore.On("GetTodo", todoID).Return(updatedTodo, nil)

		reqBody, err := json.Marshal(map[string]bool{"done": true})
		assert.NoError(t, err)

		req, err := http.NewRequest("PATCH", "/todo/update-status/"+todoID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.UpdateTodoStatus(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var todoResponse Todo
		err = json.NewDecoder(rr.Body).Decode(&todoResponse)
		assert.NoError(t, err)

		assert.Equal(t, updatedTodo.ID, todoResponse.ID)
		assert.Equal(t, updatedTodo.Title, todoResponse.Title)
		assert.Equal(t, updatedTodo.Desc, todoResponse.Desc)
		assert.Equal(t, updatedTodo.Done, todoResponse.Done)

		mockStore.AssertExpectations(t)
	})

	t.Run("Error_UpdateTodoStatus", func(t *testing.T) {
		// Reset mock expectations
		mockStore.ExpectedCalls = nil

		mockStore.On("UpdateTodoStatus", todoID).Return(assert.AnError)

		reqBody, err := json.Marshal(map[string]bool{"done": true})
		assert.NoError(t, err)

		req, err := http.NewRequest("PATCH", "/todo/update-status/"+todoID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.UpdateTodoStatus(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error updating todo status:")

		mockStore.AssertExpectations(t)
	})

	t.Run("Error_GetTodo", func(t *testing.T) {
		// Reset mock expectations
		mockStore.ExpectedCalls = nil

		mockStore.On("UpdateTodoStatus", todoID).Return(nil)
		mockStore.On("GetTodo", todoID).Return(nil, assert.AnError)

		reqBody, err := json.Marshal(map[string]bool{"done": true})
		assert.NoError(t, err)

		req, err := http.NewRequest("PATCH", "/todo/update-status/"+todoID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.UpdateTodoStatus(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error retrieving updated todo")

		mockStore.AssertExpectations(t)
	})
}

func TestDeleteTodo(t *testing.T) {
	todoID := "12345"

	mockStore := new(MockTodoStore)
	handler := &APIHandler{todoService: mockStore}

	t.Run("Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo/delete/1", nil)
		w := httptest.NewRecorder()

		handler.DeleteTodo(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	})

	t.Run("Error Todo Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/todo/delete/1", nil)
		w := httptest.NewRecorder()

		mockStore.On("DeleteTodo", "1").Return(fmt.Errorf("not found"))

		handler.DeleteTodo(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "not found\n", w.Body.String())

		mockStore.AssertExpectations(t)
	})

	t.Run("Successful_Deletion", func(t *testing.T) {
		mockStore.On("DeleteTodo", todoID).Return(nil)

		req, err := http.NewRequest("DELETE", "/todo/delete/"+todoID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.DeleteTodo(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())

		mockStore.AssertExpectations(t)
	})
}
