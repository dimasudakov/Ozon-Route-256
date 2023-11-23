package app

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

type Database struct {
	DB *sql.DB
}

func InitDB(config *Config) (*Database, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Username,
		config.Database.Password,
		config.Database.Name,
	)
	dbInstance, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	if err := dbInstance.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	dbInstance.SetMaxOpenConns(10)

	return &Database{
		DB: dbInstance,
	}, nil
}

func (d *Database) UpMigrations() {
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(d.DB, "migrations"); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
