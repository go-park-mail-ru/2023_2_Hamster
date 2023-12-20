package grpc

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := uuid.New()
	expectedLogin := "testuser"
	expectedUsername := "Test User"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	mockAuthServices.EXPECT().
		SignUp(gomock.Any(), gomock.Any()).
		Return(expectedID, expectedLogin, expectedUsername, nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	signUpRequest := &proto.SignUpRequest{
		Login:    expectedLogin,
		Username: expectedUsername,
		Password: "testpassword",
	}

	signUpResponse, err := authGRPCInstance.SignUp(context.Background(), signUpRequest)

	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse)

	assert.Equal(t, "200", signUpResponse.Status)
	assert.Equal(t, expectedID.String(), signUpResponse.Body.Id)
	assert.Equal(t, expectedLogin, signUpResponse.Body.Login)
	assert.Equal(t, expectedUsername, signUpResponse.Body.Username)
}

func TestSignUp_ErrorCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"
	expectedUsername := "Test User"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := errors.New("signup failed")
	mockAuthServices.EXPECT().
		SignUp(gomock.Any(), gomock.Any()).
		Return(uuid.Nil, "", "", expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	signUpRequest := &proto.SignUpRequest{
		Login:    expectedLogin,
		Username: expectedUsername,
		Password: "testpassword",
	}

	signUpResponse, err := authGRPCInstance.SignUp(context.Background(), signUpRequest)

	assert.Error(t, err)
	assert.Nil(t, signUpResponse)

}

func TestSignUp_UserAlreadyExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"
	expectedUsername := "Test User"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := &models.UserAlreadyExistsError{}
	mockAuthServices.EXPECT().
		SignUp(gomock.Any(), gomock.Any()).
		Return(uuid.Nil, "", "", expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	signUpRequest := &proto.SignUpRequest{
		Login:    expectedLogin,
		Username: expectedUsername,
		Password: "testpassword",
	}

	signUpResponse, err := authGRPCInstance.SignUp(context.Background(), signUpRequest)

	assert.Error(t, err)
	assert.Nil(t, signUpResponse)

	var errUserAlreadyExists *models.UserAlreadyExistsError
	assert.Equal(t, "user already exists", errUserAlreadyExists.Error())
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"
	expectedUsername := "Test User"
	expectedID := uuid.New()

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	mockAuthServices.EXPECT().
		Login(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedID, expectedLogin, expectedUsername, nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	loginRequest := &proto.LoginRequest{
		Login:    expectedLogin,
		Password: "testpassword",
	}

	loginResponse, err := authGRPCInstance.Login(context.Background(), loginRequest)

	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)

	assert.Equal(t, "200", loginResponse.Status)
	assert.Equal(t, expectedID.String(), loginResponse.Body.Id)
	assert.Equal(t, expectedLogin, loginResponse.Body.Login)
	assert.Equal(t, expectedUsername, loginResponse.Body.Username)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "nonexistentuser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := &models.NoSuchUserError{}
	mockAuthServices.EXPECT().
		Login(gomock.Any(), gomock.Eq(expectedLogin), gomock.Any()).
		Return(uuid.Nil, "", "", expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	loginRequest := &proto.LoginRequest{
		Login:    expectedLogin,
		Password: "testpassword",
	}

	loginResponse, err := authGRPCInstance.Login(context.Background(), loginRequest)

	assert.Error(t, err)
	assert.Nil(t, loginResponse)

	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestLogin_IncorrectPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := &models.IncorrectPasswordError{}
	mockAuthServices.EXPECT().
		Login(gomock.Any(), gomock.Eq(expectedLogin), gomock.Any()).
		Return(uuid.Nil, "", "", expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	loginRequest := &proto.LoginRequest{
		Login:    expectedLogin,
		Password: "wrongpassword",
	}

	loginResponse, err := authGRPCInstance.Login(context.Background(), loginRequest)

	assert.Error(t, err)
	assert.Nil(t, loginResponse)

	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	assert.Contains(t, err.Error(), "incorrect password")
}

func TestLogin_SomeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	// Simulate an IncorrectPasswordError during Login
	expectedError := errors.New("some")
	mockAuthServices.EXPECT().
		Login(gomock.Any(), gomock.Eq(expectedLogin), gomock.Any()).
		Return(uuid.Nil, "", "", expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	loginRequest := &proto.LoginRequest{
		Login:    expectedLogin,
		Password: "wrongpassword",
	}

	loginResponse, err := authGRPCInstance.Login(context.Background(), loginRequest)

	assert.Error(t, err)
	assert.Nil(t, loginResponse)

	assert.Contains(t, err.Error(), "some")
}

func TestCheckLoginUnique_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "uniqueuser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	mockAuthServices.EXPECT().
		CheckLoginUnique(gomock.Any(), gomock.Eq(expectedLogin)).
		Return(true, nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	uniqCheckRequest := &proto.UniqCheckRequest{
		Login: expectedLogin,
	}

	uniqCheckResponse, err := authGRPCInstance.CheckLoginUnique(context.Background(), uniqCheckRequest)

	assert.NoError(t, err)
	assert.NotNil(t, uniqCheckResponse)

	assert.Equal(t, "200", uniqCheckResponse.Status)
	assert.False(t, uniqCheckResponse.Body)
}

func TestCheckLoginUnique_NotUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "existinguser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	// Simulate a non-unique login during CheckLoginUnique
	mockAuthServices.EXPECT().
		CheckLoginUnique(gomock.Any(), gomock.Eq(expectedLogin)).
		Return(false, nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	uniqCheckRequest := &proto.UniqCheckRequest{
		Login: expectedLogin,
	}

	uniqCheckResponse, err := authGRPCInstance.CheckLoginUnique(context.Background(), uniqCheckRequest)

	assert.NoError(t, err)
	assert.NotNil(t, uniqCheckResponse)

	assert.Equal(t, "200", uniqCheckResponse.Status)
	assert.True(t, uniqCheckResponse.Body)
}

func TestCheckLoginUnique_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := errors.New("check login unique failed")
	mockAuthServices.EXPECT().
		CheckLoginUnique(gomock.Any(), gomock.Eq(expectedLogin)).
		Return(false, expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	uniqCheckRequest := &proto.UniqCheckRequest{
		Login: expectedLogin,
	}

	uniqCheckResponse, err := authGRPCInstance.CheckLoginUnique(context.Background(), uniqCheckRequest)

	assert.Error(t, err)
	assert.Nil(t, uniqCheckResponse)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "check login unique failed")
}

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := uuid.New()
	expectedLogin := "testuser"
	expectedUsername := "Test User"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	// Simulate a successful GetByID
	mockAuthServices.EXPECT().
		GetByID(gomock.Any(), gomock.Eq(expectedID)).
		Return(&models.User{ID: expectedID, Login: expectedLogin, Username: expectedUsername}, nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	userIDRequest := &proto.UserIdRequest{
		Id: expectedID.String(),
	}

	userResponse, err := authGRPCInstance.GetByID(context.Background(), userIDRequest)

	assert.NoError(t, err)
	assert.NotNil(t, userResponse)

	assert.Equal(t, "200", userResponse.Status)
	assert.Equal(t, expectedID.String(), userResponse.Body.Id)
	assert.Equal(t, expectedLogin, userResponse.Body.Login)
	assert.Equal(t, expectedUsername, userResponse.Body.Username)
}

func TestGetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidID := "invalid-id"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	userIDRequest := &proto.UserIdRequest{
		Id: invalidID,
	}

	userResponse, err := authGRPCInstance.GetByID(context.Background(), userIDRequest)

	assert.Error(t, err)
	assert.Nil(t, userResponse)

}

