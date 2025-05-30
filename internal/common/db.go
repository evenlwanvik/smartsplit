package common

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrUnableToCreateDB = errors.New("unable to create db connection")
)

func NewDB(dsn string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, ErrUnableToCreateDB
	}
	return db, nil
}
