package account

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/database"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/tests/fixtures"
	"testing"
)

type bankAccountRepoFixture struct {
	mockSqlDb *sqlmock.Sqlmock
	repo      Repository
}

func NewBankAccountRepoFixture(t *testing.T) (*bankAccountRepoFixture, error) {
	mockSqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	sqlDatabase, err := database.InitDBWithPool(mockSqlDb)
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	repo := NewBankAccountRepository(sqlDatabase)
	return &bankAccountRepoFixture{
		mockSqlDb: &mock,
		repo:      repo,
	}, nil
}

func TestCreateBankAccountRepo(t *testing.T) {
	t.Parallel()

	var (
		bankAccount        = fixtures.NewBankAccountBuilder().Valid().Build()
		invalidBankAccount = fixtures.NewBankAccountBuilder().Valid().HolderName("").Build()
	)

	tests := []struct {
		name            string
		bankAccount     model.BankAccount
		mockSQL         func(mock sqlmock.Sqlmock)
		expectedAccount *model.BankAccount
		expectedError   error
	}{
		{
			name:        "Success",
			bankAccount: *bankAccount,
			mockSQL: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO bank_account \(id, holder_name, balance, opening_date, bank_name\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
					WithArgs(bankAccount.ID, bankAccount.HolderName, bankAccount.Balance, sqlmock.AnyArg(), bankAccount.BankName).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedAccount: fixtures.NewBankAccountBuilder().Valid().Build(),
			expectedError:   nil,
		},
		{
			name:        "Fail, HolderName is null",
			bankAccount: *invalidBankAccount,
			mockSQL: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`^INSERT INTO bank_account \(id, holder_name, balance, opening_date, bank_name\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id$`).
					WithArgs(sqlmock.AnyArg(), "", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("null value in column \"holder_name\" violates not-null constraint"))
				mock.ExpectRollback()
			},
			expectedAccount: nil,
			expectedError:   apperr.NewInternalServerError("Internal server error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture, err := NewBankAccountRepoFixture(t)
			if err != nil {
				t.Fatalf("Error setting up test fixture: %v", err)
			}

			tc.mockSQL(*fixture.mockSqlDb)

			createdAccount, err := fixture.repo.CreateBankAccount(&tc.bankAccount)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccount, createdAccount)
			}

			err = (*fixture.mockSqlDb).ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
