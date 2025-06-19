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
	"go.mongodb.org/mongo-driver/mongo"
)

// MockAppRepository is a mock implementation of AppRepository
type MockAppRepository struct {
	mock.Mock
}

func (m *MockAppRepository) Insert(ctx context.Context, app *model.App) (*model.App, error) {
	args := m.Called(ctx, app)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppRepository) GetByName(ctx context.Context, name string) (*model.App, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.App, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.App), args.Error(1)
}

func (m *MockAppRepository) GetAll(ctx context.Context) ([]*model.App, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.App), args.Error(1)
}

func TestNewAppService(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &appServiceImpl{}, service)
}

func TestAppService_CreateApp_Success(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"
	os := "ios"

	expectedApp := &model.App{
		Id:        primitive.NewObjectID(),
		Name:      appName,
		OS:        os,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock GetByName to return nil (app doesn't exist)
	mockRepo.On("GetByName", ctx, appName).Return(nil, mongo.ErrNoDocuments)

	// Mock Insert to return the created app
	mockRepo.On("Insert", ctx, mock.AnythingOfType("*model.App")).Return(expectedApp, nil)

	// Execute
	result, err := service.CreateApp(ctx, appName, os)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedApp.Name, result.Name)
	assert.Equal(t, expectedApp.OS, result.OS)

	mockRepo.AssertExpectations(t)
}

func TestAppService_CreateApp_AlreadyExists(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "existing-app"
	os := "ios"

	existingApp := &model.App{
		Id:   primitive.NewObjectID(),
		Name: appName,
		OS:   os,
	}

	// Mock GetByName to return an existing app
	mockRepo.On("GetByName", ctx, appName).Return(existingApp, nil)

	// Execute
	result, err := service.CreateApp(ctx, appName, os)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "app with name "+appName+" already exists", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_CreateApp_InvalidOS(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"
	os := "invalid-os"

	// Execute
	result, err := service.CreateApp(ctx, appName, os)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid os", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_CreateApp_GetByNameError(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"
	os := "ios"

	// Mock GetByName to return an error (not mongo.ErrNoDocuments)
	mockRepo.On("GetByName", ctx, appName).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.CreateApp(ctx, appName, os)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_CreateApp_InsertError(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"
	os := "ios"

	// Mock GetByName to return nil (app doesn't exist)
	mockRepo.On("GetByName", ctx, appName).Return(nil, mongo.ErrNoDocuments)

	// Mock Insert to return an error
	mockRepo.On("Insert", ctx, mock.AnythingOfType("*model.App")).Return(nil, errors.New("insert error"))

	// Execute
	result, err := service.CreateApp(ctx, appName, os)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insert error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppByName_Success(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"

	expectedApp := &model.App{
		Id:   primitive.NewObjectID(),
		Name: appName,
		OS:   "ios",
	}

	// Mock GetByName to return an app
	mockRepo.On("GetByName", ctx, appName).Return(expectedApp, nil)

	// Execute
	result, err := service.GetAppByName(ctx, appName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedApp.Name, result.Name)
	assert.Equal(t, expectedApp.OS, result.OS)

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppByName_Error(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appName := "test-app"

	// Mock GetByName to return an error
	mockRepo.On("GetByName", ctx, appName).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetAppByName(ctx, appName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppById_Success(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	appIDString := appID.Hex()

	expectedApp := &model.App{
		Id:   appID,
		Name: "test-app",
		OS:   "ios",
	}

	// Mock GetById to return an app
	mockRepo.On("GetById", ctx, appID).Return(expectedApp, nil)

	// Execute
	result, err := service.GetAppById(ctx, appIDString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedApp.Id, result.Id)
	assert.Equal(t, expectedApp.Name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppById_InvalidID(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	invalidID := "invalid-id"

	// Execute
	result, err := service.GetAppById(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "error converting id", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppById_NotFound(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	appIDString := appID.Hex()

	// Mock GetById to return nil (app not found)
	mockRepo.On("GetById", ctx, appID).Return(nil, nil)

	// Execute
	result, err := service.GetAppById(ctx, appIDString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "app not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetAppById_Error(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()
	appID := primitive.NewObjectID()
	appIDString := appID.Hex()

	// Mock GetById to return an error
	mockRepo.On("GetById", ctx, appID).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetAppById(ctx, appIDString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetApps_Success(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()

	expectedApps := []*model.App{
		{
			Id:   primitive.NewObjectID(),
			Name: "app-1",
			OS:   "ios",
		},
		{
			Id:   primitive.NewObjectID(),
			Name: "app-2",
			OS:   "android",
		},
	}

	// Mock GetAll to return apps
	mockRepo.On("GetAll", ctx).Return(expectedApps, nil)

	// Execute
	result, err := service.GetApps(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedApps[0].Name, result[0].Name)
	assert.Equal(t, expectedApps[1].Name, result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestAppService_GetApps_Error(t *testing.T) {
	mockRepo := &MockAppRepository{}
	service := NewAppService(mockRepo)

	ctx := context.Background()

	// Mock GetAll to return an error
	mockRepo.On("GetAll", ctx).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetApps(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}
