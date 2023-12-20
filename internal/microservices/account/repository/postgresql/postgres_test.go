package postgresql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func Test_SharingCheck(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()
	testCases := []struct {
		name        string
		rows        *pgxmock.Rows
		rowsError   error
		expectedErr error
	}{
		{
			name: "Success",
			rows: pgxmock.NewRows([]string{"count"}).
				AddRow(1),
			rowsError:   nil,
			expectedErr: nil,
		},
		{
			name: "ForbiddenUserError",
			rows: pgxmock.NewRows([]string{"count"}).
				AddRow(0),
			rowsError:   nil,
			expectedErr: fmt.Errorf("[repo] failed %w", &models.ForbiddenUserError{}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(AccountSharingCheck)
			mock.ExpectQuery(escapedQuery).
				WithArgs(userID, accountID).
				WillReturnRows(tc.rows).
				WillReturnError(tc.rowsError)

			err := repo.SharingCheck(context.Background(), accountID, userID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_Unsubscribe(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name        string
		execResult  pgconn.CommandTag
		execError   error
		expectedErr error
	}{
		{
			name:        "Success",
			execResult:  pgconn.CommandTag{},
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Error",
			execResult:  pgconn.CommandTag{},
			execError:   errors.New("Some error"),
			expectedErr: fmt.Errorf("[repo] failed to delete from UserAccount table: %w", errors.New("Some error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the Unsubscribe query with the specified parameters
			mock.ExpectExec(regexp.QuoteMeta(Unsubscribe)).
				WithArgs(accountID, userID).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)

			err := repo.Unsubscribe(context.Background(), userID, accountID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_DeleteUserInAccount(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name        string
		execResult  pgconn.CommandTag
		execError   error
		expectedErr error
	}{
		{
			name:        "Success",
			execResult:  pgconn.CommandTag{},
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Error",
			execResult:  pgconn.CommandTag{},
			execError:   errors.New("Some error"),
			expectedErr: fmt.Errorf("[repo] failed to delete from UserAccount table: %w", errors.New("Some error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the Unsubscribe query with the specified parameters
			mock.ExpectExec(regexp.QuoteMeta(Unsubscribe)).
				WithArgs(accountID, userID).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)

			err := repo.DeleteUserInAccount(context.Background(), userID, accountID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_CheckForbidden(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name         string
		rows         *pgxmock.Rows
		rowsError    error
		expectedErr  error
		expectedBool bool
	}{
		{
			name:         "UserAllowed",
			rows:         pgxmock.NewRows([]string{"allowed"}).AddRow(true),
			rowsError:    nil,
			expectedErr:  nil,
			expectedBool: true,
		},
		{
			name:         "UserForbidden",
			rows:         pgxmock.NewRows([]string{"allowed"}).AddRow(false),
			rowsError:    nil,
			expectedErr:  nil,
			expectedBool: false,
		},
		{
			name:         "Error",
			rows:         nil,
			rowsError:    errors.New("Some error"),
			expectedErr:  errors.New("Some error"),
			expectedBool: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the AccountGetUserByID query with the specified parameters
			if tc.rows != nil {
				mock.ExpectQuery(regexp.QuoteMeta(AccountGetUserByID)).
					WithArgs(accountID, userID).
					WillReturnRows(tc.rows).
					WillReturnError(tc.rowsError)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(AccountGetUserByID)).
					WithArgs(accountID, userID).
					WillReturnError(tc.rowsError)
			}

			err := repo.CheckForbidden(context.Background(), accountID, userID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_AddUserInAccount(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name        string
		execResult  pgconn.CommandTag
		execError   error
		expectedErr error
	}{
		{
			name:        "Success",
			execResult:  pgconn.CommandTag{},
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Error",
			execResult:  pgconn.CommandTag{},
			execError:   errors.New("Some error"),
			expectedErr: fmt.Errorf("[repo] can't create accountUser %s, %w", AccountUserCreate, errors.New("Some error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the AccountUserCreate query with the specified parameters
			mock.ExpectExec(regexp.QuoteMeta(AccountUserCreate)).
				WithArgs(userID, accountID).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)

			err := repo.AddUserInAccount(context.Background(), userID, accountID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_UpdateAccount(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()
	account := &models.Accounts{
		ID:             accountID,
		Balance:        100.0,
		Accumulation:   true,
		BalanceEnabled: true,
		MeanPayment:    "account",
	}

	testCases := []struct {
		name        string
		execResult  pgconn.CommandTag
		execError   error
		expectedErr error
	}{
		{
			name:        "Success",
			execResult:  pgconn.CommandTag{},
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Error",
			execResult:  pgconn.CommandTag{},
			execError:   errors.New("Some error"),
			expectedErr: fmt.Errorf("[repo] failed update account %w", errors.New("Some error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the AccountUpdate query with the specified parameters
			mock.ExpectExec(regexp.QuoteMeta(AccountUpdate)).
				WithArgs(account.Balance, account.Accumulation, account.BalanceEnabled, account.MeanPayment, account.ID).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)

			err := repo.UpdateAccount(context.Background(), userID, account)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_CheckDuplicate(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name        string
		rows        *pgxmock.Rows
		rowsError   error
		expectedErr error
	}{
		{
			name:        "NoDuplicate",
			rows:        pgxmock.NewRows([]string{}),
			rowsError:   pgx.ErrNoRows,
			expectedErr: nil,
		},
		{
			name:        "Duplicate",
			rows:        pgxmock.NewRows([]string{"user_id", "account_id"}).AddRow(userID, accountID),
			rowsError:   nil,
			expectedErr: &models.DuplicateError{},
		},
		{
			name:        "QueryError",
			rows:        pgxmock.NewRows([]string{}),
			rowsError:   errors.New("Some error"),
			expectedErr: fmt.Errorf("[repo] query error: %w", errors.New("Some error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			// Expect a call to execute the AccountGetUserByID query with the specified parameters
			mock.ExpectQuery(regexp.QuoteMeta(AccountGetUserByID)).
				WithArgs(accountID, userID).
				WillReturnRows(tc.rows).
				WillReturnError(tc.rowsError)

			err := repo.CheckDuplicate(context.Background(), userID, accountID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
