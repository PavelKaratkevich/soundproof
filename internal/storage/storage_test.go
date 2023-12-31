package storage_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/model/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserByIDSuccess(t *testing.T) {
	req := domain.LoginRequest{
		Email:    "p_korotkevitch@mail.ru",
		Password: "12345",
	}

	user := &domain.ProfileResponse{
		ID:        1,
		FirstName: "Pavel",
		LastName:  "Karatkevich",
		Email:     "p_korotkevitch@mail.ru",
		Created:   time.Now(),
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().GetUserProfile(req).Times(1).Return(user, nil)

	rec := httptest.NewRecorder()

	userResult, err := mockStorage.GetUserProfile(req)
	require.NoError(t, err)
	require.Equal(t, user, userResult)
	// log.Printf("User %v", user)
	// log.Printf("User result %v", userResult)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserByIDFailWithNegativeID(t *testing.T) {
	// Arrange
	req := domain.LoginRequest{
		Email:    "p_korotkevitch@mail.ru",
		Password: "12345",
	}

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
	mockStorage.EXPECT().GetUserProfile(req).Times(1).Return(nil, err)

	// Act
	userResult, err1 := mockStorage.GetUserProfile(req)

	// Assert
	require.Error(t, err1)
	require.NotEqual(t, user, userResult)
	require.Equal(t, err, err1)
}

func TestRegisterUserInDBSuccess(t *testing.T) {
	req := domain.UserRegistrationRequest{
		FirstName: "Pavel",
		LastName:  "Karatkevich",
		Email:     "p.korotkevitch@gmail.com",
		Password:  "12345",
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().RegisterUserInDB(req).Times(1).Return(nil)

	err := mockStorage.RegisterUserInDB(req)
	require.NoError(t, err)
}

func TestRegisterUserInDBFail(t *testing.T) {
	req := domain.UserRegistrationRequest{
		FirstName: "Pavel",
		LastName:  "Karatkevich",
		Email:     "p.korotkevitch@gmail.com",
		Password:  "12345",
	}

	err := fmt.Errorf("Error")

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().RegisterUserInDB(req).Times(1).Return(err)

	err = mockStorage.RegisterUserInDB(req)
	require.Error(t, err)
}

func TestCheckUserCredentialsSuccess(t *testing.T) {
	loginRequest := domain.LoginRequest{
		Email:    "p.korotkevitch@gmail.com",
		Password: "12345",
	}
	loginResponse := &domain.LoginResponse{
		ID:           1,
		FirstName:    "Pavel",
		LastName:     "Karatkevich",
		Email:        "p.korotkevitch@gmail.com",
		Created:      time.Now(),
		AccessToken:  "12345",
		RefreshToken: "54321",
	}

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().CheckUserCredentials(loginRequest.Email, loginRequest.Password).Times(1).Return(true, loginResponse, nil)

	ifValid, userOutput, err := mockStorage.CheckUserCredentials(loginRequest.Email, loginRequest.Password)
	require.NoError(t, err)
	require.True(t, ifValid)
	require.Equal(t, loginResponse, userOutput)
}

func TestCheckUserCredentialsFail(t *testing.T) {
	loginRequest := domain.LoginRequest{
		Email:    "p.korotkevitch@gmail.com",
		Password: "12345",
	}
	err := fmt.Errorf("Please provide valid credentials")
	require.Error(t, err)

	ctr := gomock.NewController(t)
	mockStorage := mock.NewMockStorage(ctr)
	mockStorage.EXPECT().CheckUserCredentials(loginRequest.Email, loginRequest.Password).Times(1).Return(false, nil, err)

	ifValid, userOutput, err := mockStorage.CheckUserCredentials(loginRequest.Email, loginRequest.Password)
	require.Nil(t, userOutput)
	require.Error(t, err)
	require.False(t, ifValid)
	require.Equal(t, err, err)
}
