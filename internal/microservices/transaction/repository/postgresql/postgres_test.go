package postgresql

import (
	"context"
	"database/sql"
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

func TestGetCount(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name       string
		userID     uuid.UUID
		returnRows *pgxmock.Rows
		expected   int
		errRows    error
		err        error
	}{
		{
			name:       "ValidCount",
			userID:     userID,
			returnRows: pgxmock.NewRows([]string{"count"}).AddRow(5),
			expected:   5,
			errRows:    nil,
			err:        nil,
		},
		{
			name:       "NoRows",
			userID:     userID,
			returnRows: pgxmock.NewRows([]string{"count"}),
			expected:   0,
			errRows:    sql.ErrNoRows,
			err:        fmt.Errorf("[repo] sql: no rows in result set"),
		},
		{
			name:       "ErrorRows",
			userID:     userID,
			returnRows: pgxmock.NewRows([]string{"count"}),
			expected:   0,
			errRows:    errors.New("error getting count"),
			err:        fmt.Errorf("[repo] error getting count"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionCount)
			mock.ExpectQuery(escapedQuery).
				WithArgs(test.userID).
				WillReturnRows(test.returnRows).
				WillReturnError(test.errRows)

			count, err := repo.GetCount(context.Background(), test.userID)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if count != test.expected {
				t.Errorf("Expected count: %d, but got: %d", test.expected, count)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetFeed(t *testing.T) {
	userID := uuid.New()
	transactionID1 := uuid.New()
	categoryID := uuid.New()
	time := time.Now()
	categories := []models.CategoryName{{ID: categoryID, Name: "ffdsf"}}
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
				"category_id",
				"name",
			}).AddRow(
				categoryID,
				"ffdsf",
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
				"category_id",
				"name",
			}).AddRow(
				"dff",
				"sfd",
			),

			err:             fmt.Errorf("[repo] Scanning value error for column 'category_id': Scan: invalid UUID length: 3"),
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
				"category_id",
				"name",
			}).AddRow(
				categoryID,
				"dddd",
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
				"category_id",
				"name",
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
				"category_id",
				"name",
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
				"category_id",
				"name",
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
				"name",
			}).AddRow(
				uuid.New(),
				"ffff",
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
			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionGetFeed + " ORDER BY date DESC;")
			mock.ExpectQuery(escapedQuery).
				WithArgs(userID.String()).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			if test.errTransaction {
				escapedQueryCategory := regexp.QuoteMeta(transactionGetCategory)
				mock.ExpectQuery(escapedQueryCategory).
					WithArgs(transactionID1).
					WillReturnRows(test.rowsCategory).
					WillReturnError(test.rowsCategoryErr)
			}
			transactions, err := repo.GetFeed(context.Background(), userID, &models.QueryListOptions{})

			if !reflect.DeepEqual(transactions, test.expected) {
				t.Errorf("Expected transactions: %v, but got: %v", test.expected, transactions)
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

			logger := *logger.NewLogger(context.TODO())
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

func TestUpdateAccountBalance(t *testing.T) {
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionUpdateAccount)
			amount := 10.0
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			mock.ExpectExec(escapedQuery).
				WithArgs(amount, transactionID).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1))

			//mock.ExpectCommit()
			err := repo.updateAccountBalance(context.Background(), mock, test.returnRows, amount)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateAccountBalances(t *testing.T) {
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
			transaction: models.Transaction{Income: 10, Outcome: 10, AccountIncomeID: transactionID, AccountOutcomeID: transactionID},
			returnRows:  transactionID,
			expected:    transactionID,
			err:         fmt.Errorf("[repo] failed to update old AccountIncome balance: all expectations were already fulfilled, call to ExecQuery 'UPDATE accounts SET balance = balance - $1 WHERE id = $2;' with args [10 %s] was not expected", transactionID.String()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionUpdateAccount)
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			mock.ExpectExec(escapedQuery).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1))

			//mock.ExpectCommit()
			err := repo.updateAccountBalances(context.Background(), mock, &test.transaction)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestInsertCategories(t *testing.T) {
	transactionID := uuid.New()
	tests := []struct {
		name          string
		errRows       error
		transactionID uuid.UUID
		categories    []models.CategoryName
		err           error
		rowsErr       error
	}{
		{
			name:          "ValidCategories",
			transactionID: transactionID,
			categories:    []models.CategoryName{{ID: transactionID}},
			err:           nil,
			rowsErr:       nil,
		},
		{
			name:          "Error",
			transactionID: transactionID,
			categories:    []models.CategoryName{{ID: transactionID}},
			err:           fmt.Errorf("[repo] failed to insert category association: %w", errors.New("err")),
			rowsErr:       errors.New("err"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionCreateCategory)
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			mock.ExpectExec(escapedQuery).
				WithArgs(transactionID, transactionID).
				WillReturnResult(pgxmock.NewResult("INSERT", 1)).
				WillReturnError(test.rowsErr)

			//mock.ExpectCommit()
			err := repo.insertCategories(context.Background(), mock, transactionID, test.categories)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTransactionInfo(t *testing.T) {
	tests := []struct {
		name        string
		errRows     error
		transaction models.Transaction
		err         error
		rowsErr     error
	}{
		{
			name:        "ValidCategories",
			transaction: models.Transaction{},
			err:         nil,
			rowsErr:     nil,
		},
		{
			name:        "Error",
			transaction: models.Transaction{},
			err:         fmt.Errorf("[repo] failed to update transaction information: err"),
			rowsErr:     errors.New("err"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionUpdate)
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			// "UPDATE transaction set account_income=$2, account_outcome=$3, income=$4, outcome=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"

			mock.ExpectExec(escapedQuery).
				WithArgs(test.transaction.AccountIncomeID,
					test.transaction.AccountIncomeID,
					test.transaction.AccountOutcomeID,
					test.transaction.Income,
					test.transaction.Outcome,
					test.transaction.Date,
					test.transaction.Payer,
					test.transaction.Description).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(test.rowsErr)

			//mock.ExpectCommit()
			err := repo.updateTransactionInfo(context.Background(), mock, &test.transaction)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteExistingCategoryAssociations(t *testing.T) {
	transactionID := uuid.New()
	tests := []struct {
		name        string
		errRows     error
		transaction models.Transaction
		err         error
		rowsErr     error
	}{
		{
			name:        "ValidCategories",
			transaction: models.Transaction{},
			err:         nil,
			rowsErr:     nil,
		},
		{
			name:        "Error",
			transaction: models.Transaction{},
			err:         fmt.Errorf("[repo] failed to delete existing category associations: err"),
			rowsErr:     errors.New("err"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionDeleteCategory)
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			// "UPDATE transaction set account_income=$2, account_outcome=$3, income=$4, outcome=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"

			mock.ExpectExec(escapedQuery).
				WithArgs(transactionID).
				WillReturnResult(pgxmock.NewResult("DELETE", 1)).
				WillReturnError(test.rowsErr)

			//mock.ExpectCommit()
			err := repo.deleteExistingCategoryAssociations(context.Background(), mock, transactionID)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAccountBalance(t *testing.T) {
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
			transaction: models.Transaction{Income: 10, Outcome: 10, AccountIncomeID: transactionID, AccountOutcomeID: transactionID},
			returnRows:  transactionID,
			expected:    transactionID,
			err:         fmt.Errorf("[repo] failed to update old AccountIncome balance: err"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionUpdateAccount)
			//"UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
			mock.ExpectExec(escapedQuery).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(errors.New("err"))

			//mock.ExpectCommit()
			err := repo.deleteAccountBalance(context.Background(), mock, 10.0, 10.0, transactionID, transactionID)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCheckForbidden(t *testing.T) {
	transactionID := uuid.New()
	tests := []struct {
		name        string
		transaction uuid.UUID
		errRows     error

		returnRows *pgxmock.Rows
		expected   uuid.UUID
		err        error
	}{
		{
			name:        "ValidTransaction",
			transaction: transactionID,
			returnRows:  pgxmock.NewRows([]string{"user_id"}).AddRow(transactionID),
			expected:    transactionID,
			errRows:     nil,
			err:         nil,
		},
		{
			name:        "ValidTransaction",
			transaction: transactionID,
			returnRows:  pgxmock.NewRows([]string{"user_id"}).AddRow(transactionID),
			expected:    transactionID,
			errRows:     sql.ErrNoRows,
			err:         fmt.Errorf("[repo] No Such transaction: %s doesn't exist: sql: no rows in result set", transactionID.String()),
		},
		{
			name:        "ValidTransaction",
			transaction: transactionID,
			returnRows:  pgxmock.NewRows([]string{"user_id"}).AddRow(transactionID),
			expected:    transactionID,
			errRows:     errors.New("err"),
			err:         fmt.Errorf("[repo] failed request db SELECT user_id FROM transaction WHERE id = $1;, err"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(TransactionGetUserByID)
			mock.ExpectQuery(escapedQuery).
				WithArgs(test.transaction).
				WillReturnRows(test.returnRows).
				WillReturnError(test.errRows)

			_, err := repo.CheckForbidden(context.Background(), test.transaction)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
func TestGetExportInfo(t *testing.T) {
	userID := uuid.New()
	transactionID1 := uuid.New()
	categoryID := uuid.New()
	time := time.Now()
	categories := []models.CategoryName{{ID: categoryID, Name: "ffdsf"}}
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
				"category_id",
				"name",
			}).AddRow(
				categoryID,
				"ffdsf",
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
				"category_id",
				"name",
			}).AddRow(
				"dff",
				"sfd",
			),

			err:             fmt.Errorf("[repo] Scanning value error for column 'category_id': Scan: invalid UUID length: 3"),
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
				"category_id",
				"name",
			}).AddRow(
				categoryID,
				"dddd",
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
				"category_id",
				"name",
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
				"category_id",
				"name",
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
				"category_id",
				"name",
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
				"name",
			}).AddRow(
				uuid.New(),
				"ffff",
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
			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(transactionGetFeed + " ORDER BY date DESC;")
			mock.ExpectQuery(escapedQuery).
				WithArgs(userID.String()).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			if test.errTransaction {
				escapedQueryCategory := regexp.QuoteMeta(transactionGetFeedForExport)
				mock.ExpectQuery(escapedQueryCategory).
					WithArgs(transactionID1).
					WillReturnRows(test.rowsCategory).
					WillReturnError(test.rowsCategoryErr)
			}
			transactions, err := repo.GetTransactionForExport(context.Background(), userID, &models.QueryListOptions{})

			if !reflect.DeepEqual(transactions, test.expected) {
				t.Errorf("Expected transactions: %v, but got: %v", test.expected, transactions)
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
*/
