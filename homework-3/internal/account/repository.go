package account

import (
	"database/sql"
	"errors"
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"

	_ "github.com/lib/pq"
)

type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{
		db: db,
	}
}

func (r *BankAccountRepository) CreateBankAccount(account BankAccount) (*BankAccount, error) {
	tx, err := r.db.Begin()
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

	query := `
		INSERT INTO bank_account (holder_name, balance, opening_date, bank_name) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`
	var id int
	err = tx.QueryRow(query, account.HolderName, account.Balance, account.OpeningDate, account.BankName).Scan(&id)
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	for i, sub := range account.Subscriptions {
		query = `
			INSERT INTO subscription (account_id, subscription_name, price, start_date) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id`
		err = tx.QueryRow(query, id, sub.Name, sub.Price, sub.StartDate).Scan(&account.Subscriptions[i].ID)
		if err != nil {
			return nil, apperr.NewInternalServerError("Internal server error")
		}
		account.Subscriptions[i].AccountID = id
	}

	createdAccount := BankAccount{
		ID:            id,
		HolderName:    account.HolderName,
		Balance:       account.Balance,
		OpeningDate:   account.OpeningDate,
		BankName:      account.BankName,
		Subscriptions: account.Subscriptions,
	}

	return &createdAccount, nil
}

func (r *BankAccountRepository) GetBankAccountByID(id int) (*BankAccount, error) {
	query := "SELECT * FROM bank_account WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var bankAccount BankAccount

	if err := row.Scan(&bankAccount.ID, &bankAccount.HolderName, &bankAccount.Balance, &bankAccount.OpeningDate, &bankAccount.BankName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", id))
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

func (r *BankAccountRepository) UpdateBankAccount(id int, account BankAccount) (*BankAccount, error) {
	query := "UPDATE bank_account SET id = $1, holder_name = $2, balance = $3, bank_name = $4 WHERE id = $5 RETURNING id, holder_name, balance, opening_date, bank_name"

	var updatedAccount BankAccount
	err := r.db.QueryRow(query, account.ID, account.HolderName, account.Balance, account.BankName, id).
		Scan(&updatedAccount.ID, &updatedAccount.HolderName, &updatedAccount.Balance, &updatedAccount.OpeningDate, &updatedAccount.BankName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %d not found", account.ID))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	return &updatedAccount, nil
}

func (r *BankAccountRepository) DeleteBankAccount(id int) (*BankAccount, error) {
	tx, err := r.db.Begin()
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

	var deletedAccount BankAccount
	err = tx.QueryRow(query, id).
		Scan(&deletedAccount.ID, &deletedAccount.HolderName, &deletedAccount.Balance, &deletedAccount.OpeningDate, &deletedAccount.BankName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %d does not exist", id))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	deletedAccount.Subscriptions = deletedSubscriptions

	return &deletedAccount, nil
}

func (r *BankAccountRepository) getSubscripctionsByBankAccountID(id int) ([]subscription.Subscription, error) {
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

	rows, err := r.db.Query(query, id)
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
