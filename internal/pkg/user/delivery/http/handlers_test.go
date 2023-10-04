package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create an instance of your mock UserService and Logger
	mockUserService := mocks.NewMockUsecase(ctrl)
	mockLogger := commonHttp.NewMockLogger()

	// Create an instance of your handler with the mock dependencies
	handler := NewHandler(mockUserService, mockLogger)
	u := uuid.New()
	uInvalid := "dsfaf"
	// Define test cases
	testCases := []struct {
		name           string
		userID         uuid.UUID
		userServiceErr error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			userID:         u,
			userServiceErr: nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Balance": 100.0}`, // Replace with your expected JSON response
		},
		{
			name:           "Invalid UserID",
			userID:         "invalid-user-id",
			userServiceErr: commonHttp.NewBadRequestError("Invalid user ID"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid user ID"}`, // Replace with your expected error response
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock behavior for the GetUserBalance method
			mockUserService.EXPECT().GetUserBalance(tc.userID).Return(100.0, tc.userServiceErr)

			// Create a request with the userID as a URL parameter
			req, err := http.NewRequest("GET", "/user/balance?userId="+tc.userID.String(), nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the handler's GetUserBalance method
			handler.GetUserBalance(rr, req)

			// Check the HTTP status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check the response body
			assert.Equal(t, tc.expectedBody, rr.Body.String())
		})
	}
}
