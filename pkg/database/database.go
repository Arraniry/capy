package database

import (
	"database/sql"

	_ "github.com/lib/pq" // Import driver PostgreSQL
)

type DB struct {
	*sql.DB
}

func Connect() (*DB, error) {
	db, err := sql.Open("postgres", "user=username dbname=mydb sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
