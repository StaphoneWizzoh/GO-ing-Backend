package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(dbURL string) (*sql.DB, error){
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}