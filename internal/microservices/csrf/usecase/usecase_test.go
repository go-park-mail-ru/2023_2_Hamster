package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecaseToken_GenerateAndCheckCSRFToken(t *testing.T) {
	u := NewUsecase(*logger.NewLogger(context.TODO()))

	expectedUserID := uuid.New()

	token, err := u.GenerateCSRFToken(expectedUserID)
	assert.NoError(t, err)

	gotUserID, err := u.CheckCSRFToken(token)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, gotUserID)
}
