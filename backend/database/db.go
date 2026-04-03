package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := "user=postgres password=Rid@123SQL# dbname=diploma_platform sslmode=disable"
	return sql.Open("postgres", connStr)
}
