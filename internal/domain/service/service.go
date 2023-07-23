package service

import (
	"soundproof/internal/domain"

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

func NewUserService(logger  *zap.Logger, s domain.Storage) *UserService {
	return &UserService{
		storage: s,
		logger: logger,
	}
}
