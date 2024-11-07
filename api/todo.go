package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"sync"
	"time"
)

type Todo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	Done      bool       `json:"done"`
	CreatedAt time.Time  `json:"created_at"`
	DoneAt    *time.Time `json:"done_at"`
}

type TodoService interface {
	GetAllTodos() ([]Todo, error)
	GetTodo(id string) (*Todo, error)
	CreateTodo(todo Todo) (*Todo, error)
	UpdateTodo(id string, todo Todo) (*Todo, error)
	DeleteTodo(id string) error
	UpdateTodoStatus(id string) error
}

type DbTodoService struct {
	db *Db
	mu sync.Mutex
}

func NewDbTodoService(db *Db) *DbTodoService {
	return &DbTodoService{
		db: db,
	}
}
func generateNewID() string {
	return uuid.New().String()
}

func (s *DbTodoService) GetAllTodos() ([]Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.db.conn.Query(ctx, "SELECT id, title, description, done, created_at FROM todo ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("truy vấn thất bại: %v", err)
	}

	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Desc, &todo.Done, &todo.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan thất bại: %v", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi sau khi đọc rows: %v", err)
	}
	return todos, nil
}
func (s *DbTodoService) GetTodo(id string) (*Todo, error) {
	var todo Todo
	err := s.db.conn.QueryRow(context.Background(), "SELECT id, title, description, done, created_at, done_at FROM todo WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Desc, &todo.Done, &todo.CreatedAt, &todo.DoneAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found users")
		}
		return nil, fmt.Errorf("truy vấn thất bại: %v", err)
	}
	return &todo, nil
}
func (s *DbTodoService) CreateTodo(todo Todo) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo.ID = generateNewID()
	todo.Done = false
	_, err := s.db.conn.Exec(context.Background(),
		"INSERT INTO todo (id, title, description, done, created_at) VALUES ($1, $2, $3, $4, $5)",
		todo.ID, todo.Title, todo.Desc, todo.Done, time.Now())
	if err != nil {
		return nil, fmt.Errorf("thêm todo thất bại: %v", err)
	}
	return &todo, nil
}
func (s *DbTodoService) UpdateTodo(id string, todo Todo) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var doneAt *time.Time

	if todo.Done {
		now := time.Now()
		doneAt = &now
	} else {
		doneAt = nil
	}
	_, err := s.db.conn.Exec(context.Background(),
		"UPDATE todo SET title = $1, description = $2, done = $3, done_at = $4 WHERE id = $5",
		todo.Title, todo.Desc, todo.Done, doneAt, id) // Thêm doneAt vào câu lệnh
	if err != nil {
		return nil, fmt.Errorf("cập nhật todo thất bại: %v", err)
	}

	var updatedTodo Todo
	err = s.db.conn.QueryRow(context.Background(),
		"SELECT id, title, description, done, created_at, done_at FROM todo WHERE id = $1", id).
		Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Desc, &updatedTodo.Done, &updatedTodo.CreatedAt, &updatedTodo.DoneAt)

	if err != nil {
		return nil, fmt.Errorf("lấy todo đã cập nhật thất bại: %v", err)
	}
	return &updatedTodo, nil
}
func (s *DbTodoService) UpdateTodoStatus(id string) error {
	var currentDone bool
	err := s.db.conn.QueryRow(context.Background(), "SELECT done FROM todo WHERE id = $1", id).Scan(&currentDone)
	if err != nil {
		return fmt.Errorf("không tìm thấy todo với id %s: %v", id, err)
	}

	newDone := !currentDone
	var doneAt interface{}

	if newDone {
		doneAt = time.Now()
	} else {
		doneAt = nil
	}
	_, err = s.db.conn.Exec(context.Background(), "UPDATE todo SET done = $1, done_at = $2 WHERE id = $3", newDone, doneAt, id)
	if err != nil {
		return fmt.Errorf("cập nhật trạng thái todo thất bại: %v", err)
	}

	return nil
}
func (s *DbTodoService) DeleteTodo(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var exists bool
	err := s.db.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM todo WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("kiểm tra sự tồn tại của todo thất bại: %v", err)
	}
	if !exists {
		return fmt.Errorf("not found ID")
	}
	_, err = s.db.conn.Exec(context.Background(), "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("xóa todo thất bại: %v", err)
	}
	return nil
}
