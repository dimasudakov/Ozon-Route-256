package account

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBankAccountService struct{}

func (m *MockBankAccountService) CreateBankAccount(account BankAccount) (*BankAccount, error) {
	return &BankAccount{
		ID:         1,
		HolderName: account.HolderName,
		Balance:    account.Balance,
		BankName:   account.BankName,
	}, nil
}

func (m *MockBankAccountService) GetBankAccountById(id int) (*BankAccount, error) {
	return &BankAccount{
		ID:         1,
		HolderName: "Dima Sudakov",
		Balance:    1000,
		BankName:   "Tinkoff",
	}, nil
}

func (m *MockBankAccountService) UpdateBankAccount(id int, account BankAccount) (*BankAccount, error) {
	return &BankAccount{
		ID:         1,
		HolderName: "Dima Sudakov",
		Balance:    1500,
		BankName:   "Sberbank",
	}, nil
}

func (m *MockBankAccountService) DeleteBankAccount(id int) (*BankAccount, error) {
	return &BankAccount{
		ID:         1,
		HolderName: "Dima Sudakov",
		Balance:    1000,
		BankName:   "Sberbank",
	}, nil
}

func TestCreateBankAccount(t *testing.T) {
	controller := BankAccountController{
		service: &MockBankAccountService{},
	}

	tests := []struct {
		name           string
		requestPayload BankAccountDto
		expectedStatus int
		expectedResult string
	}{
		{
			name: "Valid Request",
			requestPayload: BankAccountDto{
				HolderName: "Dima Sudakov",
				Balance:    1000,
				BankName:   "Tinkoff",
			},
			expectedStatus: http.StatusOK,
			expectedResult: `{"id":1,"holder_name":"Dima Sudakov","balance":1000,"opening_date":"0001-01-01T00:00:00Z","bank_name":"Tinkoff","subscriptions":[]}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			requestPayload, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/bank-accounts", bytes.NewBuffer(requestPayload))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := controller.CreateBankAccount(w, r)
				assert.NoError(t, err)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedResult, rr.Body.String())
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	controller := BankAccountController{
		service: &MockBankAccountService{},
	}

	tests := []struct {
		name           string
		requestURL     string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Valid Request",
			requestURL:     "/bank-accounts/1",
			expectedStatus: http.StatusOK,
			expectedResult: `{"id":1,"holder_name":"Dima Sudakov","balance":1000,"opening_date":"0001-01-01T00:00:00Z","bank_name":"Tinkoff","subscriptions":[]}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.requestURL, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := controller.GetBankAccount(w, r)
				assert.NoError(t, err)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedResult, rr.Body.String())
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	controller := BankAccountController{
		service: &MockBankAccountService{},
	}

	tests := []struct {
		name           string
		requestURL     string
		requestPayload BankAccountDto
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Valid Request",
			requestURL:     "/bank-accounts/1",
			requestPayload: BankAccountDto{HolderName: "Dima Sudakov", Balance: 1500, BankName: "Sberbank"},
			expectedStatus: http.StatusOK,
			expectedResult: `{"id":1,"holder_name":"Dima Sudakov","balance":1500,"opening_date":"0001-01-01T00:00:00Z","bank_name":"Sberbank","subscriptions":[]}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			requestPayload, err := json.Marshal(tc.requestPayload)
			assert.NoError(t, err)

			req, err := http.NewRequest("PUT", tc.requestURL, bytes.NewBuffer(requestPayload))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := controller.UpdateBankAccount(w, r)
				assert.NoError(t, err)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedResult, rr.Body.String())
		})
	}
}

func TestDeleteBankAccount(t *testing.T) {
	controller := BankAccountController{
		service: &MockBankAccountService{},
	}

	tests := []struct {
		name           string
		requestURL     string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "Valid Request",
			requestURL:     "/bank-accounts/1",
			expectedStatus: http.StatusOK,
			expectedResult: `{"id":1,"holder_name":"Dima Sudakov","balance":1000,"opening_date":"0001-01-01T00:00:00Z","bank_name":"Sberbank","subscriptions":[]}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", tc.requestURL, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := controller.DeleteBankAccount(w, r)
				assert.NoError(t, err)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedResult, rr.Body.String())
		})
	}
}
