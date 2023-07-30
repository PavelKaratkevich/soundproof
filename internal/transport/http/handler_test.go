package transport_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/model/mock"
	transport "soundproof/internal/transport/http"
	jwtauth "soundproof/internal/transport/middleware/jwt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

///////////////////////
////  RegisterUser ////
///////////////////////

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

/////////////////////////
////  GetUserByItsID ////
/////////////////////////

func TestGetUserByItsIDFailByValidator(t *testing.T) {

	response := &domain.ProfileResponse{
		ID:        1,
		FirstName: "Adam",
		LastName:  "Smith",
		Email:     "classic_theory@economics.gov.uk",
		Created:   time.Now(),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := fmt.Sprintf("/user/profile/%d", response.ID)
	log.Println(url)

	router.GET(url, handler.GetUser)

	token, err := jwtauth.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, token)

	idInput, err := json.Marshal(response.ID)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestLoginSuccess(t *testing.T) {

	loginRequest := domain.LoginRequest{
		Email:    "classic_theory@economics.gov.uk",
		Password: "12345",
	}

	user := &domain.LoginResponse{
		ID:           1,
		FirstName:    "Adam",
		LastName:     "Smith",
		Email:        "classic_theory@economics.gov.uk",
		Created:      time.Now(),
		AccessToken:  "12345",
		RefreshToken: "54321",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := "/auth/login"

	router.POST(url, handler.Login)

	service.EXPECT().CheckCredentials(gomock.Any(), loginRequest.Email, loginRequest.Password).Return(true, user, nil)

	idInput, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestLoginFailByValidation(t *testing.T) {

	loginRequest := domain.LoginRequest{
		Email:    "",
		Password: "12345",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := "/auth/login"

	router.POST(url, handler.Login)

	idInput, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestLoginFailNotFound(t *testing.T) {

	loginRequest := domain.LoginRequest{
		Email:    "classic_theory@economics.gov.uk",
		Password: "12345",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := "/auth/login"

	router.POST(url, handler.Login)

	service.EXPECT().CheckCredentials(gomock.Any(), loginRequest.Email, loginRequest.Password).Return(false, nil, sql.ErrNoRows)

	idInput, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestLoginFail_Internal_Server_Error(t *testing.T) {

	loginRequest := domain.LoginRequest{
		Email:    "classic_theory@economics.gov.uk",
		Password: "12345",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := "/auth/login"

	router.POST(url, handler.Login)

	service.EXPECT().CheckCredentials(gomock.Any(), loginRequest.Email, loginRequest.Password).Return(false, nil, sql.ErrConnDone)

	idInput, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func Test_Login_Fail_Unauthorized(t *testing.T) {

	loginRequest := domain.LoginRequest{
		Email:    "classic_theory@economics.gov.uk",
		Password: "12345",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	url := "/auth/login"

	router.POST(url, handler.Login)

	service.EXPECT().CheckCredentials(gomock.Any(), loginRequest.Email, loginRequest.Password).Return(false, nil, fmt.Errorf("wrong password"))

	idInput, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(idInput))

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}