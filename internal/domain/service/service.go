package service

import (
	"soundproof/internal/domain"

	"github.com/gin-gonic/gin"
)

// UserService implements interface Service 
type UserService struct {
	storage domain.Storage
}

func (s *UserService) RegisterUser(c *gin.Context, req domain.UserRegistrationRequest) (int, error) {
	res, err := s.storage.RegisterUserInDB(c, req)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func NewUserService(s domain.Storage) *UserService {
	return &UserService{
		storage: s,
	}
}