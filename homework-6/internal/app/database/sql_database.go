//go:generate mockgen -source ./sql_database.go -destination=./mocks/database.go -package=mock_database

package database

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

type Database interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRows(query string, args ...interface{}) (*sql.Rows, error)
	Execute(query string, args ...interface{}) (sql.Result, error)
	BeginTx() (*sql.Tx, error)
	Close() error
	UpMigrations() error
}

type SQLDatabase struct {
	db *sql.DB
}

func (s *SQLDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *SQLDatabase) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *SQLDatabase) Execute(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SQLDatabase) BeginTx() (*sql.Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *SQLDatabase) UpMigrations() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db, "db/migrations"); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *SQLDatabase) Close() error {
	return s.db.Close()
}
