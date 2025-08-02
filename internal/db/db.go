package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"-"`
	Database string `json:"database,omitempty"`
	SSLMode  string `json:"ssl_mode,omitempty"`
}

var (
	ErrUnableToCreateDB = errors.New("unable to create db connection")
	ErrDBUnreachable    = errors.New("database is unreachable")
)

func NewDB(cred PostgresDB) (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cred.Host, cred.Port, cred.User, cred.Password, cred.Database, cred.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, ErrUnableToCreateDB
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBUnreachable, err)
	}
	return db, nil
}
