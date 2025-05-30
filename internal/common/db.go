package common

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
	return db, nil
}
