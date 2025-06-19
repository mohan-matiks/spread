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

// MockVersionRepository is a mock implementation of VersionRepository
type MockVersionRepository struct {
	mock.Mock
}

func (m *MockVersionRepository) Create(ctx context.Context, version *model.Version) (*model.Version, error) {
	args := m.Called(ctx, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) UpdateCurrentBundleId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, id, currentBundleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetByEnvironmentIdAndAppVersion(ctx context.Context, environmentId primitive.ObjectID, appVersion string) (*model.Version, error) {
	args := m.Called(ctx, environmentId, appVersion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetByIdAndEnvironmentId(ctx context.Context, versionId primitive.ObjectID, environmentId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, versionId, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetAllByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) ([]*model.Version, error) {
	args := m.Called(ctx, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetById(ctx context.Context, versionId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionRepository) GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Version, error) {
	args := m.Called(ctx, environment, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func TestNewVersionService(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &versionService{}, service)
}

func TestVersionService_CreateVersion_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	version := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   primitive.NewObjectID(),
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Mock Create to return the version
	mockRepo.On("Create", ctx, version).Return(version, nil)

	// Execute
	result, err := service.CreateVersion(ctx, version)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, version.Id, result.Id)
	assert.Equal(t, version.AppVersion, result.AppVersion)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_CreateVersion_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	version := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   primitive.NewObjectID(),
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	// Mock Create to return an error
	mockRepo.On("Create", ctx, version).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.CreateVersion(ctx, version)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_UpdateVersionCurrentBundleIdByVersionId_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	currentBundleId := primitive.NewObjectID()

	expectedVersion := &model.Version{
		Id:              versionId,
		EnvironmentId:   primitive.NewObjectID(),
		AppVersion:      "1.0.0",
		CurrentBundleId: currentBundleId,
		UpdatedAt:       time.Now(),
	}

	// Mock UpdateCurrentBundleId to return the updated version
	mockRepo.On("UpdateCurrentBundleId", ctx, versionId, currentBundleId).Return(expectedVersion, nil)

	// Execute
	result, err := service.UpdateVersionCurrentBundleIdByVersionId(ctx, versionId, currentBundleId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedVersion.Id, result.Id)
	assert.Equal(t, expectedVersion.CurrentBundleId, result.CurrentBundleId)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_UpdateVersionCurrentBundleIdByVersionId_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	currentBundleId := primitive.NewObjectID()

	// Mock UpdateCurrentBundleId to return an error
	mockRepo.On("UpdateCurrentBundleId", ctx, versionId, currentBundleId).Return(nil, errors.New("update error"))

	// Execute
	result, err := service.UpdateVersionCurrentBundleIdByVersionId(ctx, versionId, currentBundleId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "update error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndAppVersion_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()
	appVersion := "1.0.0"

	expectedVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environmentId,
		AppVersion:      appVersion,
		CurrentBundleId: primitive.NewObjectID(),
	}

	// Mock GetByEnvironmentIdAndAppVersion to return a version
	mockRepo.On("GetByEnvironmentIdAndAppVersion", ctx, environmentId, appVersion).Return(expectedVersion, nil)

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndAppVersion(ctx, environmentId, appVersion)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedVersion.AppVersion, result.AppVersion)
	assert.Equal(t, expectedVersion.EnvironmentId, result.EnvironmentId)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndAppVersion_NotFound(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()
	appVersion := "1.0.0"

	// Mock GetByEnvironmentIdAndAppVersion to return mongo.ErrNoDocuments
	mockRepo.On("GetByEnvironmentIdAndAppVersion", ctx, environmentId, appVersion).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndAppVersion(ctx, environmentId, appVersion)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndAppVersion_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()
	appVersion := "1.0.0"

	// Mock GetByEnvironmentIdAndAppVersion to return an error
	mockRepo.On("GetByEnvironmentIdAndAppVersion", ctx, environmentId, appVersion).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndAppVersion(ctx, environmentId, appVersion)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndVersionId_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	environmentId := primitive.NewObjectID()

	expectedVersion := &model.Version{
		Id:              versionId,
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	// Mock GetByIdAndEnvironmentId to return a version
	mockRepo.On("GetByIdAndEnvironmentId", ctx, versionId, environmentId).Return(expectedVersion, nil)

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndVersionId(ctx, versionId, environmentId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedVersion.Id, result.Id)
	assert.Equal(t, expectedVersion.EnvironmentId, result.EnvironmentId)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndVersionId_NotFound(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	environmentId := primitive.NewObjectID()

	// Mock GetByIdAndEnvironmentId to return mongo.ErrNoDocuments
	mockRepo.On("GetByIdAndEnvironmentId", ctx, versionId, environmentId).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndVersionId(ctx, versionId, environmentId)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndVersionId_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	environmentId := primitive.NewObjectID()

	// Mock GetByIdAndEnvironmentId to return an error
	mockRepo.On("GetByIdAndEnvironmentId", ctx, versionId, environmentId).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndVersionId(ctx, versionId, environmentId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetVersionByEnvironmentIdAndVersionId_NilVersion(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()
	environmentId := primitive.NewObjectID()

	// Mock GetByIdAndEnvironmentId to return nil version
	mockRepo.On("GetByIdAndEnvironmentId", ctx, versionId, environmentId).Return(nil, nil)

	// Execute
	result, err := service.GetVersionByEnvironmentIdAndVersionId(ctx, versionId, environmentId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "version not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetLatestVersionByEnvironmentId_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()

	expectedVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
		CreatedAt:       time.Now(),
	}

	// Mock GetLatestVersionByEnvironmentId to return a version
	mockRepo.On("GetLatestVersionByEnvironmentId", ctx, environmentId).Return(expectedVersion, nil)

	// Execute
	result, err := service.GetLatestVersionByEnvironmentId(ctx, environmentId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedVersion.Id, result.Id)
	assert.Equal(t, expectedVersion.EnvironmentId, result.EnvironmentId)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetLatestVersionByEnvironmentId_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()

	// Mock GetLatestVersionByEnvironmentId to return an error
	mockRepo.On("GetLatestVersionByEnvironmentId", ctx, environmentId).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetLatestVersionByEnvironmentId(ctx, environmentId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetAllVersionsByEnvironmentId_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()

	now := time.Now()
	expectedVersions := []*model.Version{
		{
			Id:              primitive.NewObjectID(),
			EnvironmentId:   environmentId,
			AppVersion:      "1.1.0",
			CurrentBundleId: primitive.NewObjectID(),
			CreatedAt:       now, // newer
		},
		{
			Id:              primitive.NewObjectID(),
			EnvironmentId:   environmentId,
			AppVersion:      "1.0.0",
			CurrentBundleId: primitive.NewObjectID(),
			CreatedAt:       now.Add(-time.Hour), // older
		},
	}

	// Mock GetAllByEnvironmentId to return versions
	mockRepo.On("GetAllByEnvironmentId", ctx, environmentId).Return(expectedVersions, nil)

	// Execute
	result, err := service.GetAllVersionsByEnvironmentId(ctx, environmentId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	// Should be sorted by createdAt in descending order (newest first)
	assert.Equal(t, expectedVersions[0].AppVersion, result[0].AppVersion) // 1.1.0 first (newer)
	assert.Equal(t, expectedVersions[1].AppVersion, result[1].AppVersion) // 1.0.0 second (older)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetAllVersionsByEnvironmentId_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	environmentId := primitive.NewObjectID()

	// Mock GetAllByEnvironmentId to return an error
	mockRepo.On("GetAllByEnvironmentId", ctx, environmentId).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetAllVersionsByEnvironmentId(ctx, environmentId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetByVersionId_Success(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()

	expectedVersion := &model.Version{
		Id:              versionId,
		EnvironmentId:   primitive.NewObjectID(),
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	// Mock GetById to return a version
	mockRepo.On("GetById", ctx, versionId).Return(expectedVersion, nil)

	// Execute
	result, err := service.GetByVersionId(ctx, versionId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedVersion.Id, result.Id)
	assert.Equal(t, expectedVersion.AppVersion, result.AppVersion)

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetByVersionId_NotFound(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()

	// Mock GetById to return nil version
	mockRepo.On("GetById", ctx, versionId).Return(nil, nil)

	// Execute
	result, err := service.GetByVersionId(ctx, versionId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "version not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestVersionService_GetByVersionId_Error(t *testing.T) {
	mockRepo := &MockVersionRepository{}
	service := NewVersionService(mockRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()

	// Mock GetById to return an error
	mockRepo.On("GetById", ctx, versionId).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetByVersionId(ctx, versionId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}
