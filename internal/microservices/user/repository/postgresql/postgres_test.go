package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		user       models.User
		errRows    error
		returnRows uuid.UUID
		expected   uuid.UUID
		err        error
	}{
		{
			name: "ValidUser",
			user: models.User{
				Login:    "testuser",
				Username: "Test User",
				Password: "password",
			},
			returnRows: uuid.New(),
			err:        nil,
		},
		{
			name: "InvalidUser",
			user: models.User{
				Login:    "",
				Username: "",
				Password: "",
			},
			errRows:    errors.New("Invalid user data"),
			returnRows: uuid.Nil,
			err:        errors.New("error request Invalid user data"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserCreate)

			mock.ExpectQuery(escapedQuery).
				WithArgs(test.user.Login, test.user.Username, test.user.Password).
				WillReturnError(test.errRows).
				WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(test.returnRows))

			userID, err := repo.CreateUser(context.Background(), test.user)
			assert.Equal(t, test.returnRows, userID)

			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.err, err)
			}
		})
	}
}

//func TestGetByID(t *testing.T) {
//	userID := uuid.New()
//	tests := []struct {
//		name     string
//		rows     *pgxmock.Rows
//		err      error
//		rowsErr  error
//		expected *models.User
//	}{
//		{
//			name: "ValidUser",
//			rows: pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}).
//				AddRow(userID, "testuser", "Test User", "password", 100.0, userID),
//			err: nil,
//			expected: &models.User{
//				ID:            userID,
//				Login:         "testuser",
//				Username:      "Test User",
//				Password:      "password",
//				PlannedBudget: 100.0,
//				AvatarURL:     userID,
//			},
//		},
//		{
//			name:     "UserNotFound",
//			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
//			rowsErr:  sql.ErrNoRows,
//			err:      fmt.Errorf("[repo] No Such user: %s doesn't exist: sql: no rows in result set", userID.String()),
//			expected: nil,
//		},
//		{
//			name:     "DatabaseError",
//			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
//			rowsErr:  errors.New("database error"),
//			err:      errors.New("failed request db SELECT id, login, username, password_hash, planned_budget, avatar_url FROM users WHERE id = $1;, database error"),
//			expected: nil,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			mock, _ := pgxmock.NewPool()
//			ctl := gomock.NewController(t)
//			defer ctl.Finish()
//
//			logger := *logger.NewLogger(context.TODO())
//			repo := NewRepository(mock, logger)
//
//			escapedQuery := regexp.QuoteMeta(UserIDGetByID)
//
//			mock.ExpectQuery(escapedQuery).
//				WithArgs(userID).
//				WillReturnRows(test.rows).
//				WillReturnError(test.rowsErr)
//
//			user, err := repo.GetByID(context.Background(), userID)
//
//			if !reflect.DeepEqual(test.expected, user) {
//				t.Errorf("Expected user: %+v, but got: %+v", test.expected, user)
//			}
//
//			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", test.err, err)
//			}
//
//			if err := mock.ExpectationsWereMet(); err != nil {
//				t.Errorf("There were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}

