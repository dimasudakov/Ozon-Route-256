package subscription

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"subscription_name"`
	Price     int       `db:"price"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	AccountID uuid.UUID `db:"account_id"`
}

func (s Subscription) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Price, validation.Required, validation.Min(0)),
		validation.Field(&s.EndDate, validation.Min(time.Now())),
	)
}
