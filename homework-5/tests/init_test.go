//go:build integration

package tests

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/postgres"
	"os"
	"path/filepath"
	"sync"
)

type DBConfig struct {
	TestDatabase struct {
		Name     string `mapstructure:"name"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"test-database"`
}

var (
	db   *postgres.TDB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		config, err := initConfig()
		if err != nil {
			panic(fmt.Sprintf("Can't read database config: %s", err))
		}
		connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.TestDatabase.Username,
			config.TestDatabase.Password,
			config.TestDatabase.Name,
			config.TestDatabase.Port,
		)
		dbInstance, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}

		if err := dbInstance.Ping(); err != nil {
			return
		}

		dbInstance.SetMaxOpenConns(10)

		db = postgres.NewFromEnv(dbInstance)
	})
}

func initConfig() (*DBConfig, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения текущей директории")
	}

	configPath := filepath.Join(projectRoot, "configs", "config.yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config DBConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
