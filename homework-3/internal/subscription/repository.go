package subscription

import (
	"database/sql"
	"errors"
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r SubscriptionRepository) CreateSubscription(subscription Subscription) (*Subscription, error) {
	query := `
		INSERT INTO subscription (subscription_name, price, start_date, end_date, account_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	var id int
	err := r.db.QueryRow(
		query,
		subscription.Name,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
		subscription.AccountID,
	).Scan(&id)
	if err != nil {
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	createdSubscription := Subscription{
		ID:        id,
		Name:      subscription.Name,
		Price:     subscription.Price,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		AccountID: subscription.AccountID,
	}

	return &createdSubscription, nil
}

func (r SubscriptionRepository) GetSubscriptionByID(id int) (*Subscription, error) {
	query := `
		SELECT id, subscription_name, price, start_date, end_date, account_id 
		FROM subscription 
		WHERE id = $1`

	var subscription Subscription

	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.Name,
		&subscription.Price,
		&subscription.StartDate,
		&sql.NullTime{Time: subscription.EndDate, Valid: true},
		&subscription.AccountID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFoundError(fmt.Sprintf("Subscription with ID: %d not found", id))
		}
		return nil, apperr.NewInternalServerError("Internal server error")
	}

	return &subscription, nil
}
