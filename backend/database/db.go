package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	// Используем твой реальный пароль Rid@123SQL#
	connStr := fmt.Sprintf("host=%s user=postgres password=Rid@123SQL# dbname=diploma_platform sslmode=disable", host)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	
	// Проверяем реальное соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}
