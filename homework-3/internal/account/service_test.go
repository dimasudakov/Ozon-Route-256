package account

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"testing"
	"time"
)

type MockRepository struct{}

func (m *MockRepository) CreateBankAccount(account BankAccount) (*BankAccount, error) {
	return &BankAccount{
		ID:            1,
		HolderName:    account.HolderName,
		Balance:       account.Balance,
		BankName:      account.BankName,
		OpeningDate:   account.OpeningDate,
		Subscriptions: []subscription.Subscription{},
	}, nil
}

func (m *MockRepository) GetBankAccountByID(id int) (*BankAccount, error) {
	if id < 100 {
		return &BankAccount{
			ID:            id,
			HolderName:    "Dima Sudakov",
			Balance:       1000,
			OpeningDate:   time.Now(),
			BankName:      "Tinkoff",
			Subscriptions: []subscription.Subscription{},
		}, nil
	} else {
		return nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", id))
	}
}

func (m *MockRepository) UpdateBankAccount(id int, account BankAccount) (*BankAccount, error) {
	if id <= 100 {
		return &BankAccount{
			ID:         1,
			HolderName: "Updated Dima",
			Balance:    2000,
			BankName:   "Sberbank",
		}, nil
	} else {
		return nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", id))
	}
}

func (m *MockRepository) DeleteBankAccount(id int) (*BankAccount, error) {
	if id <= 100 {
		return &BankAccount{
			ID:         1,
			HolderName: "Dima Sudakov",
			Balance:    1000,
			BankName:   "Tinkoff",
		}, nil
	} else {
		return nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", id))
	}
}

func TestBankAccountService_CreateBankAccount(t *testing.T) {
	service := NewBankAccountService(&MockRepository{})
	var emptyAccount BankAccount

	tests := []struct {
		name           string
		requestPayload BankAccount
		expectedResult BankAccount
		expectedError  error
	}{
		{
			name: "Valid Create Request",
			requestPayload: BankAccount{
				HolderName: "Dima Sudakov",
				Balance:    1000,
				BankName:   "Tinkoff",
			},
			expectedResult: BankAccount{
				HolderName:  "Dima Sudakov",
				Balance:     1000,
				OpeningDate: time.Now(),
				BankName:    "Tinkoff",
			},
			expectedError: nil,
		},
		{
			name: "Invalid Request",
			requestPayload: BankAccount{
				HolderName: "Dima 228",
				Balance:    666,
				BankName:   "Sberbank",
			},
			expectedResult: emptyAccount,
			expectedError:  apperr.NewBadRequestError("HolderName: must be in a valid format."),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.CreateBankAccount(tc.requestPayload)

			assert.Equal(t, tc.expectedError, err, "error should match")
			if err == nil {
				assert.NotEqual(t, result.ID, 0)
				assert.Equal(t, result.HolderName, tc.expectedResult.HolderName)
				assert.Equal(t, result.Balance, tc.expectedResult.Balance)
				assert.Equal(t, result.BankName, tc.expectedResult.BankName)
				assert.Equal(t, result.OpeningDate.Round(time.Second), tc.expectedResult.OpeningDate.Round(time.Second))
			} else {
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}

func TestBankAccountService_GetBankAccountById(t *testing.T) {
	service := NewBankAccountService(&MockRepository{})

	var emptyAccount BankAccount

	tests := []struct {
		name           string
		requestID      int
		expectedResult BankAccount
		expectedError  error
	}{
		{
			name:      "Valid Get Request",
			requestID: 1,
			expectedResult: BankAccount{
				ID:          1,
				HolderName:  "Dima Sudakov",
				Balance:     1000,
				OpeningDate: time.Now(),
				BankName:    "Tinkoff",
			},
			expectedError: nil,
		},
		{
			name:           "Not Found Request",
			requestID:      10329032,
			expectedResult: emptyAccount,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.GetBankAccountById(tc.requestID)

			assert.Equal(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, result.ID, tc.expectedResult.ID)
				assert.Equal(t, result.HolderName, tc.expectedResult.HolderName)
				assert.Equal(t, result.Balance, tc.expectedResult.Balance)
				assert.Equal(t, result.OpeningDate.Round(time.Second), tc.expectedResult.OpeningDate.Round(time.Second))
				assert.Equal(t, result.BankName, tc.expectedResult.BankName)
			} else {
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}

func TestBankAccountService_UpdateBankAccount(t *testing.T) {
	service := NewBankAccountService(&MockRepository{})

	var emptyAccount BankAccount

	tests := []struct {
		name           string
		requestID      int
		requestPayload BankAccount
		expectedResult BankAccount
		expectedError  error
	}{
		{
			name:      "Valid Update Request",
			requestID: 1,
			requestPayload: BankAccount{
				HolderName: "Updated Dima",
				Balance:    2000,
				BankName:   "Sberbank",
			},
			expectedResult: BankAccount{
				ID:         1,
				HolderName: "Updated Dima",
				Balance:    2000,
				BankName:   "Sberbank",
			},
			expectedError: nil,
		},
		{
			name:      "Invalid Request",
			requestID: 2,
			requestPayload: BankAccount{
				HolderName: "Invalid Name",
				Balance:    -500,
			},
			expectedResult: emptyAccount,
			expectedError:  apperr.NewBadRequestError("Balance: must be no less than 0; BankName: cannot be blank."),
		},
		{
			name:      "Not Found Request",
			requestID: 10329032,
			requestPayload: BankAccount{
				HolderName: "Dima Sudakov",
				Balance:    1000,
				BankName:   "Sberbank",
			},
			expectedResult: emptyAccount,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.UpdateBankAccount(tc.requestID, tc.requestPayload)

			assert.Equal(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, result.ID, tc.expectedResult.ID)
				assert.Equal(t, result.HolderName, tc.expectedResult.HolderName)
				assert.Equal(t, result.Balance, tc.expectedResult.Balance)
				assert.Equal(t, result.BankName, tc.expectedResult.BankName)
			} else {
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}

func TestBankAccountService_DeleteBankAccount(t *testing.T) {
	service := NewBankAccountService(&MockRepository{})

	var emptyAccount *BankAccount

	tests := []struct {
		name           string
		requestID      int
		expectedResult *BankAccount
		expectedError  error
	}{
		{
			name:      "Valid Delete Request",
			requestID: 1,
			expectedResult: &BankAccount{
				ID:         1,
				HolderName: "Dima Sudakov",
				Balance:    1000,
				BankName:   "Tinkoff",
			},
			expectedError: nil,
		},
		{
			name:           "Not Found Request",
			requestID:      10329032,
			expectedResult: emptyAccount,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.DeleteBankAccount(tc.requestID)

			assert.Equal(t, err, tc.expectedError)
			if err == nil {
				assert.Equal(t, result.ID, tc.expectedResult.ID)
				assert.Equal(t, result.HolderName, tc.expectedResult.HolderName)
				assert.Equal(t, result.Balance, tc.expectedResult.Balance)
				assert.Equal(t, result.BankName, tc.expectedResult.BankName)
				assert.WithinDuration(t, result.OpeningDate, tc.expectedResult.OpeningDate, time.Second)
			} else {
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}
