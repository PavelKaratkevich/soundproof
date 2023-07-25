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

func TestGetUserByIDFail(t *testing.T) {

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