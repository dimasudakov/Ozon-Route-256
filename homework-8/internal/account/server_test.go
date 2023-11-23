package account

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock_account "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/mocks"
	app_mock "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/mocks"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/bank_accounts"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/fixtures"
	"testing"
)

type bankAccountControllerFixture struct {
	ctrl        *gomock.Controller
	controller  bank_accounts.BankAccountServiceServer
	mockService *mock_account.MockService
	mockLogger  *app_mock.MockLogger
}

func NewBankAccountControllerFixture(t *testing.T) *bankAccountControllerFixture {
	ctrl := gomock.NewController(t)
	mockService := mock_account.NewMockService(ctrl)
	mockLogger := app_mock.NewMockLogger(ctrl)
	controller := NewBankAccountGrpcImpl(mockService)
	return &bankAccountControllerFixture{
		ctrl:        ctrl,
		controller:  controller,
		mockService: mockService,
		mockLogger:  mockLogger,
	}
}

func TestCreateBankAccount(t *testing.T) {
	t.Parallel()

	var (
		ctx                   = context.Background()
		bankAccountDto        = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		bankAccount           = fixtures.NewBankAccountBuilder().Valid().Build()
		createdBankAccount    = fixtures.NewBankAccountBuilder().Valid().Build()
		createdBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		invalidBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().ID("invalid id").Build()
	)

	tests := []struct {
		name            string
		requestPayload  bank_accounts.BankAccountDto
		mockService     func(service *mock_account.MockService)
		mockLogger      func(logger *app_mock.MockLogger)
		expectedError   error
		expectedAccount *bank_accounts.BankAccountDto
	}{
		{
			name:           "Valid Request",
			requestPayload: *bankAccountDto,
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().CreateBankAccount(ctx, bankAccount).Return(createdBankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedError:   nil,
			expectedAccount: createdBankAccountDto,
		},
		{
			name:            "Fail, invalid ID",
			requestPayload:  *invalidBankAccountDto,
			expectedError:   errors.New("invalid id"),
			expectedAccount: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountControllerFixture(t)
			if tc.mockService != nil {
				tc.mockService(fixture.mockService)
			}

			createBankAccountRequest := &bank_accounts.CreateBankAccountRequest{
				Account: &tc.requestPayload,
			}
			result, err := fixture.controller.CreateBankAccount(ctx, createBankAccountRequest)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			cmp.Equal(tc.expectedAccount, result.GetAccount(), cmpopts.IgnoreFields(bank_accounts.BankAccountDto{},
				"Id", "OpeningDate", "state", "sizeCache", "unknownFields"))
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	t.Parallel()

	var (
		ctx                    = context.Background()
		bankAccount            = fixtures.NewBankAccountBuilder().Valid().Build()
		expectedBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		notFoundId, _          = uuid.Parse("331684ab-5af8-439a-8a4f-62a571013283")
	)

	tests := []struct {
		name             string
		requestAccountId string
		mockService      func(service *mock_account.MockService)
		mockLogger       func(logger *app_mock.MockLogger)
		expectedError    error
		expectedResult   *bank_accounts.BankAccountDto
	}{
		{
			name:             "Valid Request",
			requestAccountId: bankAccount.ID.String(),
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().GetBankAccountById(ctx, bankAccount.ID).Return(bankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedError:  nil,
			expectedResult: expectedBankAccountDto,
		},
		{
			name:             "Not Found Request",
			requestAccountId: notFoundId.String(),
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().GetBankAccountById(ctx, notFoundId).Return(
					nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", notFoundId)))
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
				logger.EXPECT().Error(gomock.Any())
			},
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", notFoundId)),
			expectedResult: &bank_accounts.BankAccountDto{},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fixture := NewBankAccountControllerFixture(t)
			if tc.mockService != nil {
				tc.mockService(fixture.mockService)
			}

			getBankAccountRequest := &bank_accounts.GetBankAccountByIdRequest{
				Id: &bank_accounts.UUID{Value: tc.requestAccountId},
			}
			receivedBankAccount, err := fixture.controller.GetBankAccountById(ctx, getBankAccountRequest)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				cmp.Equal(tc.expectedResult, receivedBankAccount, cmpopts.IgnoreFields(bank_accounts.BankAccountDto{},
					"OpeningDate", "state", "sizeCache", "unknownFields"))
			}
		})
	}
}
