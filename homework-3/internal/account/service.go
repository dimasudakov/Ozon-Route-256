package account

import (
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"time"
)

type Repository interface {
	CreateBankAccount(account BankAccount) (*BankAccount, error)
	GetBankAccountByID(id int) (*BankAccount, error)
	UpdateBankAccount(id int, account BankAccount) (*BankAccount, error)
	DeleteBankAccount(id int) (*BankAccount, error)
}

type BankAccountService struct {
	repository Repository
}

func NewBankAccountService(repository Repository) *BankAccountService {
	return &BankAccountService{repository: repository}
}

func (b *BankAccountService) CreateBankAccount(account BankAccount) (*BankAccount, error) {
	if err := account.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}

	account.OpeningDate = time.Now()

	bankAccount, err := b.repository.CreateBankAccount(account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) GetBankAccountById(id int) (*BankAccount, error) {
	bankAccount, err := b.repository.GetBankAccountByID(id)
	if err != nil {
		return nil, err
	}
	return bankAccount, nil
}

func (b *BankAccountService) UpdateBankAccount(id int, account BankAccount) (*BankAccount, error) {
	if err := account.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}

	bankAccount, err := b.repository.UpdateBankAccount(id, account)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}

func (b *BankAccountService) DeleteBankAccount(id int) (*BankAccount, error) {
	bankAccount, err := b.repository.DeleteBankAccount(id)
	if err != nil {
		return nil, err
	}

	return bankAccount, nil
}