func TestGetByID_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := uuid.New()

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := &models.NoSuchUserError{}
	mockAuthServices.EXPECT().
		GetByID(gomock.Any(), gomock.Eq(expectedID)).
		Return(nil, expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	userIDRequest := &proto.UserIdRequest{
		Id: expectedID.String(),
	}

	userResponse, err := authGRPCInstance.GetByID(context.Background(), userIDRequest)

	assert.Error(t, err)
	assert.Nil(t, userResponse)

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := uuid.New()

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := errors.New("get user by ID failed")
	mockAuthServices.EXPECT().
		GetByID(gomock.Any(), gomock.Eq(expectedID)).
		Return(nil, expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	userIDRequest := &proto.UserIdRequest{
		Id: expectedID.String(),
	}

	userResponse, err := authGRPCInstance.GetByID(context.Background(), userIDRequest)

	assert.Error(t, err)
	assert.Nil(t, userResponse)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "get user by ID failed")
}

func TestChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	mockAuthServices.EXPECT().
		ChangePassword(gomock.Any(), auth.ChangePasswordInput{
			Login:       expectedLogin,
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}).
		Return(nil)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	changePasswordRequest := &proto.ChangePasswordRequest{
		Login:       expectedLogin,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	emptyResponse, err := authGRPCInstance.ChangePassword(context.Background(), changePasswordRequest)

	assert.NoError(t, err)
	assert.NotNil(t, emptyResponse)
}

func TestChangePassword_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedLogin := "testuser"
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	mockAuthServices := mocks.NewMockUsecase(ctrl)

	expectedError := errors.New("change password failed")
	mockAuthServices.EXPECT().
		ChangePassword(gomock.Any(), auth.ChangePasswordInput{
			Login:       expectedLogin,
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}).
		Return(expectedError)

	authGRPCInstance := NewAuthGRPC(mockAuthServices, *logger.NewLogger(context.TODO()))

	changePasswordRequest := &proto.ChangePasswordRequest{
		Login:       expectedLogin,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	emptyResponse, err := authGRPCInstance.ChangePassword(context.Background(), changePasswordRequest)

	assert.Error(t, err)
	assert.Nil(t, emptyResponse)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "change password failed")
}
