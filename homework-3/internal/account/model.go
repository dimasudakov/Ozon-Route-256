package account

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"regexp"
	"time"
)

type BankAccount struct {
	ID            int                         `db:"id"`
	HolderName    string                      `db:"holder_name"`
	Balance       int                         `db:"balance"`
	OpeningDate   time.Time                   `db:"opening_date"`
	BankName      string                      `db:"bank_name"`
	Subscriptions []subscription.Subscription `db:"subscriptions"`
}

func (a BankAccount) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.HolderName,
			validation.Required,
			validation.Length(3, 100),
			validation.Match(regexp.MustCompile("^[a-zA-Z ]+$")),
		),
		validation.Field(&a.Balance,
			validation.Min(0),
		),
		validation.Field(&a.BankName,
			validation.Required,
			validation.Length(3, 20),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]+$")),
		),
	)
}
