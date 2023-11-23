//go:generate mockgen -source=./service.go -destination=./mocks/repository.go -package=mock_account

package account

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"time"
)

type Repository interface {
	CreateBankAccount(account *model.BankAccount) (*model.BankAccount, error)
	GetBankAccountByID(id uuid.UUID) (*model.BankAccount, error)
	UpdateBankAccount(id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error)
	DeleteBankAccount(id uuid.UUID) (*model.BankAccount, error)
}

type BankAccountService struct {
	repository Repository
}

func NewBankAccountService(repository Repository) *BankAccountService {
	return &BankAccountService{repository: repository}
}

func (b *BankAccountService) CreateBankAccount(account *model.BankAccount) (*model.BankAccount, error) {
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

	bankAccount, err := b.repository.CreateBankAccount(account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) GetBankAccountById(id uuid.UUID) (*model.BankAccount, error) {
	bankAccount, err := b.repository.GetBankAccountByID(id)
	if err != nil {
		return nil, err
	}
	return bankAccount, nil
}

func (b *BankAccountService) UpdateBankAccount(id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error) {
	if err := account.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}
	str := id.String()
	fmt.Println(str)
	bankAccount, err := b.repository.UpdateBankAccount(id, account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) DeleteBankAccount(id uuid.UUID) (*model.BankAccount, error) {
	bankAccount, err := b.repository.DeleteBankAccount(id)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}
