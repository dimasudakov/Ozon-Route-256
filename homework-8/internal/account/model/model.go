package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/bank_accounts"
	"google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"time"
)

type BankAccount struct {
	ID            uuid.UUID                   `db:"id"`
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

func MapFromDto(dto *bank_accounts.BankAccountDto) (*BankAccount, error) {
	id, err := uuid.Parse(dto.GetId().GetValue())
	if err != nil {
		return nil, errors.New("invalid id")
	}
	subs, err := subscription.MapFromDtoList(dto.Subscriptions)
	if err != nil {
		return nil, err
	}
	return &BankAccount{
		ID:            id,
		HolderName:    dto.GetHolderName(),
		Balance:       int(dto.GetBalance()),
		OpeningDate:   dto.GetOpeningDate().GetValue().AsTime(),
		BankName:      dto.GetBankName(),
		Subscriptions: subs,
	}, nil
}

func (a BankAccount) MapToDto() *bank_accounts.BankAccountDto {
	return &bank_accounts.BankAccountDto{
		Id:            &bank_accounts.UUID{Value: a.ID.String()},
		HolderName:    a.HolderName,
		Balance:       int32(a.Balance),
		OpeningDate:   &bank_accounts.Timestamp{Value: timestamppb.New(a.OpeningDate)},
		BankName:      a.BankName,
		Subscriptions: subscription.MapToDtoList(a.Subscriptions),
	}
}
