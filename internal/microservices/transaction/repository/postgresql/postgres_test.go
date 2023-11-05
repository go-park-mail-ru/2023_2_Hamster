package postgresql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func TestGetFeed(t *testing.T) {
	userID := uuid.New()
	transactionID1 := uuid.New()
	categoryID := uuid.New()
	time := time.Now()
	categories := []uuid.UUID{categoryID}
	tests := []struct {
		name            string
		rows            *pgxmock.Rows
		rowsCategory    *pgxmock.Rows
		rowsCategoryErr error
		err             error
		rowsErr         error
		expected        []models.Transaction
		expectedLast    bool
		errTransaction  bool
	}{
		{
			name: "ValidFeed",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).AddRow(
				transactionID1, userID, transactionID1, transactionID1, 100.0, 0.0, time, "John Doe", "Transaction 1",
			),
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}).AddRow(
				categoryID,
			),

			err:             nil,
			rowsErr:         nil,
			rowsCategoryErr: nil,
			expected:        []models.Transaction{{ID: transactionID1, UserID: userID, AccountIncomeID: transactionID1, AccountOutcomeID: transactionID1, Income: 100.0, Outcome: 0.0, Date: time, Payer: "John Doe", Description: "Transaction 1", Categories: categories}},
			expectedLast:    false,
			errTransaction:  true,
		},
		{
			name: "Invalid Scan Category",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).AddRow(
				transactionID1, userID, transactionID1, transactionID1, 100.0, 0.0, time, "John Doe", "Transaction 1",
			),
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}).AddRow(
				"sfd",
			),

			err:             fmt.Errorf("[repo] Scanning value error for column 'catygory_id': Scan: invalid UUID length: 3"),
			rowsErr:         nil,
			rowsCategoryErr: nil,
			expected:        nil,
			expectedLast:    false,
			errTransaction:  true,
		},
		{
			name: "Invalid Scan transaction",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).AddRow(
				"fff", userID, transactionID1, transactionID1, 100.0, 0.0, time, "John Doe", "Transaction 1",
			),

			err:             fmt.Errorf("[repo] Scanning value error for column 'id': Scan: invalid UUID length: 3"),
			rowsErr:         nil,
			rowsCategoryErr: nil,
			expected:        nil,
			expectedLast:    false,
			errTransaction:  false,
		},
		{
			name: "INValid category",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).AddRow(
				transactionID1, userID, transactionID1, transactionID1, 100.0, 0.0, time, "John Doe", "Transaction 1",
			),
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}).AddRow(
				categoryID,
			),

			err:             fmt.Errorf("[repo] err"),
			rowsErr:         nil,
			rowsCategoryErr: errors.New("err"),
			expected:        nil,
			expectedLast:    false,
			errTransaction:  true,
		},
		{
			name: "Rows err",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).RowError(0, errors.New("err")),
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}),

			err:             fmt.Errorf("[repo] err"),
			rowsErr:         nil,
			rowsCategoryErr: nil,
			expected:        nil,
			expectedLast:    false,
			errTransaction:  false,
		},
		{
			name: "Rows err",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}).AddRow(
				transactionID1, userID, transactionID1, transactionID1, 100.0, 0.0, time, "John Doe", "Transaction 1",
			),
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}).RowError(0, errors.New("err")),
			err:             fmt.Errorf("[repo] err"),
			rowsErr:         nil,
			rowsCategoryErr: nil,
			expected:        nil,
			expectedLast:    false,
			errTransaction:  true,
		},
		{
			name: "NoRows",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}),
			rowsErr:         nil,
			rowsCategoryErr: nil,
			rowsCategory: pgxmock.NewRows([]string{
				"catygory_id",
			}),
			err:            fmt.Errorf("[repo] %w: <nil>", &models.NoSuchTransactionError{UserID: userID}),
			expected:       nil,
			expectedLast:   false,
			errTransaction: false,
		},
		{
			name: "DatabaseError",
			rows: pgxmock.NewRows([]string{
				"id", "user_id", "account_income_id", "account_outcome_id", "income", "outcome", "date", "payer", "description",
			}),
			rowsErr: errors.New("err"),
			rowsCategory: pgxmock.NewRows([]string{
				"category_id",
			}).AddRow(
				uuid.New(),
			),
			rowsCategoryErr: errors.New("err"),
			err:             fmt.Errorf("[repo] err"),
			expected:        nil,
			expectedLast:    false,
			errTransaction:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			pageSize, page := 1, 5
			logger := *logger.InitLogger()
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionGetFeed)
			mock.ExpectQuery(escapedQuery).
				WithArgs(userID, pageSize, (page-1)*pageSize).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			if test.errTransaction {
				escapedQueryCategory := regexp.QuoteMeta(transactionGetCategory)
				mock.ExpectQuery(escapedQueryCategory).
					WithArgs(transactionID1).
					WillReturnRows(test.rowsCategory).
					WillReturnError(test.rowsCategoryErr)
			}
			transactions, last, err := repo.GetFeed(context.Background(), userID, page, pageSize)

			if !reflect.DeepEqual(transactions, test.expected) {
				t.Errorf("Expected transactions: %v, but got: %v", test.expected, transactions)
			}

			if last != test.expectedLast {
				t.Errorf("Expected last: %v, but got: %v", test.expectedLast, last)
			}

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestInsertTransaction(t *testing.T) {
	transactionID := uuid.New()
	tests := []struct {
		name        string
		transaction models.Transaction
		errRows     error
		returnRows  uuid.UUID
		expected    uuid.UUID
		err         error
	}{
		{
			name:        "ValidTransaction",
			transaction: models.Transaction{},
			returnRows:  transactionID,
			expected:    transactionID,
			err:         nil,
		},
		{
			name:        "InvalidTransaction",
			transaction: models.Transaction{},

			errRows:    errors.New("Invalid user data"),
			returnRows: uuid.Nil,
			err:        errors.New("[repo] failed create transaction: Invalid user data"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.InitLogger()
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionCreate)
			//mock.ExpectBegin()
			mock.ExpectQuery(escapedQuery).WithArgs(test.transaction.UserID,
				test.transaction.AccountIncomeID,
				test.transaction.AccountOutcomeID,
				test.transaction.Income,
				test.transaction.Outcome,
				test.transaction.Date,
				test.transaction.Payer,
				test.transaction.Description).
				WillReturnError(test.errRows).
				WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(test.returnRows))
			//mock.ExpectCommit()

			transactinID, err := repo.insertTransaction(context.Background(), mock, &test.transaction)
			if !reflect.DeepEqual(transactinID, test.expected) {
				t.Errorf("Expected transactions: %v, but got: %v", test.expected, transactinID)
			}

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
