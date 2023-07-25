package storage_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/model/mock"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserByIDSuccess(t *testing.T) {
	id := 1
	user := &domain.ProfileResponse{
		ID:        1,
		FirstName: "Pavel",
		LastName:  "Karatkevich",
		Email:     "p_korotkevitch@mail.ru",
		Created:   time.Now(),
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().GetUserByID(&gin.Context{}, id).Times(1).Return(user, nil)

	rec := httptest.NewRecorder()

	userResult, err := mockStorage.GetUserByID(&gin.Context{}, id)
	require.NoError(t, err)
	require.Equal(t, user, userResult)
	log.Printf("User %v", user)
	log.Printf("User result %v", userResult)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserByIDFailWithNegativeID(t *testing.T) {

	// Arrange
	id := -1
	user := &domain.ProfileResponse{
		ID:        1,
		FirstName: "Pavel",
		LastName:  "Karatkevich",
		Email:     "p_korotkevitch@mail.ru",
		Created:   time.Now(),
	}
	err := fmt.Errorf("Key: '' Error:Field validation for '' failed on the 'number' tag")

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().GetUserByID(gomock.Any(), id).Times(1).Return(nil, err)

	// Act
	userResult, err1 := mockStorage.GetUserByID(&gin.Context{}, id)
	
	// Assert
	require.Error(t, err1)
	require.NotEqual(t, user, userResult)
	require.Equal(t, err, err1)
}

func TestRegisterUserInDBSuccess(t *testing.T) {
	req := domain.UserRegistrationRequest{
		FirstName: "Pavel",
		LastName: "Karatkevich",
		Email: "p.korotkevitch@gmail.com",
		Password: "12345",
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().RegisterUserInDB(gomock.Any(), req).Times(1).Return(1, nil)

	userResult, err := mockStorage.RegisterUserInDB(&gin.Context{}, req)
	require.NoError(t, err)
	require.Equal(t, 1, userResult)
}

func TestRegisterUserInDBFail(t *testing.T) {
	req := domain.UserRegistrationRequest{
		FirstName: "Pavel",
		LastName: "Karatkevich",
		Email: "p.korotkevitch@gmail.com",
		Password: "12345",
	}

	err := fmt.Errorf("Error")

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().RegisterUserInDB(gomock.Any(), req).Times(1).Return(0, err)

	_, err1 := mockStorage.RegisterUserInDB(&gin.Context{}, req)
	require.Error(t, err)
	require.Equal(t, err1, err)
}

func TestCheckUserCredentialsSuccess(t *testing.T) {
	loginRequest := domain.LoginRequest{
		Email: "p.korotkevitch@gmail.com",
		Password: "12345",
	}
	loginResponse := &domain.LoginResponse	{
		ID: 1,
		FirstName: "Pavel",
		LastName: "Karatkevich",
		Email: "p.korotkevitch@gmail.com",
		Created: time.Now(),
		AccessToken: "12345",
		RefreshToken: "54321",
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().CheckUserCredentials(gomock.Any(), loginRequest).Times(1).Return(true, loginResponse, nil)

	ifValid, userOutput, err := mockStorage.CheckUserCredentials(&gin.Context{}, loginRequest)
	require.NoError(t, err)
	require.True(t, ifValid)
	require.Equal(t, loginResponse, userOutput)
}

func TestCheckUserCredentialsFail(t *testing.T) {
	loginRequest := domain.LoginRequest{
		Email: "p.korotkevitch@gmail.com",
		Password: "12345",
	}
	error := fmt.Errorf("Please provide valid credentials")

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().CheckUserCredentials(gomock.Any(), loginRequest).Times(1).Return(false, nil, error)

	ifValid, userOutput, err := mockStorage.CheckUserCredentials(&gin.Context{}, loginRequest)
	require.Nil(t, userOutput)
	require.Error(t, err)
	require.False(t, ifValid)
	require.Equal(t, err, error)
}
