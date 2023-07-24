package service

import (
	domain "soundproof/internal/domain/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserService struct {
	logger  *zap.Logger
	storage domain.Storage
}

func (s *UserService) RegisterUser(c *gin.Context, req domain.UserRegistrationRequest) (int, error) {
	res, err := s.storage.RegisterUserInDB(c, req)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *UserService) GetByID(c *gin.Context, id int) (*domain.ProfileResponse, error) {
	res, err := s.storage.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) CheckCredentials(c *gin.Context, req domain.LoginRequest) (bool, *domain.LoginResponse, error) {
	ifValid, user, err := s.storage.CheckUserCredentials(c, req)
	if err != nil {
		return false, nil, err
	}
	return ifValid, user, nil
}

func NewUserService(logger *zap.Logger, s domain.Storage) *UserService {
	return &UserService{
		storage: s,
		logger:  logger,
	}
}
