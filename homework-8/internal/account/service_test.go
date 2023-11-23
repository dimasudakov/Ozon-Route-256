package account

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_account "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/mocks"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/fixtures"
	"testing"
)

type bankAccountServiceFixture struct {
	ctrl     *gomock.Controller
	service  Service
	mockRepo *mock_account.MockRepository
}

func NewBankAccountServiceFixture(t *testing.T) bankAccountServiceFixture {
	ctrl := gomock.NewController(t)
	mockRepo := mock_account.NewMockRepository(ctrl)
	service := NewBankAccountService(mockRepo)
	return bankAccountServiceFixture{
		ctrl:     ctrl,
		service:  service,
		mockRepo: mockRepo,
	}
}

func (f *bankAccountServiceFixture) tearDown() {
	f.ctrl.Finish()
}

func TestBankAccountService_CreateBankAccount(t *testing.T) {
	t.Parallel()
	var (
		ctx                = context.Background()
		emptyBankAccount   model.BankAccount
		bankAccount        = fixtures.NewBankAccountBuilder().Valid().Build()
		createdBankAccount = fixtures.NewBankAccountBuilder().Valid().Build()
		invalidBankAccount = fixtures.NewBankAccountBuilder().Invalid().Build()
	)

	tests := []struct {
		name           string
		requestPayload model.BankAccount
		expectedResult *model.BankAccount
		mockRepo       func(repository *mock_account.MockRepository)
		expectedError  error
	}{
		{
			name:           "Valid Create Request",
			requestPayload: *bankAccount,
			expectedResult: bankAccount,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any()).Return(createdBankAccount, nil)
			},
			expectedError: nil,
		},
		{
			name:           "Invalid Create Request",
			requestPayload: *invalidBankAccount,
			expectedResult: &emptyBankAccount,
			mockRepo:       nil,
			expectedError:  apperr.NewBadRequestError("HolderName: must be in a valid format."),
		},
		{
			name:           "Fail Repository",
			requestPayload: *bankAccount,
			expectedResult: nil,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any()).Return(nil, apperr.NewInternalServerError("Internal server error"))
			},
			expectedError: apperr.NewInternalServerError("Internal server error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountServiceFixture(t)
			if tc.mockRepo != nil {
				tc.mockRepo(fixture.mockRepo)
			}

			result, err := fixture.service.CreateBankAccount(ctx, &tc.requestPayload)

			if tc.expectedError == nil {
				require.NoError(t, err)
				cmp.Equal(tc.expectedResult, result, cmpopts.IgnoreFields(model.BankAccount{}, "OpeningDate"))
			} else {
				require.EqualError(t, tc.expectedError, err.Error())
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}

func TestBankAccountService_GetBankAccountById(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.Background()
		emptyAccount    *model.BankAccount
		expectedAccount = fixtures.NewBankAccountBuilder().Valid().Build()
		notFoundId, _   = uuid.Parse("331684ab-5af8-439a-8a4f-62a571013283")
	)

	tests := []struct {
		name           string
		requestID      uuid.UUID
		mockRepo       func(repository *mock_account.MockRepository)
		expectedResult *model.BankAccount
		expectedError  error
	}{
		{
			name:      "Valid QueryRow Request",
			requestID: expectedAccount.ID,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().GetBankAccountByID(ctx, expectedAccount.ID).Return(expectedAccount, nil)
			},
			expectedResult: expectedAccount,
			expectedError:  nil,
		},
		{
			name:      "Not Found Request",
			requestID: notFoundId,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().GetBankAccountByID(ctx, notFoundId).
					Return(nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)))
			},
			expectedResult: emptyAccount,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountServiceFixture(t)
			if tc.mockRepo != nil {
				tc.mockRepo(fixture.mockRepo)
			}

			result, err := fixture.service.GetBankAccountById(ctx, tc.requestID)

			assert.Equal(t, err, tc.expectedError)
			if tc.expectedError == nil {
				require.NoError(t, err)
				cmp.Equal(tc.expectedResult, result, cmpopts.IgnoreFields(model.BankAccount{}, "OpeningDate"))
			} else {
				assert.EqualError(t, tc.expectedError, err.Error())
				assert.Nil(t, result, "result should be nil when there is an error")
			}
		})
	}
}

