package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"

	sessionMocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_GetSessionByCookie(t *testing.T) {
	testCase := []struct {
		expectedSession models.Session
		expectedError   error
		name            string
	}{
		{
			expectedSession: models.Session{
				UserId: uuid.New(),
				Cookie: uuid.New().String(),
			},
			expectedError: sessions.ErrSessionIsAlreadyCreated,
			name:          "Successfull getting session",
		},
		{
			expectedSession: models.Session{
				UserId: uuid.New(),
				Cookie: uuid.New().String(),
			},
			expectedError: sessions.ErrSessionNotFound,
			name:          "Session not found",
		},
		{
			expectedSession: models.Session{
				UserId: uuid.New(),
				Cookie: uuid.New().String(),
			},
			expectedError: sessions.ErrInternalServer,
			name:          "Internal error",
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authRepository := sessionMocks.NewMockRepository(ctl)
	usecase := NewSessionUsecase(authRepository)

	for _, tc := range testCase {
		authRepository.EXPECT().GetSessionByCookie(context.TODO(), "").Return(tc.expectedSession, tc.expectedError).Times(1)
		session, err := usecase.GetSessionByCookie(context.TODO(), "")

		require.Error(t, err, tc.expectedError)
		require.Equal(t, session, tc.expectedSession, tc.name)
	}
}

func Test_DeletSessionByCookie_OK(t *testing.T) {
	testCase := struct {
		expectedSession models.Session
		expectedError   error
		name            string
	}{
		expectedSession: models.Session{
			UserId: uuid.New(),
			Cookie: uuid.New().String(),
		},
		expectedError: nil,
		name:          "Successfull getting session",
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	sessionRep := sessionMocks.NewMockRepository(ctl)
	usecase := NewSessionUsecase(sessionRep)

	sessionRep.EXPECT().DeleteSession(context.TODO(), testCase.expectedSession.Cookie).Return(nil).Times(1)

	err := usecase.DeleteSessionByCookie(context.TODO(), testCase.expectedSession.Cookie)

	require.NoError(t, testCase.expectedError, err)
}

func Test_CreateSession_OK(t *testing.T) {
	testCase := []struct {
		expectedSession models.Session
		expectedError   error
		name            string
	}{
		{
			expectedSession: models.Session{
				UserId: uuid.New(),
				Cookie: uuid.New().String(),
			},
			expectedError: nil,
			name:          "Successfull creating session",
		},
		{
			expectedSession: models.Session{
				UserId: uuid.New(),
				Cookie: uuid.New().String(),
			},
			expectedError: nil,
			name:          "Successfull getting session",
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	sessionRep := sessionMocks.NewMockRepository(ctl)
	usecase := NewSessionUsecase(sessionRep)

	for _, tc := range testCase {
		sessionRep.EXPECT().CreateSession(context.TODO(), tc.expectedSession.Cookie).Return(nil).Times(1)

		_, err := usecase.CreateSessionById(context.TODO(), tc.expectedSession.UserId)

		require.NoError(t, tc.expectedError, err)
	}
}
