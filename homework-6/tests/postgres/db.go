//go:build integration

package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/database"
	"log"
	"sync"
	"testing"
)

type TDB struct {
	DB database.Database
	sync.Mutex
}

func NewFromEnv(sqlDB *sql.DB) *TDB {
	db, err := database.InitDBWithPool(sqlDB)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, args ...interface{}) {
	t.Helper()
	d.Lock()
	d.Truncate()
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate()
}

func (d *TDB) Truncate() {
	rows, err := d.DB.QueryRows("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		tableNames = append(tableNames, tableName)
	}

	if len(tableNames) == 0 {
		panic("run migrations")
	}
}

func (d *TDB) ExistsById(id uuid.UUID, tableName string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM " + tableName + " WHERE id = $1)"
	var exists bool
	err := d.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *TDB) DeleteBankAccounts(accountIDs []uuid.UUID) error {
	query := `DELETE FROM bank_account WHERE id = $1`

	for _, id := range accountIDs {
		_, err := d.DB.Execute(query, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *TDB) InsertSampleData(accounts []model.BankAccount) error {
	tx, err := d.DB.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, account := range accounts {
		_, err := tx.Exec("INSERT INTO bank_account (id, holder_name, balance, opening_date, bank_name) VALUES ($1, $2, $3, $4, $5)",
			account.ID, account.HolderName, account.Balance, account.OpeningDate, account.BankName)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
