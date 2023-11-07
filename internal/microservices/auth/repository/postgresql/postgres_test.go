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
	"github.com/pashagolub/pgxmock"
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
