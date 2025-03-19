package service

import (
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/utils"
)

type AuthKeyService interface {
	CreateAuthKey(name string, username string) (string, error)
	GetByAuthKey(key string) (*model.AuthKey, error)
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

func (s *authKeyService) GetByAuthKey(key string) (*model.AuthKey, error) {
	authKey, err := s.authKeyRepository.GetById(key)
	if err != nil {
		return nil, err
	}
	return authKey, nil
}
