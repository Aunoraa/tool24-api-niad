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

	_, err := s.db.conn.Exec(context.Background(),
		"UPDATE todo SET title = $1, description = $2, done = $3 WHERE id = $4",
		todo.Title, todo.Desc, todo.Done, id)

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
func (s *DbTodoService) DeleteTodo(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.conn.Exec(context.Background(), "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("xóa todo thất bại: %v", err)
	}

	return nil
}

func (s *DbTodoService) UpdateTodoStatus(id string) error {
	var currentDone bool
	err := s.db.conn.QueryRow(context.Background(), "SELECT done FROM todo WHERE id = $1", id).Scan(&currentDone)
	if err != nil {
		return fmt.Errorf("không tìm thấy todo với id %s: %v", id, err)
	}

	// Đảo ngược trạng thái
	newDone := !currentDone
	var doneAt interface{}

	if newDone {
		// Nếu mới là true, cập nhật done_at với thời gian hiện tại
		doneAt = time.Now()
	} else {
		// Nếu mới là false, đặt done_at thành NULL
		doneAt = nil
	}

	// Cập nhật cơ sở dữ liệu
	_, err = s.db.conn.Exec(context.Background(), "UPDATE todo SET done = $1, done_at = $2 WHERE id = $3", newDone, doneAt, id)
	if err != nil {
		return fmt.Errorf("cập nhật trạng thái todo thất bại: %v", err)
	}

	return nil
}