func TestGetUserByLogin(t *testing.T) {
	login := "testuser"
	userID := uuid.New()
	tests := []struct {
		name     string
		rows     *pgxmock.Rows
		err      error
		rowsErr  error
		expected *models.User
	}{
		{
			name: "ValidUser",
			rows: pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}).
				AddRow(userID, login, "Test User", "password", 100.0, userID),
			err: nil,
			expected: &models.User{
				ID:            userID,
				Login:         login,
				Username:      "Test User",
				Password:      "password",
				PlannedBudget: 100.0,
				AvatarURL:     userID,
			},
		},
		{
			name:     "UserNotFound",
			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
			rowsErr:  pgx.ErrNoRows,
			err:      fmt.Errorf("[repo] failed request db no rows in result set"),
			expected: nil,
		},
		{
			name:     "DatabaseError",
			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
			rowsErr:  errors.New("database error"),
			err:      fmt.Errorf("[repo] failed request db %w", errors.New("database error")),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserGetByUserName)

			mock.ExpectQuery(escapedQuery).
				WithArgs(login).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			user, err := repo.GetUserByLogin(context.Background(), login)

			if !reflect.DeepEqual(test.expected, user) {
				t.Errorf("Expected user: %+v, but got: %+v", test.expected, user)
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

func TestGetUserBalance(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		row      *pgxmock.Rows
		err      error
		rowsErr  error
		expected float64
	}{
		{
			name:     "ValidBalance",
			row:      pgxmock.NewRows([]string{"sum"}).AddRow(100.0),
			err:      nil,
			rowsErr:  nil,
			expected: 100.0,
		},
		{
			name:     "NoRows",
			row:      pgxmock.NewRows([]string{"sum"}),
			rowsErr:  sql.ErrNoRows,
			err:      fmt.Errorf("[repo] %w: %v", &models.NoSuchUserIdBalanceError{UserID: userID}, sql.ErrNoRows),
			expected: 0,
		},
		{
			name:     "DatabaseError",
			row:      pgxmock.NewRows([]string{"sum"}).AddRow(nil),
			rowsErr:  errors.New("err"),
			err:      fmt.Errorf("[repo] failed request db err"),
			expected: 0,
		},
		{
			name:     "InvalidBalance",
			row:      pgxmock.NewRows([]string{"sum"}).AddRow(nil),
			err:      nil,
			rowsErr:  nil,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(AccountBalance)

			mock.ExpectQuery(escapedQuery).
				WithArgs(userID).
				WillReturnRows(test.row).
				WillReturnError(test.rowsErr)

			balance, err := repo.GetUserBalance(context.Background(), userID)

			if balance != test.expected {
				t.Errorf("Expected balance: %f, but got: %f", test.expected, balance)
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

func TestGetPlannedBudget(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		row      *pgxmock.Rows
		err      error
		rowsErr  error
		expected float64
	}{
		{
			name:     "ValidPlannedBudget",
			row:      pgxmock.NewRows([]string{"planned_budget"}).AddRow(200.0),
			err:      nil,
			rowsErr:  nil,
			expected: 200.0,
		},
		{
			name:     "NoRows",
			row:      pgxmock.NewRows([]string{"planned_budget"}),
			rowsErr:  sql.ErrNoRows,
			err:      fmt.Errorf("[repo] %w: %v", &models.NoSuchPlannedBudgetError{UserID: userID}, sql.ErrNoRows),
			expected: 0,
		},
		{
			name:     "DatabaseError",
			row:      pgxmock.NewRows([]string{"planned_budget"}).AddRow(nil),
			rowsErr:  errors.New("err"),
			err:      fmt.Errorf("[repository] failed request db err"),
			expected: 0,
		},
		{
			name:     "InvalidPlannedBudget",
			row:      pgxmock.NewRows([]string{"planned_budget"}).AddRow(nil),
			err:      nil,
			rowsErr:  nil,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserGetPlannedBudget)

			mock.ExpectQuery(escapedQuery).
				WithArgs(userID).
				WillReturnRows(test.row).
				WillReturnError(test.rowsErr)

			plannedBudget, err := repo.GetPlannedBudget(context.Background(), userID)

			if plannedBudget != test.expected {
				t.Errorf("Expected planned budget: %f, but got: %f", test.expected, plannedBudget)
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

func TestGetCurrentBudget(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		row      *pgxmock.Rows
		err      error
		rowsErr  error
		expected float64
	}{
		{
			name:     "ValidCurrentBudget",
			row:      pgxmock.NewRows([]string{"total_sum"}).AddRow(500.0),
			err:      nil,
			rowsErr:  nil,
			expected: 500.0,
		},

		{
			name:     "NoRows",
			row:      pgxmock.NewRows([]string{"total_sum"}),
			rowsErr:  sql.ErrNoRows,
			err:      fmt.Errorf("[repository] failed request db sql: no rows in result set"),
			expected: 0,
		},
		{
			name:     "DatabaseError",
			row:      pgxmock.NewRows([]string{"total_sum"}).AddRow(nil),
			rowsErr:  errors.New("err"),
			err:      fmt.Errorf("[repository] failed request db err"),
			expected: 0,
		},
		{
			name:     "InvalidPlannedBudget",
			row:      pgxmock.NewRows([]string{"total_sum"}).AddRow(nil),
			err:      nil,
			rowsErr:  nil,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(`SELECT SUM(outcome) AS total_sum
			FROM transaction
			WHERE date_part('month', date) = date_part('month', CURRENT_DATE)
			  AND date_part('year', date) = date_part('year', CURRENT_DATE)
			AND outcome > 0
			AND account_income = account_outcome
			AND user_id = $1;`)

			mock.ExpectQuery(escapedQuery).
				WithArgs(userID).
				WillReturnRows(test.row).
				WillReturnError(test.rowsErr)

			currentBudget, err := repo.GetCurrentBudget(context.Background(), userID)

			if currentBudget != test.expected {
				t.Errorf("Expected current budget: %f, but got: %f", test.expected, currentBudget)
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

// func TestGetAccounts(t *testing.T) {
// 	userID := uuid.New()
// 	tests := []struct {
// 		name     string
// 		rows     *pgxmock.Rows
// 		err      error
// 		rowsErr  error
// 		expected []models.Accounts
// 	}{
// 		{
// 			name: "ValidAccounts",
// 			rows: pgxmock.NewRows([]string{"id", "balance", "sharing_id", "accumulation", "balance_enabled", "mean_payment"}).
// 				AddRow(userID, 100.0, userID, true, false, "Кошелек").
// 				AddRow(userID, 200.0, userID, true, false, "Наличка"),
// 			err: nil,
// 			expected: []models.Accounts{
// 				{
// 					ID:             userID,
// 					Balance:        100.0,
// 					SharingID:      userID,
// 					Accumulation:   true,
// 					BalanceEnabled: false,
// 					MeanPayment:    "Кошелек",
// 				},
// 				{
// 					ID:             userID,
// 					Balance:        200.0,
// 					SharingID:      userID,
// 					Accumulation:   true,
// 					BalanceEnabled: false,
// 					MeanPayment:    "Наличка",
// 				},
// 			},
// 		},
// 		{
// 			name: "ValidAccounts",
// 			rows: pgxmock.NewRows([]string{"id", "balance", "sharing_id", "accumulation", "balance_enabled", "mean_payment"}).
// 				AddRow("fff", 100.0, userID, true, false, "Кошелек").
// 				AddRow(userID, 200.0, userID, true, false, "Наличка"),
// 			err:      fmt.Errorf("[repo] Scanning value error for column 'id': Scan: invalid UUID length: 3"),
// 			expected: nil,
// 		},
// 		{
// 			name:     "NoAccountsFound",
// 			rows:     pgxmock.NewRows([]string{"id", "balance", "sharing_id", "accumulation", "balance_enabled", "mean_payment"}),
// 			rowsErr:  nil,
// 			err:      fmt.Errorf("[repo] No Such Accounts from user: %s doesn't exist: <nil>", userID.String()),
// 			expected: nil,
// 		},
// 		{
// 			name:     "Rows error",
// 			rows:     pgxmock.NewRows([]string{"id", "balance", "sharing_id", "accumulation", "balance_enabled", "mean_payment"}).RowError(0, errors.New("err")),
// 			rowsErr:  nil,
// 			err:      fmt.Errorf("[repo] %w", errors.New("err")),
// 			expected: nil,
// 		},
// 		{
// 			name:     "DatabaseError",
// 			rows:     pgxmock.NewRows([]string{"id", "balance", "sharing_id", "accumulation", "balance_enabled", "mean_payment"}),
// 			rowsErr:  errors.New("database error"),
// 			err:      fmt.Errorf("[repo] %w", errors.New("database error")),
// 			expected: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			mock, _ := pgxmock.NewPool()
// 			ctl := gomock.NewController(t)
// 			defer ctl.Finish()

// 			logger := *logger.NewLogger(context.TODO())
// 			repo := NewRepository(mock, logger)

// 			escapedQuery := regexp.QuoteMeta(AccountGet)

// 			mock.ExpectQuery(escapedQuery).
// 				WithArgs(userID).
// 				WillReturnRows(test.rows).
// 				WillReturnError(test.rowsErr)

// 			accounts, err := repo.GetAccounts(context.Background(), userID)

// 			if !reflect.DeepEqual(test.expected, accounts) {
// 				t.Errorf("Expected accounts: %+v, but got: %+v", test.expected, accounts)
// 			}

// 			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
// 				t.Errorf("Expected error: %v, but got: %v", test.err, err)
// 			}

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("There were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

// func TestCheckUser(t *testing.T) {
// 	userID := uuid.New()
// 	tests := []struct {
// 		name     string
// 		rows     *pgxmock.Rows
// 		err      error
// 		rowsErr  error
// 		expected error
// 	}{
// 		{
// 			name: "UserExists",
// 			rows: pgxmock.NewRows([]string{"exists"}).
// 				AddRow(true),
// 			err:      nil,
// 			expected: nil,
// 		},
// 		{
// 			name:     "UserDoesNotExist",
// 			rows:     pgxmock.NewRows([]string{"exists"}),
// 			rowsErr:  sql.ErrNoRows,
// 			err:      fmt.Errorf("failed request checkUser %w: %v", &models.NoSuchUserError{UserID: userID}, sql.ErrNoRows),
// 			expected: &models.NoSuchUserError{UserID: userID},
// 		},
// 		{
// 			name:     "DatabaseError",
// 			rows:     pgxmock.NewRows([]string{"exists"}),
// 			rowsErr:  errors.New("database error"),
// 			err:      fmt.Errorf("[repo] failed request checkUser %w", errors.New("database error")),
// 			expected: errors.New("database error"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			mock, _ := pgxmock.NewPool()
// 			ctl := gomock.NewController(t)
// 			defer ctl.Finish()

// 			logger := *logger.NewLogger(context.TODO())
// 			repo := NewRepository(mock, logger)

// 			escapedQuery := regexp.QuoteMeta(UserCheck)

// 			mock.ExpectQuery(escapedQuery).
// 				WithArgs(userID).
// 				WillReturnRows(test.rows).
// 				WillReturnError(test.rowsErr)

// 			err := repo.CheckUser(context.Background(), userID)

// 			if (test.err == nil && err != nil) || (test.err != nil && err == nil) || (test.err != nil && err != nil && test.err.Error() != err.Error()) {
// 				t.Errorf("Expected error: %v, but got: %v", test.err, err)
// 			}

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("There were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

func TestUpdateUser(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		user     models.User
		err      error
		rowsErr  error
		expected error
	}{
		{
			name: "ValidUserUpdate",
			user: models.User{
				ID:            userID,
				Username:      "Updated User",
				PlannedBudget: 1000.0,
				AvatarURL:     userID,
			},
			rowsErr:  nil,
			expected: nil,
		},
		{
			name: "InvalidUserUpdate",
			user: models.User{
				ID:            userID, // Invalid ID
				Username:      "Updated User",
				PlannedBudget: 1000.0,
				AvatarURL:     userID,
			},
			rowsErr:  errors.New("Update failed"),
			expected: errors.New("[repo] failed update user Update failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserUpdate)

			mock.ExpectExec(escapedQuery).
				WithArgs(test.user.ID, test.user.Username, test.user.PlannedBudget, test.user.AvatarURL).
				WillReturnError(test.rowsErr).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1))

			err := repo.UpdateUser(context.Background(), &test.user)

			if (test.expected == nil && err != nil) || (test.expected != nil && err == nil) || (test.expected != nil && err != nil && test.expected.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}

func TestUpdatePhoto(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name     string
		userID   uuid.UUID
		path     uuid.UUID
		errRows  error
		expected error
	}{
		{
			name:     "ValidUpdate",
			userID:   userID,
			path:     userID,
			errRows:  nil,
			expected: nil,
		},
		{
			name:     "InvalidUser",
			userID:   userID,
			path:     userID,
			errRows:  sql.ErrNoRows,
			expected: fmt.Errorf("[repo] No Such user: %s doesn't exist: sql: no rows in result set", userID.String()),
		},
		{
			name:     "InvalidDB",
			userID:   userID,
			path:     userID,
			errRows:  errors.New("err"),
			expected: fmt.Errorf("[repo] failed request db: %s, %w", UserUpdatePhoto, errors.New("err")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserUpdatePhoto)
			mock.ExpectExec(escapedQuery).
				WithArgs(test.userID, test.path).
				WillReturnError(test.errRows).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1))

			err := repo.UpdatePhoto(context.Background(), test.userID, test.path)

			if (test.expected == nil && err != nil) || (test.expected != nil && err == nil) || (test.expected != nil && err != nil && test.expected.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", test.expected, err)
			}
		})
	}
}
