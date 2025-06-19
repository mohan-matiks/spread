package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockAuthKeyRepository is a mock implementation of AuthKeyRepository
type MockAuthKeyRepository struct {
	mock.Mock
}

func (m *MockAuthKeyRepository) Insert(authKey *model.AuthKey) (*model.AuthKey, error) {
	args := m.Called(authKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthKey), args.Error(1)
}

func (m *MockAuthKeyRepository) GetById(key string) (*model.AuthKey, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthKey), args.Error(1)
}

func (m *MockAuthKeyRepository) GetAll(ctx context.Context) ([]*model.AuthKey, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.AuthKey), args.Error(1)
}

func TestNewAuthKeyService(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &authKeyService{}, service)
}

func TestAuthKeyService_CreateAuthKey_Success(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	name := "test-auth-key"
	username := "testuser"

	// Mock Insert to return a created auth key
	expectedAuthKey := &model.AuthKey{
		Id:        primitive.NewObjectID(),
		Name:      name,
		Key:       "generated-key-123",
		IsValid:   true,
		CreatedBy: username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("Insert", mock.AnythingOfType("*model.AuthKey")).Return(expectedAuthKey, nil)

	// Execute
	result, err := service.CreateAuthKey(name, username)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, expectedAuthKey.Key, result)

	// Verify the auth key was created with correct fields
	mockRepo.AssertExpectations(t)
	insertCall := mockRepo.Calls[0]
	authKeyArg := insertCall.Arguments[0].(*model.AuthKey)
	assert.Equal(t, name, authKeyArg.Name)
	assert.Equal(t, username, authKeyArg.CreatedBy)
	assert.True(t, authKeyArg.IsValid)
	assert.NotEmpty(t, authKeyArg.Key) // Should be generated
}

func TestAuthKeyService_CreateAuthKey_InsertError(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	name := "test-auth-key"
	username := "testuser"

	// Mock Insert to return an error
	mockRepo.On("Insert", mock.AnythingOfType("*model.AuthKey")).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.CreateAuthKey(name, username)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetByAuthKey_Success(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	key := "test-key-123"

	// Mock GetById to return an auth key
	expectedAuthKey := &model.AuthKey{
		Id:        primitive.NewObjectID(),
		Name:      "test-auth-key",
		Key:       key,
		IsValid:   true,
		CreatedBy: "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("GetById", key).Return(expectedAuthKey, nil)

	// Execute
	result, err := service.GetByAuthKey(key)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAuthKey.Key, result.Key)
	assert.Equal(t, expectedAuthKey.Name, result.Name)
	assert.Equal(t, expectedAuthKey.CreatedBy, result.CreatedBy)

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetByAuthKey_NotFound(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	key := "nonexistent-key"

	// Mock GetById to return nil (not found)
	mockRepo.On("GetById", key).Return(nil, nil)

	// Execute
	result, err := service.GetByAuthKey(key)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetByAuthKey_RepositoryError(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	key := "test-key-123"

	// Mock GetById to return an error
	mockRepo.On("GetById", key).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetByAuthKey(key)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetAllAuthKeys_Success(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	ctx := context.Background()

	// Mock GetAll to return auth keys
	expectedAuthKeys := []*model.AuthKey{
		{
			Id:        primitive.NewObjectID(),
			Name:      "auth-key-1",
			Key:       "key-1",
			IsValid:   true,
			CreatedBy: "user1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        primitive.NewObjectID(),
			Name:      "auth-key-2",
			Key:       "key-2",
			IsValid:   false,
			CreatedBy: "user2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockRepo.On("GetAll", ctx).Return(expectedAuthKeys, nil)

	// Execute
	result, err := service.GetAllAuthKeys(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedAuthKeys[0].Name, result[0].Name)
	assert.Equal(t, expectedAuthKeys[1].Name, result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetAllAuthKeys_EmptyList(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	ctx := context.Background()

	// Mock GetAll to return empty list
	mockRepo.On("GetAll", ctx).Return([]*model.AuthKey{}, nil)

	// Execute
	result, err := service.GetAllAuthKeys(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_GetAllAuthKeys_RepositoryError(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	ctx := context.Background()

	// Mock GetAll to return an error
	mockRepo.On("GetAll", ctx).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetAllAuthKeys(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAuthKeyService_CreateAuthKey_GeneratesUniqueKeys(t *testing.T) {
	mockRepo := &MockAuthKeyRepository{}
	service := NewAuthKeyService(mockRepo)

	name := "test-auth-key"
	username := "testuser"

	// Mock Insert to return different keys for each call
	authKey1 := &model.AuthKey{
		Id:        primitive.NewObjectID(),
		Name:      name,
		Key:       "unique-key-1",
		IsValid:   true,
		CreatedBy: username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	authKey2 := &model.AuthKey{
		Id:        primitive.NewObjectID(),
		Name:      name,
		Key:       "unique-key-2",
		IsValid:   true,
		CreatedBy: username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Insert", mock.AnythingOfType("*model.AuthKey")).Return(authKey1, nil).Once()
	mockRepo.On("Insert", mock.AnythingOfType("*model.AuthKey")).Return(authKey2, nil).Once()

	// Execute twice
	result1, err1 := service.CreateAuthKey(name, username)
	result2, err2 := service.CreateAuthKey(name, username)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, result1)
	assert.NotEmpty(t, result2)
	assert.NotEqual(t, result1, result2) // Keys should be different

	mockRepo.AssertExpectations(t)
}
