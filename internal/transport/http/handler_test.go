package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/model/mock"
	"soundproof/internal/domain/service"
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
		Email:     "economy@gov.uk",
		Password:  "1234567890",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock.NewMockStorage(ctrl)
	service := service.NewUserService(zap.NewNop(), storage)
	handler := transport.NewHandler(zap.NewNop(), service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	res, err := json.Marshal(req)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(res))
	require.NoError(t, err)

	// Serving the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Asserting the response code
	require.Equal(t, http.StatusOK, recorder.Code)

	// res, err := json.Marshal(req)
	// require.NoError(t, err)

	// request, err := http.NewRequest(http.MethodPost, "/auth/register/", bytes.NewReader(res))
	// require.NoError(t, err)
	// require.Equal(t, http.StatusOK, recorder.Code)

	// router := gin.Default()
	// router.ServeHTTP(recorder, request)
	// require.NoError(t, err)
}

// func TestCreateUser(t *testing.T) {

// 	user, psw := randomUser(t)

// 	var arg = gin.H{
// 		"username":  user.Username,
// 		"full_name": user.FullName,
// 		"email":     user.Email,
// 		"password":  psw,
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := mockdb.NewMockStore(ctrl)

// 	// build stubs
// 	store.EXPECT().
// 		CreateUser(gomock.Any(), gomock.Any()).
// 		Times(1).
// 		Return(user, nil)

// 	server := newTestServer(t, store)
// 	recorder := httptest.NewRecorder()

// 	res, err := json.Marshal(arg)
// 	require.NoError(t, err)

// 	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(res))
// 	require.NoError(t, err)
// 	require.Equal(t, http.StatusOK, recorder.Code)

// 	server.router.ServeHTTP(recorder, request)
// 	require.NoError(t, err)
// }
