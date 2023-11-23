//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/dtos"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BankAccountServiceFixture struct {
	coreController *app.Controller
}

func NewBankAccountServiceFixture() *BankAccountServiceFixture {
	InitDB()

	bankAccountRepository := account.NewBankAccountRepository(db.DB)
	bankAccountService := account.NewBankAccountService(bankAccountRepository)
	bankAccountController := account.NewBankAccountController(bankAccountService)

	subscriptionRepository := subscription.NewSubscriptionRepository(db.DB)
	subscriptionService := subscription.NewSubscriptionService(subscriptionRepository)
	subscriptionController := subscription.NewSubscriptionController(subscriptionService)

	coreController := app.NewCoreController(bankAccountController, subscriptionController)

	return &BankAccountServiceFixture{
		coreController: coreController,
	}
}

func TestCreateBankAccount(t *testing.T) {
	var (
		validBankAccountDto   = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		invalidBankAccountDto = fixtures.NewBankAccountDtoBuilder().Invalid().Build()
	)

	tests := []struct {
		name            string
		requestPayload  dtos.BankAccountDto
		expectedError   error
		expectedAccount dtos.BankAccountDto
		idExists        bool
	}{
		{
			name:            "Valid create request",
			requestPayload:  *validBankAccountDto,
			expectedError:   nil,
			expectedAccount: *validBankAccountDto,
			idExists:        true,
		},
		{
			name:            "Invalid create request",
			requestPayload:  *invalidBankAccountDto,
			expectedError:   apperr.NewBadRequestError("Balance: must be no less than 0; BankName: cannot be blank; HolderName: must be in a valid format."),
			expectedAccount: dtos.BankAccountDto{},
			idExists:        false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fixture := NewBankAccountServiceFixture()
			db.SetUp(t)
			defer func() {
				err := db.DeleteBankAccounts([]uuid.UUID{tc.requestPayload.ID})
				if err != nil {
					panic(fmt.Sprintf("Failed to clear data after test: %s", err))
				}
				db.TearDown()
			}()

			requestBody, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/bank-accounts", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.coreController.AccountController.CreateBankAccount(w, r)
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
				var createdAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&createdAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedAccount, createdAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}

			exists, err := db.ExistsById(tc.requestPayload.ID, "bank_account")
			assert.NoError(t, err)
			assert.Equal(t, tc.idExists, exists)
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	var (
		validBankAccountDto   = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		invalidBankAccountDto = fixtures.NewBankAccountDtoBuilder().Invalid().Build()
	)

	tests := []struct {
		name            string
		requestId       uuid.UUID
		expectedAccount dtos.BankAccountDto
		expectedError   error
		isExist         bool
	}{
		{
			name:            "Valid get request",
			requestId:       validBankAccountDto.ID,
			expectedAccount: *validBankAccountDto,
			expectedError:   nil,
			isExist:         true,
		},
		{
			name:            "Not found get request",
			requestId:       invalidBankAccountDto.ID,
			expectedAccount: dtos.BankAccountDto{},
			expectedError:   apperr.NewNotFoundError(fmt.Sprintf("Bank account with ID: %s not found", invalidBankAccountDto.ID.String())),
			isExist:         false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fixture := NewBankAccountServiceFixture()
			db.SetUp(t)
			defer db.TearDown()

			if tc.isExist {
				err := db.InsertSampleData([]model.BankAccount{tc.expectedAccount.MapToModel()})
				assert.NoError(t, err, "Failed to insert test data into the database")
				defer func() {
					err = db.DeleteBankAccounts([]uuid.UUID{tc.requestId})
					if err != nil {
						panic(fmt.Sprintf("Failed to clear data after test: %s", err))
					}
				}()
			}

			req, err := http.NewRequest("GET", fmt.Sprintf("/bank-accounts/%s", tc.requestId), nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.coreController.AccountController.GetBankAccount(w, r)
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
				var receivedAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&receivedAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedAccount, receivedAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}

			exists, err := db.ExistsById(tc.requestId, "bank_account")
			assert.NoError(t, err)
			assert.Equal(t, tc.isExist, exists)
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	var (
		validBankAccountDto      = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		validUpdatedAccountDto   = fixtures.NewBankAccountDtoBuilder().Valid().Balance(9500).Build()
		notExistingAccountId     = uuid.New()
		invalidUpdatedAccountDto = fixtures.NewBankAccountDtoBuilder().Invalid().Balance(5000).Build()
	)

	tests := []struct {
		name            string
		requestId       uuid.UUID
		oldAccount      dtos.BankAccountDto
		requestPayload  dtos.BankAccountDto
		expectedError   error
		expectedAccount dtos.BankAccountDto
		idExists        bool
	}{
		{
			name:            "Valid update request",
			requestId:       validBankAccountDto.ID,
			oldAccount:      *validBankAccountDto,
			requestPayload:  *validUpdatedAccountDto,
			expectedError:   nil,
			expectedAccount: *validUpdatedAccountDto,
			idExists:        true,
		},
		{
			name:            "Not found update request",
			requestId:       notExistingAccountId,
			requestPayload:  *validBankAccountDto,
			expectedError:   apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %s not found", notExistingAccountId.String())),
			expectedAccount: dtos.BankAccountDto{},
			idExists:        false,
		},
		{
			name:            "Invalid update request",
			requestId:       validBankAccountDto.ID,
			oldAccount:      *validBankAccountDto,
			requestPayload:  *invalidUpdatedAccountDto,
			expectedError:   apperr.NewBadRequestError("BankName: cannot be blank; HolderName: must be in a valid format."),
			expectedAccount: dtos.BankAccountDto{},
			idExists:        true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fixture := NewBankAccountServiceFixture()
			db.SetUp(t)
			defer db.TearDown()

			if tc.idExists {
				err := db.InsertSampleData([]model.BankAccount{tc.oldAccount.MapToModel()})
				assert.NoError(t, err, "Failed to insert test data into the database")
				defer func() {
					err = db.DeleteBankAccounts([]uuid.UUID{tc.requestId})
					if err != nil {
						panic(fmt.Sprintf("Failed to clear data after test: %s", err))
					}
				}()
			}

			exists, err := db.ExistsById(tc.requestId, "bank_account")
			assert.NoError(t, err)
			assert.Equal(t, tc.idExists, exists)

			requestBody, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("PUT", fmt.Sprintf("/bank-accounts/%s", tc.requestId), bytes.NewBuffer(requestBody))
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.coreController.AccountController.UpdateBankAccount(w, r)
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
				var updatedBankAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&updatedBankAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedAccount, updatedBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}
		})
	}
}

func TestDeleteBankAccount(t *testing.T) {
	var (
		validBankAccountDto  = fixtures.NewBankAccountDtoBuilder().Valid().Build()
		notExistinsAccountId = uuid.New()
	)

	tests := []struct {
		name            string
		requestId       uuid.UUID
		expectedError   error
		expectedAccount dtos.BankAccountDto
		idExists        bool
		idExistsAfter   bool
	}{
		{
			name:            "Valid delete request",
			requestId:       validBankAccountDto.ID,
			expectedError:   nil,
			expectedAccount: *validBankAccountDto,
			idExists:        true,
			idExistsAfter:   false,
		},
		{
			name:            "Not found delete request",
			requestId:       notExistinsAccountId,
			expectedError:   apperr.NewNotFoundError(fmt.Sprintf("bank account with ID: %s does not exist", notExistinsAccountId.String())),
			expectedAccount: dtos.BankAccountDto{},
			idExists:        false,
			idExistsAfter:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fixture := NewBankAccountServiceFixture()
			db.SetUp(t)
			defer db.TearDown()

			if tc.idExists {
				err := db.InsertSampleData([]model.BankAccount{tc.expectedAccount.MapToModel()})
				assert.NoError(t, err, "Failed to insert test data into the database")
				defer func() {
					err = db.DeleteBankAccounts([]uuid.UUID{tc.requestId})
					if err != nil {
						panic(fmt.Sprintf("Failed to clear data after test: %s", err))
					}
				}()
			}

			exists, err := db.ExistsById(tc.requestId, "bank_account")
			assert.NoError(t, err)
			assert.Equal(t, tc.idExists, exists)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/bank-accounts/%s", tc.requestId), nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestId.String()})
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fixture.coreController.AccountController.DeleteBankAccount(w, r)
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
				var deletedBankAccount dtos.BankAccountDto
				err = json.NewDecoder(rr.Body).Decode(&deletedBankAccount)
				assert.NoError(t, err)
				cmp.Equal(tc.expectedAccount, deletedBankAccount, cmpopts.IgnoreFields(dtos.BankAccountDto{}, "OpeningDate"))
			}

			exists, err = db.ExistsById(tc.requestId, "bank_account")
			assert.NoError(t, err)
			assert.Equal(t, tc.idExistsAfter, exists)
		})
	}
}
