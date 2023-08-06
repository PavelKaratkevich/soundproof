package service

import (
	domain "soundproof/internal/domain/model"
	"soundproof/pkg/eth"

	"go.uber.org/zap"
)

type UserService struct {
	logger  *zap.Logger
	storage domain.Storage
}

func (s *UserService) RegisterUser(req domain.UserRegistrationRequest) error {
	err := s.storage.RegisterUserInDB(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUser(req domain.UpdateUserProfileRequest) error {
	address, err := eth.ParseMetamaskSignedString(req.SignedMessage, req.Signature)
	if err != nil {
		return err
	}

	err = s.storage.UpdateUserProfile(address, req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserProfile(req domain.LoginRequest) (*domain.ProfileResponse, error) {
	res, err := s.storage.GetUserProfile(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) CheckCredentials(email, password string) (bool, *domain.LoginResponse, error) {
	ifValid, user, err := s.storage.CheckUserCredentials(email, password)
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
