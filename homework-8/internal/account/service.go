//go:generate mockgen -source=./service.go -destination=./mocks/service.go -package=mock_account

package account

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"time"
)

type Service interface {
	CreateBankAccount(ctx context.Context, account *model.BankAccount) (*model.BankAccount, error)
	GetBankAccountById(ctx context.Context, id uuid.UUID) (*model.BankAccount, error)
	UpdateBankAccount(ctx context.Context, id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error)
	DeleteBankAccount(ctx context.Context, id uuid.UUID) (*model.BankAccount, error)
}

type BankAccountService struct {
	repository Repository
}

func NewBankAccountService(repository Repository) *BankAccountService {
	return &BankAccountService{repository: repository}
}

func (b *BankAccountService) CreateBankAccount(ctx context.Context, account *model.BankAccount) (*model.BankAccount, error) {
	if err := account.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}

	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}

	for i := range account.Subscriptions {
		account.Subscriptions[i].AccountID = account.ID
		if account.Subscriptions[i].ID == uuid.Nil {
			account.Subscriptions[i].ID = uuid.New()
		}
	}

	account.OpeningDate = time.Now()

	bankAccount, err := b.repository.CreateBankAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) GetBankAccountById(ctx context.Context, id uuid.UUID) (*model.BankAccount, error) {
	bankAccount, err := b.repository.GetBankAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return bankAccount, nil
}

func (b *BankAccountService) UpdateBankAccount(ctx context.Context, id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error) {
	if err := account.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}
	str := id.String()
	fmt.Println(str)
	bankAccount, err := b.repository.UpdateBankAccount(ctx, id, account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) DeleteBankAccount(ctx context.Context, id uuid.UUID) (*model.BankAccount, error) {
	bankAccount, err := b.repository.DeleteBankAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}
