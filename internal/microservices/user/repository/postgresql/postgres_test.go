package postgresql

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	logger := logger.GetLogger()
	repo := NewRepository(mock, logger)

	user := models.User{
		Login:    "testuser",
		Username: "Test User",
		Password: "password",
	}

	id := uuid.New()

	// Expect the query with the exact SQL string
	escapedQuery := regexp.QuoteMeta(UserCreate)
	mock.ExpectQuery(escapedQuery).
		WithArgs(user.Login, user.Username, user.Password).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(id))

	result, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, id, result)
}
