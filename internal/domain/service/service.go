package service

import (
	domain "soundproof/internal/domain/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"soundproof/pkg/eth"
)

type UserService struct {
	logger  *zap.Logger
	storage domain.Storage
}

func (s *UserService) RegisterUser(c *gin.Context, req domain.UserRegistrationRequest) error {
	err := s.storage.RegisterUserInDB(c, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUser(c *gin.Context, req domain.UpdateUserProfileRequest) error {
	
	address, err := eth.ParseMetamaskSignedString(req.SignedMessage, req.Signature)
	if err != nil {
		return err
	}

	err = s.storage.UpdateUserProfile(c, address, req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserProfile(c *gin.Context, req domain.LoginRequest) (*domain.ProfileResponse, error) {
	res, err := s.storage.GetUserProfile(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) CheckCredentials(c *gin.Context, email, password string) (bool, *domain.LoginResponse, error) {
	ifValid, user, err := s.storage.CheckUserCredentials(c, email, password)
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
