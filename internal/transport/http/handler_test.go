package transport_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/model/mock"
	transport "soundproof/internal/transport/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRegisterUserSuccess(t *testing.T) {

	req := domain.UserRegistrationRequest{
		FirstName: "Adam",
		LastName:  "Smith",
		Email:     "classic_theory@economics.gov.uk",
		Password:  "1234567890",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	res, err := json.Marshal(req)
	require.NoError(t, err)

	service.EXPECT().RegisterUser(gomock.Any(), req).Times(1).Return(nil)

	request, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(res))
	require.NoError(t, err)

	log.Println(request)

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestRegisterUserFail(t *testing.T) {

	req := domain.UserRegistrationRequest{
		FirstName: "Adam",
		LastName:  "Smith",
		Email:     "classic_theory@.economics.gov.uk",
		Password:  "1234567890",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	res, err := json.Marshal(req)
	require.NoError(t, err)

	request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(res))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestRegisterUserInternalError(t *testing.T) {

	req := domain.UserRegistrationRequest{
		FirstName: "Adam",
		LastName:  "Smith",
		Email:     "classic_theory@economics.gov.uk",
		Password:  "1234567890",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	res, err := json.Marshal(req)
	require.NoError(t, err)

	service.EXPECT().RegisterUser(gomock.Any(), req).Times(1).Return(fmt.Errorf("error"))

	request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(res))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
