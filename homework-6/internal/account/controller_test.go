package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/dtos"
	mock_account "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/mocks"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app"
	app_mock "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/mocks"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bankAccountControllerFixture struct {
	ctrl        *gomock.Controller
	controller  app.AccountController
	mockService *mock_account.MockService
	mockLogger  *app_mock.MockLogger
}

func NewBankAccountControllerFixture(t *testing.T) *bankAccountControllerFixture {
	ctrl := gomock.NewController(t)
	mockService := mock_account.NewMockService(ctrl)
	mockLogger := app_mock.NewMockLogger(ctrl)
	controller := NewBankAccountController(mockService, mockLogger)
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
		bankAccountDto        = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		bankAccount           = fixtures.NewBankAccountBuilder().Valid().Build()
		createdBankAccount    = fixtures.NewBankAccountBuilder().Valid().Build()
		createdBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
	)

	tests := []struct {
		name            string
		requestPayload  dtos.BankAccountDto
		mockService     func(service *mock_account.MockService)
		mockLogger      func(logger *app_mock.MockLogger)
		expectedStatus  int
		expectedAccount dtos.BankAccountDto
	}{
		{
			name:           "Valid Request",
			requestPayload: *bankAccountDto,
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().CreateBankAccount(bankAccount).Return(createdBankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedStatus:  http.StatusOK,
			expectedAccount: *createdBankAccountDto,
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
			if tc.mockLogger != nil {
				tc.mockLogger(fixture.mockLogger)
			}

			requestPayload, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/bank-accounts", bytes.NewBuffer(requestPayload))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.controller.CreateBankAccount(w, r)
				assert.NoError(t, err)
			})

			handler.ServeHTTP(rr, req)

			var createdBankAccount dtos.BankAccountDto
			err = json.NewDecoder(rr.Body).Decode(&createdBankAccount)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			cmp.Equal(tc.expectedAccount, createdBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "ID", "OpeningDate"))
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	t.Parallel()

	var (
		bankAccount            = fixtures.NewBankAccountBuilder().Valid().Build()
		expectedBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		notFoundId, _          = uuid.Parse("331684ab-5af8-439a-8a4f-62a571013283")
	)

	tests := []struct {
		name             string
		requestAccountId uuid.UUID
		requestURL       string
		mockService      func(service *mock_account.MockService)
		mockLogger       func(logger *app_mock.MockLogger)
		expectedError    error
		expectedResult   dtos.BankAccountDto
	}{
		{
			name:             "Valid Request",
			requestAccountId: bankAccount.ID,
			requestURL:       fmt.Sprintf("/bank-accounts/%s", bankAccount.ID),
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().GetBankAccountById(bankAccount.ID).Return(bankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedError:  nil,
			expectedResult: *expectedBankAccountDto,
		},
		{
			name:             "Not Found Request",
			requestAccountId: notFoundId,
			requestURL:       fmt.Sprintf("/bank-accounts/%s", notFoundId),
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().GetBankAccountById(notFoundId).Return(
					nil, apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", notFoundId)))
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
				logger.EXPECT().Error(gomock.Any())
			},
			expectedError:  apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %d not found", notFoundId)),
			expectedResult: dtos.BankAccountDto{},
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
			if tc.mockLogger != nil {
				tc.mockLogger(fixture.mockLogger)
			}

			req, err := http.NewRequest("GET", tc.requestURL, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestAccountId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.controller.GetBankAccount(w, r)
				if tc.expectedError != nil {
					assert.EqualError(t, tc.expectedError, err.Error())
				} else {
					assert.NoError(t, err)
				}
			})

			handler.ServeHTTP(rr, req)

			if tc.expectedError != nil {
				assert.Empty(t, rr.Body)
			} else {
				var receivedBankAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&receivedBankAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedResult, receivedBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	t.Parallel()

	var (
		bankAccountDto        = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		bankAccount           = bankAccountDto.MapToModel()
		updatedBankAccount    = &bankAccount
		updatedBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
	)

	tests := []struct {
		name             string
		requestURL       string
		requestAccountId uuid.UUID
		requestPayload   dtos.BankAccountDto
		mockService      func(service *mock_account.MockService)
		mockLogger       func(logger *app_mock.MockLogger)
		expectedError    error
		expectedResult   dtos.BankAccountDto
	}{
		{
			name:             "Valid Request",
			requestURL:       fmt.Sprintf("/bank-accounts/%s", bankAccountDto.ID),
			requestAccountId: bankAccountDto.ID,
			requestPayload:   *bankAccountDto,
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().UpdateBankAccount(bankAccountDto.ID, &bankAccount).Return(updatedBankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedError:  nil,
			expectedResult: *updatedBankAccountDto,
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
			if tc.mockLogger != nil {
				tc.mockLogger(fixture.mockLogger)
			}

			requestPayload, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("PUT", tc.requestURL, bytes.NewBuffer(requestPayload))
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestAccountId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.controller.UpdateBankAccount(w, r)
				if tc.expectedError != nil {
					assert.EqualError(t, tc.expectedError, err.Error())
				} else {
					assert.NoError(t, err)
				}
			})

			handler.ServeHTTP(rr, req)

			if tc.expectedError != nil {
				assert.Empty(t, rr.Body)
			} else {
				var receivedBankAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&receivedBankAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedResult, receivedBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}
		})
	}
}

func TestDeleteBankAccount(t *testing.T) {
	t.Parallel()

	var (
		deletedBankAccount    = fixtures.NewBankAccountBuilder().Valid().Build()
		deletedBankAccountDto = fixtures.NewBankAccountDtoBuilder().Valid().Build()
	)

	tests := []struct {
		name             string
		requestURL       string
		requestAccountId uuid.UUID
		mockService      func(service *mock_account.MockService)
		mockLogger       func(logger *app_mock.MockLogger)
		expectedError    error
		expectedResult   dtos.BankAccountDto
	}{
		{
			name:             "Valid Request",
			requestURL:       fmt.Sprintf("/bank-accounts/%s", deletedBankAccount.ID),
			requestAccountId: deletedBankAccount.ID,
			mockService: func(service *mock_account.MockService) {
				service.EXPECT().DeleteBankAccount(deletedBankAccountDto.ID).Return(deletedBankAccount, nil)
			},
			mockLogger: func(logger *app_mock.MockLogger) {
				logger.EXPECT().Log(gomock.Any())
			},
			expectedError:  nil,
			expectedResult: *deletedBankAccountDto,
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
			if tc.mockLogger != nil {
				tc.mockLogger(fixture.mockLogger)
			}

			req, err := http.NewRequest("DELETE", tc.requestURL, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestAccountId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.controller.DeleteBankAccount(w, r)
				if tc.expectedError != nil {
					assert.EqualError(t, tc.expectedError, err.Error())
				} else {
					assert.NoError(t, err)
				}
			})

			handler.ServeHTTP(rr, req)

			if tc.expectedError != nil {
				assert.Empty(t, rr.Body)
			} else {
				var receivedBankAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&receivedBankAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedResult, receivedBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}
		})
	}
}
