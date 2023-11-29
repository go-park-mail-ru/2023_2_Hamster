package postgres

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
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func Test_CreateTag(t *testing.T) {
	testUUID := uuid.New()
	testCases := []struct {
		name      string
		tag       models.Category
		returnID  uuid.UUID
		returnErr error
		expected  uuid.UUID
		err       error
	}{
		{
			name:      "Success",
			tag:       models.Category{UserID: uuid.New(), Name: "TestTag", ShowIncome: true, ShowOutcome: true, Regular: true},
			returnID:  testUUID,
			expected:  testUUID,
			returnErr: nil,
			err:       nil,
		},
		{
			name:      "Failure",
			tag:       models.Category{UserID: uuid.New(), Name: "TestTag", ShowIncome: true, ShowOutcome: true, Regular: true},
			returnID:  uuid.Nil,
			expected:  uuid.Nil,
			returnErr: errors.New("some database error"),
			err:       fmt.Errorf("[repo] failed create new tag: %w", errors.New("some database error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(CategoryCreate)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.tag.UserID, nil, tc.tag.Name, tc.tag.ShowIncome, tc.tag.ShowOutcome, tc.tag.Regular).
				WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(tc.returnID)).
				WillReturnError(tc.returnErr)

			id, err := repo.CreateTag(context.Background(), tc.tag)

			if (tc.err == nil && err != nil) || (tc.err != nil && err == nil) || (tc.err != nil && err != nil && tc.err.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.err, err)
			}

			if tc.expected != id {
				t.Errorf("Expected ID: %v, got: %v", tc.expected, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_UpdateTag(t *testing.T) {
	testCases := []struct {
		name        string
		tag         *models.Category
		rowExists   bool
		execError   error
		expectedErr error
		sqlNoRows   error
	}{
		{
			name:        "Success",
			tag:         &models.Category{ID: uuid.New(), ParentID: uuid.Nil, Name: "UpdatedTag", ShowIncome: true, ShowOutcome: true, Regular: true},
			rowExists:   true,
			execError:   nil,
			expectedErr: nil,
		},
		{
			name:        "UpdateError",
			tag:         &models.Category{ID: uuid.New(), ParentID: uuid.Nil, Name: "UpdatedTag", ShowIncome: true, ShowOutcome: true, Regular: true},
			rowExists:   true,
			execError:   errors.New("some database error"),
			expectedErr: fmt.Errorf("[repo] failed to update category info: %s, %w", CategoryUpdate, errors.New("some database error")),
			sqlNoRows:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)
			/*
				escapedQuery := regexp.QuoteMeta(CategoryGet)
				mock.ExpectQuery(escapedQuery).
					WithArgs(tc.tag.ID).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(tc.rowExists)).
					WillReturnError(tc.sqlNoRows)
			*/

			escapedQuery := regexp.QuoteMeta(CategoryUpdate)
			mock.ExpectExec(escapedQuery).
				WithArgs(nil, tc.tag.Name, tc.tag.ShowIncome, tc.tag.ShowOutcome, tc.tag.Regular, tc.tag.ID).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(tc.execError)

			err := repo.UpdateTag(context.Background(), tc.tag)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
	 func Test_DeleteTag(t *testing.T) {
		testCases := []struct {
			name        string
			tagID       uuid.UUID
			rowExists   bool
			execError   error
			expectedErr error
			sqlNoRows   error
			transaction bool
		}{
			{
				name:        "Success",
				tagID:       uuid.New(),
				rowExists:   true,
				execError:   nil,
				expectedErr: nil,
				transaction: true,
			},
			{
				name:        "TagNotFound",
				tagID:       uuid.New(),
				rowExists:   false,
				execError:   nil,
				expectedErr: fmt.Errorf("[repo] tag doesn't exist Error: %v", sql.ErrNoRows),
				sqlNoRows:   sql.ErrNoRows,
				transaction: false,
			},
			{
				name:        "DeleteError",
				tagID:       uuid.New(),
				rowExists:   true,
				execError:   errors.New("some database error"),
				expectedErr: fmt.Errorf("[repo] failed to delete category %s, %w", CategoryDelete, errors.New("some database error")),
				transaction: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mock, _ := pgxmock.NewPool()

				logger := *logger.NewLogger(context.TODO())
				repo := NewRepository(mock, logger)

					escapedQuery := regexp.QuoteMeta(CategoryGet)
					mock.ExpectQuery(escapedQuery).
						WithArgs(tc.tagID).
						WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(tc.rowExists)).
						WillReturnError(tc.sqlNoRows)

				mock.ExpectBegin()
				if tc.transaction {
					mock.ExpectCommit()
				} else {
					mock.ExpectRollback()
				}

				escapedQuery := regexp.QuoteMeta(CategoryDelete)
				mock.ExpectExec(escapedQuery).
					WithArgs(tc.tagID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)).
					WillReturnError(tc.execError)

				err := repo.DeleteTag(context.Background(), tc.tagID)
				if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
					t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
				}

				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("There were unfulfilled expectations: %s", err)
				}
			})
		}
	}
*/
func Test_GetTags(t *testing.T) {
	tagId := uuid.New()
	userId := uuid.New()
	parentId := uuid.New()
	testCases := []struct {
		name        string
		userID      uuid.UUID
		rows        *pgxmock.Rows
		rowsError   error
		expected    []models.Category
		expectedErr error
	}{
		{
			name:   "Success",
			userID: uuid.New(),
			rows: pgxmock.NewRows([]string{"id", "user_id", "parent_id", "name", "show_income", "show_outcome", "regular"}).
				AddRow(tagId, userId, parentId, "Tag1", true, true, true).
				AddRow(tagId, userId, parentId, "Tag2", false, true, false),
			rowsError: nil,
			expected: []models.Category{{
				ID:          tagId,
				UserID:      userId,
				ParentID:    parentId,
				Name:        "Tag1",
				ShowIncome:  true,
				ShowOutcome: true,
				Regular:     true,
			}, {
				ID:          tagId,
				UserID:      userId,
				ParentID:    parentId,
				Name:        "Tag2",
				ShowIncome:  false,
				ShowOutcome: true,
				Regular:     false,
			},
			},
			expectedErr: nil,
		},
		{
			name:        "NoTagsFound",
			userID:      uuid.New(),
			rows:        pgxmock.NewRows([]string{}),
			rowsError:   sql.ErrNoRows,
			expected:    nil,
			expectedErr: fmt.Errorf("[repo] Error no tags found: %v", sql.ErrNoRows),
		},
		{
			name:        "RowsError",
			userID:      uuid.New(),
			rows:        pgxmock.NewRows([]string{}).RowError(0, errors.New("SOME ERROR")),
			rowsError:   nil,
			expected:    nil,
			expectedErr: fmt.Errorf("[repo] Error rows error: %v", errors.New("SOME ERROR")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(CategoeyAll)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.userID).
				WillReturnRows(tc.rows).
				WillReturnError(tc.rowsError)

			result, err := repo.GetTags(context.Background(), tc.userID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("Expected tags: %v, but got: %v", tc.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_CheckNameUniq(t *testing.T) {
	testCases := []struct {
		name        string
		userID      uuid.UUID
		parentID    uuid.UUID
		nameToCheck string
		rowExists   bool
		rowError    error
		expected    bool
		expectedErr error
	}{
		{
			name:        "SuccessUniq",
			userID:      uuid.New(),
			parentID:    uuid.New(),
			nameToCheck: "NewTagName",
			rowExists:   false,
			rowError:    nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "NotUniq",
			userID:      uuid.New(),
			parentID:    uuid.New(),
			nameToCheck: "ExistingTagName",
			rowExists:   true,
			rowError:    nil,
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "RowsError",
			userID:      uuid.New(),
			parentID:    uuid.New(),
			nameToCheck: "TestName",
			rowExists:   false,
			rowError:    errors.New("some database error"),
			expected:    false,
			expectedErr: fmt.Errorf("[repo] Error: %v", errors.New("some database error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(CategoryNameCheck)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.userID, tc.parentID, tc.nameToCheck).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(tc.rowExists)).
				WillReturnError(tc.rowError)

			result, err := repo.CheckNameUniq(context.Background(), tc.userID, tc.parentID, tc.nameToCheck)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if result != tc.expected {
				t.Errorf("Expected result: %v, but got: %v", tc.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_CheckExist(t *testing.T) {
	testCases := []struct {
		name        string
		userID      uuid.UUID
		tagID       uuid.UUID
		rowExists   bool
		rowError    error
		expected    bool
		expectedErr error
	}{
		{
			name:        "SuccessExist",
			userID:      uuid.New(),
			tagID:       uuid.New(),
			rowExists:   true,
			rowError:    nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "NotExist",
			userID:      uuid.New(),
			tagID:       uuid.New(),
			rowExists:   false,
			rowError:    nil,
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "RowsError",
			userID:      uuid.New(),
			tagID:       uuid.New(),
			rowExists:   true,
			rowError:    errors.New("some database error"),
			expected:    false,
			expectedErr: fmt.Errorf("[repo] Error: %v", errors.New("some database error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			logger := *logger.NewLogger(context.TODO())
			repo := NewRepository(mock, logger)

			escapedQuery := regexp.QuoteMeta(CategoryExistCheck)
			mock.ExpectQuery(escapedQuery).
				WithArgs(tc.userID, tc.tagID).
				WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(tc.rowExists)).
				WillReturnError(tc.rowError)

			result, err := repo.CheckExist(context.Background(), tc.userID, tc.tagID)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if result != tc.expected {
				t.Errorf("Expected result: %v, but got: %v", tc.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
