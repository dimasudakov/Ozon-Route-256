package fixtures

import (
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/bank_accounts"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/subscriptions"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type BankAccountDtoBuilder struct {
	instance *bank_accounts.BankAccountDto
}

func NewBankAccountDtoBuilder() *BankAccountDtoBuilder {
	return &BankAccountDtoBuilder{instance: &bank_accounts.BankAccountDto{}}
}

func (b *BankAccountDtoBuilder) ID(val string) *BankAccountDtoBuilder {
	b.instance.Id = &bank_accounts.UUID{Value: val}
	return b
}

func (b *BankAccountDtoBuilder) HolderName(val string) *BankAccountDtoBuilder {
	b.instance.HolderName = val
	return b
}

func (b *BankAccountDtoBuilder) Balance(val int) *BankAccountDtoBuilder {
	b.instance.Balance = int32(val)
	return b
}

func (b *BankAccountDtoBuilder) OpeningDate(val time.Time) *BankAccountDtoBuilder {
	b.instance.OpeningDate = &bank_accounts.Timestamp{Value: timestamppb.New(val)}
	return b
}

func (b *BankAccountDtoBuilder) BankName(val string) *BankAccountDtoBuilder {
	b.instance.BankName = val
	return b
}

func (b *BankAccountDtoBuilder) Subscriptions(val []*subscriptions.SubscriptionDto) *BankAccountDtoBuilder {
	b.instance.Subscriptions = val
	return b
}

func (b *BankAccountDtoBuilder) Build() *bank_accounts.BankAccountDto {
	return b.instance
}

func (b *BankAccountDtoBuilder) Valid() *BankAccountDtoBuilder {
	//id, _ := uuid.Parse("dc8a075c-6c39-4761-9740-df64bf3b7678")
	return NewBankAccountDtoBuilder().
		ID("dc8a075c-6c39-4761-9740-df64bf3b7678").
		HolderName("Dima Sudakov").
		Balance(1000).
		OpeningDate(time.Time{}).
		BankName("Sberbank").
		Subscriptions(make([]*subscriptions.SubscriptionDto, 0))
}

func (b *BankAccountDtoBuilder) Invalid() *BankAccountDtoBuilder {
	//id, _ := uuid.Parse("dc8a075c-6c39-4761-9740-df64bf3b7678")
	return NewBankAccountDtoBuilder().
		ID("dc8a075c-6c39-4761-9740-df64bf3b7678").
		HolderName("Dima Sudakov 2003").
		Balance(-1000).
		OpeningDate(time.Time{}).
		Subscriptions(make([]*subscriptions.SubscriptionDto, 0))
}
