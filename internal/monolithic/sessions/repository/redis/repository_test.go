package redis

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedis_GetSession_Success(t *testing.T) {
	userID := uuid.New()
	cookie := uuid.New().String()
	expectedSession := models.Session{
		UserId: userID,
		Cookie: cookie,
	}

	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectGet(cookie).SetVal(userID.String())

	session, err := repo.GetSessionByCookie(context.TODO(), cookie)

	require.NoError(t, err)
	require.Equal(t, expectedSession, session)
}

func TestRedis_GetSession_NotFoundFalier(t *testing.T) {
	cookie := uuid.New().String()
	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectGet(cookie).SetErr(redis.Nil)

	_, err := repo.GetSessionByCookie(context.TODO(), cookie)

	require.Error(t, err, sessions.ErrSessionNotFound)
}

func TestRedis_GetSession_ParseUUID_Falier(t *testing.T) {
	userID := "invalid-userid"
	cookie := uuid.New().String()

	// cookie := "invalid-uuid"
	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectGet(cookie).SetVal(userID)

	_, err := repo.GetSessionByCookie(context.TODO(), cookie)

	require.ErrorIs(t, err, sessions.ErrInvalidUUID)
}

func TestRedis_DeleteSession_Success(t *testing.T) {
	cookie := uuid.New().String()
	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectDel(cookie).SetVal(1)

	err := repo.DeleteSession(context.TODO(), cookie)

	require.NoError(t, err)
}

func TestRedis_DeleteSession_NotFoundFalier(t *testing.T) {
	cookie := uuid.New().String()
	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectDel(cookie).SetErr(redis.Nil)

	err := repo.DeleteSession(context.TODO(), cookie)

	require.Error(t, err, sessions.ErrSessionNotFound)
}

func TestRedis_CreateSession_Success(t *testing.T) {
	session := models.Session{
		UserId: uuid.New(),
		Cookie: uuid.New().String(),
	}

	mockedClient, mock := redismock.NewClientMock()
	repo := NewSessionRepository(mockedClient)

	mock.ExpectSet(session.Cookie, session.UserId.String(), 0).SetVal(session.Cookie)

	err := repo.CreateSession(context.TODO(), session)

	require.NoError(t, err)
}
