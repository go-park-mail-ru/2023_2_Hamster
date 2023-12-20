package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func Test_CheckLoginUniq(t *testing.T) {
	tastCase := []struct {
		name    string
		login   string
		errRows error

		returnRows *pgxmock.Rows
		expected   bool
		err        error
	}{
		{
			name:       "Success uniqness",
			login:      "sema",
			returnRows: pgxmock.NewRows([]string{"count"}).AddRow(0),
			expected:   true,
			errRows:    nil,
			err:        nil,
		},
		{
			name:       "Falier uniqness",
			login:      "grisha",
			returnRows: pgxmock.NewRows([]string{"count"}),
			expected:   false,
			errRows:    sql.ErrNoRows,
			err:        fmt.Errorf("[repo] failed login unique check %w", sql.ErrNoRows),
		},
	}

	for _, tc := range tastCase {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserCheckLoginUnique)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.login).
				WillReturnRows(tc.returnRows).
				WillReturnError(tc.errRows)

			_, err := repo.CheckLoginUnique(context.Background(), tc.login)

			if (tc.err == nil && err != nil) || (tc.err != nil && err == nil) || (tc.err != nil && err != nil && tc.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.err, err)
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
			rowsErr:  pgx.ErrNoRows,
			err:      fmt.Errorf("[repo] %w, %v", &models.NoSuchUserError{}, pgx.ErrNoRows),
			expected: nil,
		},
		{
			name:     "DatabaseError",
			rows:     pgxmock.NewRows([]string{"id", "login", "username", "password", "planned_budget", "avatar_url"}),
			rowsErr:  errors.New("database error"),
			err:      fmt.Errorf("[repo] failed request db: %w", errors.New("database error")),
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
			err:        fmt.Errorf("(Repo) failed to scan from query: %w", errors.New("Invalid user data")), // errors.New("error request Invalid user data"),
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

func TestGetByID(t *testing.T) {
	staticUUID := uuid.New()
	tests := []struct {
		name        string
		userID      uuid.UUID
		rows        *pgxmock.Rows
		expected    *models.User
		expectedErr error
		err         error
	}{
		{
			name:   "ValidUser",
			userID: staticUUID,
			rows: pgxmock.NewRows([]string{
				"id", "login", "username", "password", "planned_budget", "avatar_url",
			}).AddRow(
				staticUUID, "testuser", "Test User", "password", 1000.0, staticUUID,
			),
			expected: &models.User{
				ID:            staticUUID,
				Login:         "testuser",
				Username:      "Test User",
				Password:      "password",
				PlannedBudget: 1000.0,
				AvatarURL:     staticUUID,
			},
			expectedErr: nil,
			err:         nil,
		},
		{
			name:        "NoSuchUser",
			userID:      staticUUID,
			rows:        pgxmock.NewRows([]string{}).RowError(0, sql.ErrNoRows),
			expected:    nil,
			expectedErr: fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: staticUUID}, sql.ErrNoRows),
			err:         sql.ErrNoRows,
		},

		{
			name:        "DatabaseError",
			userID:      uuid.New(),
			rows:        pgxmock.NewRows([]string{}).RowError(0, errors.New("database error")),
			expected:    nil,
			expectedErr: fmt.Errorf("failed request db %s, %w", UserIDGetByID, errors.New("database error")),
			err:         errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserIDGetByID)

			mock.ExpectQuery(escapedQuery).
				WithArgs(staticUUID).
				WillReturnRows(test.rows).
				WillReturnError(test.err)

			user, err := repo.GetByID(context.Background(), staticUUID)

			assert.Equal(t, test.expected, user)
			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	staticUUID := uuid.New()
	newPassword := "newpassword"

	tests := []struct {
		name        string
		userID      uuid.UUID
		newPassword string
		execError   error
		expectedErr error
	}{
		{
			name:        "Success",
			userID:      staticUUID,
			newPassword: newPassword,
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "DatabaseError",
			userID:      staticUUID,
			newPassword: newPassword,
			execError:   errors.New("database error"),
			expectedErr: fmt.Errorf("[repo] failed to update password: %w", errors.New("database error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(UserChangePassword)

			mock.ExpectExec(escapedQuery).
				WithArgs(newPassword, test.userID).
				WillReturnResult(pgxmock.NewResult("0", 1)).
				WillReturnError(test.execError)

			err := repo.ChangePassword(context.Background(), test.userID, test.newPassword)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err)
			}
		})
	}
}
