package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Insert(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &userService{}, service)
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	createRequest := &types.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Roles:    []string{"user"},
	}

	// Mock GetByUsername to return nil (user doesn't exist)
	mockRepo.On("GetByUsername", ctx, "testuser").Return(nil, nil)

	// Mock Insert to return a created user
	expectedUser := &model.User{
		Id:        primitive.NewObjectID(),
		Username:  "testuser",
		Password:  "hashedpassword",
		Roles:     []string{"user"},
		IsValid:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("Insert", ctx, mock.AnythingOfType("*model.User")).Return(expectedUser, nil)

	// Execute
	result, err := service.Create(createRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Roles, result.Roles)
	assert.NotEqual(t, createRequest.Password, result.Password) // Password should be hashed

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_UserAlreadyExists(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	createRequest := &types.CreateUserRequest{
		Username: "existinguser",
		Password: "password123",
		Roles:    []string{"user"},
	}

	// Mock GetByUsername to return an existing user
	existingUser := &model.User{
		Id:       primitive.NewObjectID(),
		Username: "existinguser",
	}
	mockRepo.On("GetByUsername", ctx, "existinguser").Return(existingUser, nil)

	// Execute
	result, err := service.Create(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user already exists", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_GetByUsernameError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	createRequest := &types.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Roles:    []string{"user"},
	}

	// Mock GetByUsername to return an error
	mockRepo.On("GetByUsername", ctx, "testuser").Return(nil, errors.New("database error"))

	// Execute
	result, err := service.Create(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_InsertError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	createRequest := &types.CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Roles:    []string{"user"},
	}

	// Mock GetByUsername to return nil (user doesn't exist)
	mockRepo.On("GetByUsername", ctx, "testuser").Return(nil, nil)

	// Mock Insert to return an error
	mockRepo.On("Insert", ctx, mock.AnythingOfType("*model.User")).Return(nil, errors.New("insert error"))

	// Execute
	result, err := service.Create(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insert error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	loginRequest := &types.LoginUserRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Create a hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Mock GetByUsername to return an existing user
	existingUser := &model.User{
		Id:       primitive.NewObjectID(),
		Username: "testuser",
		Password: string(hashedPassword),
		Roles:    []string{"user"},
	}
	mockRepo.On("GetByUsername", ctx, "testuser").Return(existingUser, nil)

	// Execute
	result, err := service.Login(loginRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify the token is valid
	token, parseErr := jwt.Parse(*result, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenSecret), nil
	})
	assert.NoError(t, parseErr)
	assert.True(t, token.Valid)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	loginRequest := &types.LoginUserRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}

	// Mock GetByUsername to return nil (user doesn't exist)
	mockRepo.On("GetByUsername", ctx, "nonexistentuser").Return(nil, nil)

	// Execute
	result, err := service.Login(loginRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	loginRequest := &types.LoginUserRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Create a hashed password for correct password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// Mock GetByUsername to return an existing user
	existingUser := &model.User{
		Id:       primitive.NewObjectID(),
		Username: "testuser",
		Password: string(hashedPassword),
		Roles:    []string{"user"},
	}
	mockRepo.On("GetByUsername", ctx, "testuser").Return(existingUser, nil)

	// Execute
	result, err := service.Login(loginRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid password", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_GetByUsernameError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	loginRequest := &types.LoginUserRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Mock GetByUsername to return an error
	mockRepo.On("GetByUsername", ctx, "testuser").Return(nil, errors.New("database error"))

	// Execute
	result, err := service.Login(loginRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	userIDString := userID.Hex()

	// Mock GetById to return a user
	expectedUser := &model.User{
		Id:       userID,
		Username: "testuser",
		Roles:    []string{"user"},
	}
	mockRepo.On("GetById", ctx, userID).Return(expectedUser, nil)

	// Execute
	result, err := service.GetUser(userIDString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Id, result.Id)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_InvalidID(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	// Execute with invalid ID
	result, err := service.GetUser("invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	userIDString := userID.Hex()

	// Mock GetById to return mongo.ErrNoDocuments
	mockRepo.On("GetById", ctx, userID).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetUser(userIDString)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	userIDString := userID.Hex()

	// Mock GetById to return an error
	mockRepo.On("GetById", ctx, userID).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetUser(userIDString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Count_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	expectedCount := int64(5)

	// Mock Count to return a count
	mockRepo.On("Count", ctx).Return(expectedCount, nil)

	// Execute
	result, err := service.Count(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Count_Error(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	ctx := context.Background()

	// Mock Count to return an error
	mockRepo.On("Count", ctx).Return(int64(0), errors.New("database error"))

	// Execute
	result, err := service.Count(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}
