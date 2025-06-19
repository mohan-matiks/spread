package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockAppService is a mock implementation of AppService
type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) CreateApp(ctx context.Context, appName string, os string) (*model.App, error) {
	args := m.Called(ctx, appName, os)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppService) GetAppByName(ctx context.Context, appName string) (*model.App, error) {
	args := m.Called(ctx, appName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppService) GetApps(ctx context.Context) ([]*model.App, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.App), args.Error(1)
}

func (m *MockAppService) GetAppById(ctx context.Context, id string) (*model.App, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

// MockEnvironmentRepository is a mock implementation of EnvironmentRepository
type MockEnvironmentRepository struct {
	mock.Mock
}

func (m *MockEnvironmentRepository) Insert(ctx context.Context, environment *model.Environment) (*model.Environment, error) {
	args := m.Called(ctx, environment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentRepository) GetByAppIdAndName(ctx context.Context, appId primitive.ObjectID, name string) (*model.Environment, error) {
	args := m.Called(ctx, appId, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentRepository) GetByKey(ctx context.Context, key string) (*model.Environment, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentRepository) GetAllByAppId(ctx context.Context, appId primitive.ObjectID) ([]*model.Environment, error) {
	args := m.Called(ctx, appId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Environment), args.Error(1)
}

func (m *MockEnvironmentRepository) GetByIdAndAppId(ctx context.Context, id primitive.ObjectID, appId primitive.ObjectID) (*model.Environment, error) {
	args := m.Called(ctx, id, appId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func TestNewEnvironmentService(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &environmentServiceImpl{}, service)
}

func TestEnvironmentService_CreateEnvironment_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "dev"
	appID := primitive.NewObjectID()

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return an app
	expectedApp := &model.App{
		Id:   appID,
		Name: appName,
		OS:   "ios",
	}
	mockAppService.On("GetAppByName", ctx, appName).Return(expectedApp, nil)

	// Mock GetByAppIdAndName to return nil (environment doesn't exist)
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(nil, mongo.ErrNoDocuments)

	// Mock Insert to return a created environment
	expectedEnvironment := &model.Environment{
		Id:        primitive.NewObjectID(),
		AppId:     appID,
		Name:      environmentName,
		Key:       "generated-key",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockEnvRepo.On("Insert", ctx, mock.AnythingOfType("*model.Environment")).Return(expectedEnvironment, nil)

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEnvironment.Name, result.Name)
	assert.Equal(t, expectedEnvironment.AppId, result.AppId)
	assert.NotEmpty(t, result.Key)

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_CreateEnvironment_AppNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "nonexistent-app"
	environmentName := "dev"

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return mongo.ErrNoDocuments
	mockAppService.On("GetAppByName", ctx, appName).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "app not found", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_CreateEnvironment_AppServiceError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "dev"

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return an error
	mockAppService.On("GetAppByName", ctx, appName).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_CreateEnvironment_EnvironmentAlreadyExists(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "dev"
	appID := primitive.NewObjectID()

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return an app
	expectedApp := &model.App{
		Id:   appID,
		Name: appName,
		OS:   "ios",
	}
	mockAppService.On("GetAppByName", ctx, appName).Return(expectedApp, nil)

	// Mock GetByAppIdAndName to return an existing environment
	existingEnvironment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: appID,
		Name:  environmentName,
	}
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(existingEnvironment, nil)

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "environment with name "+environmentName+" already exists for app "+appName, err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_CreateEnvironment_GetByAppIdAndNameError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "dev"
	appID := primitive.NewObjectID()

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return an app
	expectedApp := &model.App{
		Id:   appID,
		Name: appName,
		OS:   "ios",
	}
	mockAppService.On("GetAppByName", ctx, appName).Return(expectedApp, nil)

	// Mock GetByAppIdAndName to return an error (not mongo.ErrNoDocuments)
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_CreateEnvironment_InsertError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appName := "test-app"
	environmentName := "dev"
	appID := primitive.NewObjectID()

	createRequest := &types.CreateEnvironmentRequest{
		AppName:         appName,
		EnvironmentName: environmentName,
	}

	// Mock GetAppByName to return an app
	expectedApp := &model.App{
		Id:   appID,
		Name: appName,
		OS:   "ios",
	}
	mockAppService.On("GetAppByName", ctx, appName).Return(expectedApp, nil)

	// Mock GetByAppIdAndName to return nil (environment doesn't exist)
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(nil, mongo.ErrNoDocuments)

	// Mock Insert to return an error
	mockEnvRepo.On("Insert", ctx, mock.AnythingOfType("*model.Environment")).Return(nil, errors.New("insert error"))

	// Execute
	result, err := service.CreateEnvironment(createRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insert error", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndName_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	environmentName := "dev"

	// Mock GetByAppIdAndName to return an environment
	expectedEnvironment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: appID,
		Name:  environmentName,
		Key:   "env-key",
	}
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(expectedEnvironment, nil)

	// Execute
	result, err := service.GetEnvironmentByAppIdAndName(ctx, appID, environmentName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEnvironment.Name, result.Name)
	assert.Equal(t, expectedEnvironment.AppId, result.AppId)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndName_NotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	environmentName := "nonexistent"

	// Mock GetByAppIdAndName to return mongo.ErrNoDocuments
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetEnvironmentByAppIdAndName(ctx, appID, environmentName)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndName_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	environmentName := "dev"

	// Mock GetByAppIdAndName to return an error
	mockEnvRepo.On("GetByAppIdAndName", ctx, appID, environmentName).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetEnvironmentByAppIdAndName(ctx, appID, environmentName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByKey_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	key := "test-key"

	// Mock GetByKey to return an environment
	expectedEnvironment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: primitive.NewObjectID(),
		Name:  "dev",
		Key:   key,
	}
	mockEnvRepo.On("GetByKey", ctx, key).Return(expectedEnvironment, nil)

	// Execute
	result, err := service.GetEnvironmentByKey(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEnvironment.Key, result.Key)
	assert.Equal(t, expectedEnvironment.Name, result.Name)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByKey_NotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	key := "nonexistent-key"

	// Mock GetByKey to return mongo.ErrNoDocuments
	mockEnvRepo.On("GetByKey", ctx, key).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetEnvironmentByKey(ctx, key)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByKey_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	key := "test-key"

	// Mock GetByKey to return an error
	mockEnvRepo.On("GetByKey", ctx, key).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetEnvironmentByKey(ctx, key)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetAllEnvironmentsByAppId_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()

	// Mock GetAllByAppId to return environments
	expectedEnvironments := []*model.Environment{
		{
			Id:    primitive.NewObjectID(),
			AppId: appID,
			Name:  "dev",
			Key:   "dev-key",
		},
		{
			Id:    primitive.NewObjectID(),
			AppId: appID,
			Name:  "staging",
			Key:   "staging-key",
		},
	}
	mockEnvRepo.On("GetAllByAppId", ctx, appID).Return(expectedEnvironments, nil)

	// Execute
	result, err := service.GetAllEnvironmentsByAppId(ctx, appID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedEnvironments[0].Name, result[0].Name)
	assert.Equal(t, expectedEnvironments[1].Name, result[1].Name)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetAllEnvironmentsByAppId_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()

	// Mock GetAllByAppId to return an error
	mockEnvRepo.On("GetAllByAppId", ctx, appID).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetAllEnvironmentsByAppId(ctx, appID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndEnvironmentId_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	environmentID := primitive.NewObjectID()
	environmentIDString := environmentID.Hex()

	// Mock GetByIdAndAppId to return an environment
	expectedEnvironment := &model.Environment{
		Id:    environmentID,
		AppId: appID,
		Name:  "dev",
		Key:   "dev-key",
	}
	mockEnvRepo.On("GetByIdAndAppId", ctx, environmentID, appID).Return(expectedEnvironment, nil)

	// Execute
	result, err := service.GetEnvironmentByAppIdAndEnvironmentId(ctx, appID, environmentIDString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEnvironment.Id, result.Id)
	assert.Equal(t, expectedEnvironment.AppId, result.AppId)

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndEnvironmentId_InvalidID(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	invalidID := "invalid-id"

	// Execute
	result, err := service.GetEnvironmentByAppIdAndEnvironmentId(ctx, appID, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "error converting environment id", err.Error())

	mockEnvRepo.AssertExpectations(t)
}

func TestEnvironmentService_GetEnvironmentByAppIdAndEnvironmentId_NotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvRepo := &MockEnvironmentRepository{}
	service := NewEnvironmentService(mockAppService, mockEnvRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	environmentID := primitive.NewObjectID()
	environmentIDString := environmentID.Hex()

	// Mock GetByIdAndAppId to return an error
	mockEnvRepo.On("GetByIdAndAppId", ctx, environmentID, appID).Return(nil, errors.New("not found"))

	// Execute
	result, err := service.GetEnvironmentByAppIdAndEnvironmentId(ctx, appID, environmentIDString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "environment not found", err.Error())

	mockEnvRepo.AssertExpectations(t)
}
