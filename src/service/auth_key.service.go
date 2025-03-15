package service

import (
	"errors"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/utils"
)

type AuthKeyService interface {
	CreateAuthKey(name string, username string) (string, error)
	ValidateAuthKey(key string) (bool, error)
}

type authKeyService struct {
	authKeyRepository repository.AuthKeyRepository
}

func NewAuthKeyService(authKeyRepository repository.AuthKeyRepository) AuthKeyService {
	return &authKeyService{authKeyRepository: authKeyRepository}
}

func (s *authKeyService) CreateAuthKey(name string, username string) (string, error) {
	authKey := &model.AuthKey{
		Name:      name,
		Key:       utils.GenerateAuthKey(),
		CreatedBy: username,
	}
	authKey, err := s.authKeyRepository.Insert(authKey)
	if err != nil {
		return "", err
	}
	return authKey.Key, nil
}

func (s *authKeyService) ValidateAuthKey(key string) (bool, error) {
	authKey, err := s.authKeyRepository.GetById(key)
	if err != nil {
		return false, err
	}
	if authKey == nil {
		return false, errors.New("auth key not found")
	}
	return true, nil
}
