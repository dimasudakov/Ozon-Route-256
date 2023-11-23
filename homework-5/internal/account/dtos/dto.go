package dtos

import (
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"time"
)

type BankAccountDto struct {
	ID            uuid.UUID                      `json:"id"`
	HolderName    string                         `json:"holder_name"`
	Balance       int                            `json:"balance"`
	OpeningDate   time.Time                      `json:"opening_date"`
	BankName      string                         `json:"bank_name"`
	Subscriptions []subscription.SubscriptionDto `json:"subscriptions"`
}

func (a *BankAccountDto) MapFromModel(model model.BankAccount) *BankAccountDto {

	a.ID = model.ID
	a.HolderName = model.HolderName
	a.Balance = model.Balance
	a.OpeningDate = model.OpeningDate
	a.BankName = model.BankName
	a.Subscriptions = subscription.MapFromModels(model.Subscriptions)

	return a
}

func (a *BankAccountDto) MapToModel() model.BankAccount {

	return model.BankAccount{
		ID:            a.ID,
		HolderName:    a.HolderName,
		Balance:       a.Balance,
		OpeningDate:   a.OpeningDate,
		BankName:      a.BankName,
		Subscriptions: subscription.MapToModels(a.Subscriptions),
	}

}
