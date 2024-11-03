package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type Db struct {
	conn *pgx.Conn
}

func NewDb() (*Db, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("không thể tải file .env: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối đến cơ sở dữ liệu: %v", err)
	}

	return &Db{conn: conn}, nil
}

func (db *Db) Close() {
	db.conn.Close(context.Background())
}