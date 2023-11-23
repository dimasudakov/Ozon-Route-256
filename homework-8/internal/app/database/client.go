package database

import (
	"database/sql"
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app"
)

func InitDB(config *app.Config) (Database, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Database.Username,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
	)
	dbInstance, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	if err := dbInstance.Ping(); err != nil {
		return nil, err
	}

	dbInstance.SetMaxOpenConns(10)

	return &SQLDatabase{
		db: dbInstance,
	}, nil
}

func InitDBWithPool(DB *sql.DB) (Database, error) {
	return &SQLDatabase{
		db: DB,
	}, nil
}
