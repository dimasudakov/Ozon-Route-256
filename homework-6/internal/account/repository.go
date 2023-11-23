//go:generate mockgen -source=./repository.go -destination=./mocks/repository.go -package=mock_account

package account

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/database"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"

	_ "github.com/lib/pq"
)

type BankAccountRepository struct {
	db database.Database
}

func NewBankAccountRepository(db database.Database) *BankAccountRepository {
	return &BankAccountRepository{
		db: db,
	}
}

func (r *BankAccountRepository) CreateBankAccount(account *model.BankAccount) (*model.BankAccount, error) {
	tx, err := r.db.BeginTx()
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO bank_account (id, holder_name, balance, opening_date, bank_name) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_, err = tx.Exec(query, account.ID, account.HolderName, account.Balance, account.OpeningDate, account.BankName)
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	for _, sub := range account.Subscriptions {
		query = `
			INSERT INTO subscription (id, account_id, subscription_name, price, start_date) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id`
		_, err = tx.Exec(query, sub.ID, account.ID, sub.Name, sub.Price, sub.StartDate)
		if err != nil {
			return nil, apperr.NewInternalServerError(fmt.Sprintf("Internal server error: %s", err))
		}
	}

	return account, nil
}

func (r *BankAccountRepository) GetBankAccountByID(id uuid.UUID) (*model.BankAccount, error) {
	query := "SELECT * FROM bank_account WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var bankAccount model.BankAccount

	if err := row.Scan(&bankAccount.ID, &bankAccount.HolderName, &bankAccount.Balance, &bankAccount.OpeningDate, &bankAccount.BankName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %s not found", id.String()))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	subscriptions, err := r.getSubscripctionsByBankAccountID(id)
	if err != nil {
		return nil, err
	}
	bankAccount.Subscriptions = subscriptions

	return &bankAccount, nil
}

func (r *BankAccountRepository) UpdateBankAccount(id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error) {
	query := "UPDATE bank_account SET id = $1, holder_name = $2, balance = $3, bank_name = $4 WHERE id = $5 RETURNING id, holder_name, balance, opening_date, bank_name"

	var updatedAccount model.BankAccount
	err := r.db.QueryRow(query, account.ID, account.HolderName, account.Balance, account.BankName, id).
		Scan(&updatedAccount.ID, &updatedAccount.HolderName, &updatedAccount.Balance, &updatedAccount.OpeningDate, &updatedAccount.BankName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %s not found", id))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	return &updatedAccount, nil
}

func (r *BankAccountRepository) DeleteBankAccount(id uuid.UUID) (*model.BankAccount, error) {
	tx, err := r.db.BeginTx()
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	rows, err := tx.Query("DELETE FROM subscription WHERE account_id = $1 RETURNING id, subscription_name, price, start_date, end_date, account_id", id)
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}
	defer rows.Close()

	var deletedSubscriptions []subscription.Subscription
	for rows.Next() {
		var sub subscription.Subscription
		err := rows.Scan(&sub.ID, &sub.Name, &sub.Price, &sub.StartDate, &sql.NullTime{Time: sub.EndDate, Valid: true}, &sub.AccountID)
		if err != nil {
			return nil, apperr.NewInternalServerError("Internal server error")
		}
		deletedSubscriptions = append(deletedSubscriptions, sub)
	}

	query := `
		DELETE FROM bank_account 
		WHERE id = $1 
		RETURNING id, holder_name, balance, opening_date, bank_name`

	var deletedAccount model.BankAccount
	err = tx.QueryRow(query, id).
		Scan(&deletedAccount.ID, &deletedAccount.HolderName, &deletedAccount.Balance, &deletedAccount.OpeningDate, &deletedAccount.BankName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %s does not exist", id.String()))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	deletedAccount.Subscriptions = deletedSubscriptions

	return &deletedAccount, nil
}

func (r *BankAccountRepository) getSubscripctionsByBankAccountID(id uuid.UUID) ([]subscription.Subscription, error) {
	query := `
        SELECT 
            s.id,
            s.subscription_name,
            s.price,
            s.start_date,
            s.end_date,
            s.account_id
        FROM subscription s
        WHERE s.account_id = $1`

	rows, err := r.db.QueryRows(query, id)
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}
	defer rows.Close()

	var subscriptions []subscription.Subscription
	for rows.Next() {
		var sub subscription.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.Name,
			&sub.Price,
			&sub.StartDate,
			&sql.NullTime{Time: sub.EndDate, Valid: true},
			&sub.AccountID,
		)
		if err != nil {
			return nil, apperr.NewInternalServerError("Internal server error")
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}
