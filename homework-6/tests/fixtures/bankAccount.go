package fixtures

import (
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"time"
)

type BankAccountBuilder struct {
	instance *model.BankAccount
}

func NewBankAccountBuilder() *BankAccountBuilder {
	return &BankAccountBuilder{instance: &model.BankAccount{}}
}

func (b *BankAccountBuilder) ID(val uuid.UUID) *BankAccountBuilder {
	b.instance.ID = val
	return b
}

func (b *BankAccountBuilder) HolderName(val string) *BankAccountBuilder {
	b.instance.HolderName = val
	return b
}

func (b *BankAccountBuilder) Balance(val int) *BankAccountBuilder {
	b.instance.Balance = val
	return b
}

func (b *BankAccountBuilder) OpeningDate(val time.Time) *BankAccountBuilder {
	b.instance.OpeningDate = val
	return b
}

func (b *BankAccountBuilder) BankName(val string) *BankAccountBuilder {
	b.instance.BankName = val
	return b
}

func (b *BankAccountBuilder) Subscriptions(val []subscription.Subscription) *BankAccountBuilder {
	b.instance.Subscriptions = val
	return b
}

func (b *BankAccountBuilder) Build() *model.BankAccount {
	return b.instance
}

func (b *BankAccountBuilder) Valid() *BankAccountBuilder {
	id, _ := uuid.Parse("a7115d4e-65af-487f-a3ca-bf7ca9747c4c")
	return NewBankAccountBuilder().
		ID(id).
		HolderName("Dima Sudakov").
		Balance(1000).
		OpeningDate(time.Time{}).
		BankName("Sberbank").
		Subscriptions(make([]subscription.Subscription, 0))
}

func (b *BankAccountBuilder) Invalid() *BankAccountBuilder {
	id, _ := uuid.Parse("a7115d4e-65af-487f-a3ca-bf7ca9747c4c")
	return NewBankAccountBuilder().
		ID(id).
		HolderName("Dima Sudakov 2003").
		Balance(1000).
		OpeningDate(time.Time{}).
		BankName("Sberbank")
}
