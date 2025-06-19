package service

import (
	"context"
	"errors"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/types"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user *types.CreateUserRequest) (*model.User, error)
	Login(user *types.LoginUserRequest) (*string, error)
	GetUser(id string) (*model.User, error)
	Count(ctx context.Context) (int64, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Create(user *types.CreateUserRequest) (*model.User, error) {
	// check if user already exists
	existingUser, err := s.userRepository.GetByUsername(context.Background(), user.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}
	// encrypto password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userModel := &model.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Roles:    user.Roles,
	}
	createdUser, err := s.userRepository.Insert(context.Background(), userModel)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *userService) Login(user *types.LoginUserRequest) (*string, error) {
	existingUser, err := s.userRepository.GetByUsername(context.Background(), user.Username)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}
	// create a access token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": existingUser.Id,
	}).SignedString([]byte(config.TokenSecret))
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func (s *userService) GetUser(id string) (*model.User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetById(context.Background(), userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Count(ctx context.Context) (int64, error) {
	return s.userRepository.Count(ctx)
}
