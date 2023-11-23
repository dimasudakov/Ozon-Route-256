package fixtures

import (
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/dtos"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"time"
)

type BankAccountDtoBuilder struct {
	instance *dtos.BankAccountDto
}

func NewBankAccountDtoBuilder() *BankAccountDtoBuilder {
	return &BankAccountDtoBuilder{instance: &dtos.BankAccountDto{}}
}

func (b *BankAccountDtoBuilder) ID(val uuid.UUID) *BankAccountDtoBuilder {
	b.instance.ID = val
	return b
}

func (b *BankAccountDtoBuilder) HolderName(val string) *BankAccountDtoBuilder {
	b.instance.HolderName = val
	return b
}

func (b *BankAccountDtoBuilder) Balance(val int) *BankAccountDtoBuilder {
	b.instance.Balance = val
	return b
}

func (b *BankAccountDtoBuilder) OpeningDate(val time.Time) *BankAccountDtoBuilder {
	b.instance.OpeningDate = val
	return b
}

func (b *BankAccountDtoBuilder) BankName(val string) *BankAccountDtoBuilder {
	b.instance.BankName = val
	return b
}

func (b *BankAccountDtoBuilder) Subscriptions(val []subscription.SubscriptionDto) *BankAccountDtoBuilder {
	b.instance.Subscriptions = val
	return b
}

func (b *BankAccountDtoBuilder) Build() *dtos.BankAccountDto {
	return b.instance
}

func (b *BankAccountDtoBuilder) Valid() *BankAccountDtoBuilder {
	id, _ := uuid.Parse("a7115d4e-65af-487f-a3ca-bf7ca9747c4c")
	return NewBankAccountDtoBuilder().
		ID(id).
		HolderName("Dima Sudakov").
		Balance(1000).
		OpeningDate(time.Time{}).
		BankName("Sberbank").
		Subscriptions(make([]subscription.SubscriptionDto, 0))
}

func (b *BankAccountDtoBuilder) Invalid() *BankAccountDtoBuilder {
	id, _ := uuid.Parse("dc8a075c-6c39-4761-9740-df64bf3b7678")
	return NewBankAccountDtoBuilder().
		ID(id).
		HolderName("Dima Sudakov 2003").
		Balance(-1000).
		OpeningDate(time.Time{}).
		Subscriptions(make([]subscription.SubscriptionDto, 0))
}
