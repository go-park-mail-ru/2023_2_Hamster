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

			logger := *logger.InitLogger()
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

func TestGetByID(t *testing.T) {
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
				AddRow(userID, "testuser", "Test User", "password", 100.0, userID),
			err: nil,
			expected: &models.User{
				ID:            userID,
				Login:         "testuser",
				Username:      "Test User",
				Password:      "password",
				PlannedBudget: 100.0,
				AvatarURL:     userID,
			},
		},
		{
			name:     "UserNotFound",
			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
			rowsErr:  sql.ErrNoRows,
			err:      fmt.Errorf("[repo] No Such user: %s doesn't exist: sql: no rows in result set", userID.String()),
			expected: nil,
		},
		{
			name:     "DatabaseError",
			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
			rowsErr:  errors.New("database error"),
			err:      errors.New("failed request db SELECT * FROM users WHERE id = $1;, database error"),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.InitLogger()
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserIDGetByID)

			mock.ExpectQuery(escapedQuery).
				WithArgs(userID).
				WillReturnRows(test.rows).
				WillReturnError(test.rowsErr)

			user, err := repo.GetByID(context.Background(), userID)

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
			rowsErr:  sql.ErrNoRows,
			err:      fmt.Errorf("[repo] nothing found for this request %w", sql.ErrNoRows),
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

			logger := *logger.InitLogger()
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
			err:      fmt.Errorf("[repo] invalid type balance"),
			rowsErr:  nil,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.InitLogger()
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
			err:      fmt.Errorf("[repo] invalid planned budget"),
			rowsErr:  nil,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.InitLogger()
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
