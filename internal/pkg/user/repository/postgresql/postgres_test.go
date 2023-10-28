package postgresql

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockConn := pgxpoolmock.NewMockPgxPool(ctl)
	logger := logger.CreateCustomLogger()

	repo := NewRepository(mockConn, *logger)

	user := models.User{
		Login:    "testuser",
		Username: "Test User",
		Password: "password",
	}

	id := uuid.New()
	columns := []string{"id"}
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(id.String())
	mockConn.EXPECT().QueryRow(gomock.Any(), UserCreate, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows)

	result, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, id, result)
}
