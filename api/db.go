package main

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type Db struct {
	conn *pgxpool.Pool
}

func NewDb() (*Db, error) {

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("thiếu biến môi trường database")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối đến cơ sở dữ liệu: %v", err)
	}

	return &Db{conn: conn}, nil
}

func Migrate() error {
	connString := fmt.Sprintf("cockroachdb://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	m, err := migrate.New(
		"file://db/migrations",
		connString,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if err.Error() == "Dirty database version 1. Fix and force version." {
			fmt.Println("Database đang ở trạng thái dirty, đang thực hiện Force version về 1")
			if forceErr := m.Force(1); forceErr != nil {
				return fmt.Errorf("không thể Force version: %v", forceErr)
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("lỗi khi thực hiện lại migration: %v", err)
			}
		} else {
			return fmt.Errorf("lỗi khi thực hiện migration: %v", err)
		}
	}

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("lỗi khi lấy phiên bản migration: %v", err)
	}

	if dirty {
		fmt.Println("Database đang ở trạng thái dirty, cần phải được sửa chữa.")
	}

	fmt.Printf("Migration thành công, phiên bản hiện tại: %d\n", version)
	return nil
}