func TestBankAccountService_UpdateBankAccount(t *testing.T) {
	t.Parallel()
	var (
		ctx                = context.Background()
		bankAccount        = fixtures.NewBankAccountBuilder().Valid().Balance(2000).Build()
		invalidBankAccount = fixtures.NewBankAccountBuilder().Invalid().Balance(-500).BankName("").Build()
		notFoundId, _      = uuid.Parse("331684ab-5af8-439a-8a4f-62a571013283")
	)

	tests := []struct {
		name           string
		requestID      uuid.UUID
		requestPayload model.BankAccount
		mockRepo       func(repository *mock_account.MockRepository)
		expectedResult *model.BankAccount
		expectedError  error
	}{
		{
			name:           "Valid Update Request",
			requestID:      bankAccount.ID,
			requestPayload: *bankAccount,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().UpdateBankAccount(ctx, bankAccount.ID, bankAccount).Return(bankAccount, nil)
			},
			expectedResult: bankAccount,
			expectedError:  nil,
		},
		{
			name:           "Invalid Request",
			requestID:      bankAccount.ID,
			requestPayload: *invalidBankAccount,
			mockRepo:       nil,
			expectedResult: nil,
			expectedError:  apperr.NewBadRequestError("Balance: must be no less than 0; BankName: cannot be blank; HolderName: must be in a valid format."),
		},
		{
			name:           "Not Found Request",
			requestID:      notFoundId,
			requestPayload: *bankAccount,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().UpdateBankAccount(ctx, notFoundId, bankAccount).Return(
					nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)))
			},
			expectedResult: nil,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountServiceFixture(t)
			if tc.mockRepo != nil {
				tc.mockRepo(fixture.mockRepo)
			}

			result, err := fixture.service.UpdateBankAccount(ctx, tc.requestID, &tc.requestPayload)

			if tc.expectedError == nil {
				require.NoError(t, err)
				cmp.Equal(tc.expectedResult, result, cmpopts.IgnoreFields(model.BankAccount{}, "OpeningDate"))
			} else {
				assert.EqualError(t, tc.expectedError, err.Error())
				assert.Nil(t, result)
			}
		})
	}
}

func TestBankAccountService_DeleteBankAccount(t *testing.T) {
	t.Parallel()
	var (
		ctx           = context.Background()
		bankAccount   = fixtures.NewBankAccountBuilder().Valid().Build()
		notFoundId, _ = uuid.Parse("331684ab-5af8-439a-8a4f-62a571013283")
	)

	tests := []struct {
		name           string
		requestID      uuid.UUID
		mockRepo       func(repository *mock_account.MockRepository)
		expectedResult *model.BankAccount
		expectedError  error
	}{
		{
			name:      "Valid Delete Request",
			requestID: bankAccount.ID,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().DeleteBankAccount(ctx, bankAccount.ID).Return(bankAccount, nil)
			},
			expectedResult: bankAccount,
			expectedError:  nil,
		},
		{
			name:      "Not Found Request",
			requestID: notFoundId,
			mockRepo: func(repository *mock_account.MockRepository) {
				repository.EXPECT().DeleteBankAccount(ctx, notFoundId).Return(
					nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)))
			},
			expectedResult: nil,
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", 10329032)),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountServiceFixture(t)
			if tc.mockRepo != nil {
				tc.mockRepo(fixture.mockRepo)
			}

			result, err := fixture.service.DeleteBankAccount(ctx, tc.requestID)

			if tc.expectedError == nil {
				require.NoError(t, err)
				cmp.Equal(tc.expectedResult, result, cmpopts.IgnoreFields(model.BankAccount{}, "OpeningDate"))
			} else {
				assert.EqualError(t, tc.expectedError, err.Error())
				assert.Nil(t, result)
			}
		})
	}
}
