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

// MockVersionService is a mock implementation of VersionService
type MockVersionService struct {
	mock.Mock
}

func (m *MockVersionService) CreateVersion(ctx context.Context, version *model.Version) (*model.Version, error) {
	args := m.Called(ctx, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) UpdateVersionCurrentBundleIdByVersionId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, id, currentBundleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) GetVersionByEnvironmentIdAndAppVersion(ctx context.Context, environmentId primitive.ObjectID, appVersion string) (*model.Version, error) {
	args := m.Called(ctx, environmentId, appVersion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) GetVersionByEnvironmentIdAndVersionId(ctx context.Context, versionId primitive.ObjectID, environmentId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, versionId, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) GetAllVersionsByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) ([]*model.Version, error) {
	args := m.Called(ctx, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Version), args.Error(1)
}

func (m *MockVersionService) GetByVersionId(ctx context.Context, versionId primitive.ObjectID) (*model.Version, error) {
	args := m.Called(ctx, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

func (m *MockVersionService) GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Version, error) {
	args := m.Called(ctx, environment, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Version), args.Error(1)
}

// MockEnvironmentService is a mock implementation of EnvironmentService
type MockEnvironmentService struct {
	mock.Mock
}

func (m *MockEnvironmentService) CreateEnvironment(createRequest *types.CreateEnvironmentRequest) (*model.Environment, error) {
	args := m.Called(createRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentService) GetEnvironmentByAppIdAndName(ctx context.Context, appId primitive.ObjectID, name string) (*model.Environment, error) {
	args := m.Called(ctx, appId, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentService) GetEnvironmentByKey(ctx context.Context, key string) (*model.Environment, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

func (m *MockEnvironmentService) GetAllEnvironmentsByAppId(ctx context.Context, appId primitive.ObjectID) ([]*model.Environment, error) {
	args := m.Called(ctx, appId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Environment), args.Error(1)
}

func (m *MockEnvironmentService) GetEnvironmentByAppIdAndEnvironmentId(ctx context.Context, appId primitive.ObjectID, environmentId string) (*model.Environment, error) {
	args := m.Called(ctx, appId, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Environment), args.Error(1)
}

// MockBundleRepository is a mock implementation of BundleRepository
type MockBundleRepository struct {
	mock.Mock
}

func (m *MockBundleRepository) CreateBundle(ctx context.Context, bundle *model.Bundle) (*model.Bundle, error) {
	args := m.Called(ctx, bundle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetByHashAndVersionId(ctx context.Context, hash string, versionId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(ctx, hash, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) UpdateVersionIdById(ctx context.Context, id primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(ctx, id, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Bundle, error) {
	args := m.Called(ctx, environment, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetNextSeqByEnvironmentIdAndVersionId(ctx context.Context, environmentId primitive.ObjectID, versionId primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, environmentId, versionId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBundleRepository) GetBySequenceIdEnvironmentIdAndVersionId(ctx context.Context, sequenceId int64, environmentId primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(ctx, sequenceId, environmentId, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetByLabelAndEnvironmentId(ctx context.Context, label string, environmentId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(ctx, label, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetAllByVersionId(ctx context.Context, versionId primitive.ObjectID) ([]*model.Bundle, error) {
	args := m.Called(ctx, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) UpdateIsMandatoryById(ctx context.Context, id primitive.ObjectID, isMandatory bool) (*model.Bundle, error) {
	args := m.Called(ctx, id, isMandatory)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) UpdateIsValid(ctx context.Context, id primitive.ObjectID, isValid bool) (*model.Bundle, error) {
	args := m.Called(ctx, id, isValid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleRepository) AddActive(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleRepository) AddFailed(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleRepository) AddInstalled(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleRepository) DecrementActive(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewBundleService(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}

	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &bundleService{}, service)
}

func TestBundleService_GetBundleById_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleID := primitive.NewObjectID()

	expectedBundle := &model.Bundle{
		Id:            bundleID,
		AppId:         primitive.NewObjectID(),
		EnvironmentId: primitive.NewObjectID(),
		DownloadFile:  "test-bundle.js",
		Size:          1024,
		Hash:          "test-hash",
		Description:   "Test bundle",
		IsMandatory:   false,
		IsValid:       true,
	}

	mockBundleRepo.On("GetById", ctx, bundleID).Return(expectedBundle, nil)

	result, err := service.GetBundleById(bundleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBundle.Id, result.Id)
	assert.Equal(t, expectedBundle.DownloadFile, result.DownloadFile)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleById_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleID := primitive.NewObjectID()

	mockBundleRepo.On("GetById", ctx, bundleID).Return(nil, errors.New("database error"))

	result, err := service.GetBundleById(bundleID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByLabelAndEnvironmentId_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	label := "v1x1"
	environmentId := primitive.NewObjectID()

	expectedBundle := &model.Bundle{
		Id:            primitive.NewObjectID(),
		Label:         label,
		EnvironmentId: environmentId,
		DownloadFile:  "test-bundle.js",
		Size:          1024,
		Hash:          "test-hash",
	}

	mockBundleRepo.On("GetByLabelAndEnvironmentId", ctx, label, environmentId).Return(expectedBundle, nil)

	result, err := service.GetBundleByLabelAndEnvironmentId(label, environmentId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBundle.Label, result.Label)
	assert.Equal(t, expectedBundle.EnvironmentId, result.EnvironmentId)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByLabelAndEnvironmentId_NotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	label := "v1x1"
	environmentId := primitive.NewObjectID()

	mockBundleRepo.On("GetByLabelAndEnvironmentId", ctx, label, environmentId).Return(nil, mongo.ErrNoDocuments)

	result, err := service.GetBundleByLabelAndEnvironmentId(label, environmentId)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByLabelAndEnvironmentId_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	label := "v1x1"
	environmentId := primitive.NewObjectID()

	mockBundleRepo.On("GetByLabelAndEnvironmentId", ctx, label, environmentId).Return(nil, errors.New("database error"))

	result, err := service.GetBundleByLabelAndEnvironmentId(label, environmentId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundlesByVersionId_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()

	expectedBundles := []*model.Bundle{
		{
			Id:           primitive.NewObjectID(),
			VersionId:    versionId,
			DownloadFile: "bundle1.js",
			Size:         1024,
			Hash:         "hash1",
			CreatedAt:    time.Now().Add(-time.Hour),
		},
		{
			Id:           primitive.NewObjectID(),
			VersionId:    versionId,
			DownloadFile: "bundle2.js",
			Size:         2048,
			Hash:         "hash2",
			CreatedAt:    time.Now(),
		},
	}

	mockBundleRepo.On("GetAllByVersionId", ctx, versionId).Return(expectedBundles, nil)

	result, err := service.GetBundlesByVersionId(versionId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	// Should be sorted by createdAt in descending order (newest first)
	// bundle2 (newer) should be first, bundle1 (older) should be second
	assert.Equal(t, "https://dev-spread-bucket.justswish.in/bundle2.js", result[0].DownloadFile) // bundle2 first (newer)
	assert.Equal(t, "https://dev-spread-bucket.justswish.in/bundle1.js", result[1].DownloadFile) // bundle1 second (older)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundlesByVersionId_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	versionId := primitive.NewObjectID()

	mockBundleRepo.On("GetAllByVersionId", ctx, versionId).Return(nil, errors.New("database error"))

	result, err := service.GetBundlesByVersionId(versionId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_ToggleMandatory_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	existingBundle := &model.Bundle{
		Id:          bundleId,
		IsMandatory: false,
	}

	updatedBundle := &model.Bundle{
		Id:          bundleId,
		IsMandatory: true,
	}

	mockBundleRepo.On("GetById", ctx, bundleId).Return(existingBundle, nil)
	mockBundleRepo.On("UpdateIsMandatoryById", ctx, bundleId, true).Return(updatedBundle, nil)

	err := service.ToggleMandatory(bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_ToggleMandatory_GetByIdError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("GetById", ctx, bundleId).Return(nil, errors.New("database error"))

	err := service.ToggleMandatory(bundleId)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_ToggleActive_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	existingBundle := &model.Bundle{
		Id:      bundleId,
		IsValid: false,
	}

	updatedBundle := &model.Bundle{
		Id:      bundleId,
		IsValid: true,
	}

	mockBundleRepo.On("GetById", ctx, bundleId).Return(existingBundle, nil)
	mockBundleRepo.On("UpdateIsValid", ctx, bundleId, true).Return(updatedBundle, nil)

	err := service.ToggleActive(bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_ToggleActive_GetByIdError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("GetById", ctx, bundleId).Return(nil, errors.New("database error"))

	err := service.ToggleActive(bundleId)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_AddActive_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("AddActive", ctx, bundleId).Return(nil)

	err := service.AddActive(ctx, bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_AddActive_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("AddActive", ctx, bundleId).Return(errors.New("database error"))

	err := service.AddActive(ctx, bundleId)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_AddFailed_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("AddFailed", ctx, bundleId).Return(nil)

	err := service.AddFailed(ctx, bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_AddInstalled_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("AddInstalled", ctx, bundleId).Return(nil)

	err := service.AddInstalled(ctx, bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_DecrementActive_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockBundleRepo)

	ctx := context.Background()
	bundleId := primitive.NewObjectID()

	mockBundleRepo.On("DecrementActive", ctx, bundleId).Return(nil)

	err := service.DecrementActive(ctx, bundleId)

	assert.NoError(t, err)

	mockBundleRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByHashAndVersionId_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	hash := "test-hash"
	versionId := primitive.NewObjectID()

	expectedBundle := &model.Bundle{
		Id:        primitive.NewObjectID(),
		Hash:      hash,
		VersionId: versionId,
	}

	// Mock GetByHashAndVersionId to return a bundle
	mockRepo.On("GetByHashAndVersionId", ctx, hash, versionId).Return(expectedBundle, nil)

	// Execute
	result, err := service.GetBundleByHashAndVersionId(hash, versionId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBundle.Hash, result.Hash)
	assert.Equal(t, expectedBundle.VersionId, result.VersionId)

	mockRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByHashAndVersionId_NotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	hash := "test-hash"
	versionId := primitive.NewObjectID()

	// Mock GetByHashAndVersionId to return mongo.ErrNoDocuments
	mockRepo.On("GetByHashAndVersionId", ctx, hash, versionId).Return(nil, mongo.ErrNoDocuments)

	// Execute
	result, err := service.GetBundleByHashAndVersionId(hash, versionId)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestBundleService_GetBundleByHashAndVersionId_Error(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	hash := "test-hash"
	versionId := primitive.NewObjectID()

	// Mock GetByHashAndVersionId to return an error
	mockRepo.On("GetByHashAndVersionId", ctx, hash, versionId).Return(nil, errors.New("database error"))

	// Execute
	result, err := service.GetBundleByHashAndVersionId(hash, versionId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestBundleService_CreateNewBundle_NewVersion_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	createdBy := "test-user"
	payload := &types.CreateNewBundleRequest{
		AppName:      "test-app",
		Environment:  "dev",
		DownloadFile: "test-bundle.js",
		Description:  "Test bundle",
		AppVersion:   "1.0.0",
		Size:         1024,
		Hash:         "test-hash",
	}

	app := &model.App{
		Id:   primitive.NewObjectID(),
		Name: payload.AppName,
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: app.Id,
		Name:  payload.Environment,
	}

	expectedBundle := &model.Bundle{
		Id:            primitive.NewObjectID(),
		AppId:         app.Id,
		EnvironmentId: environment.Id,
		DownloadFile:  payload.DownloadFile,
		Size:          payload.Size,
		Hash:          payload.Hash,
		Description:   payload.Description,
		IsMandatory:   false,
		Failed:        0,
		Installed:     0,
		IsValid:       false,
		Label:         "v1x1",
		CreatedBy:     createdBy,
		SequenceId:    1,
	}

	expectedVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      payload.AppVersion,
		VersionNumber:   1.0,
		CurrentBundleId: expectedBundle.Id,
	}

	// Mock GetAppByName to return an app
	mockAppService.On("GetAppByName", ctx, payload.AppName).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndName to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndName", ctx, app.Id, payload.Environment).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndAppVersion to return nil (version doesn't exist)
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, payload.AppVersion).Return(nil, nil)

	// Mock CreateBundle to return the created bundle
	mockRepo.On("CreateBundle", ctx, mock.AnythingOfType("*model.Bundle")).Return(expectedBundle, nil)

	// Mock CreateVersion to return the created version
	mockVersionService.On("CreateVersion", ctx, mock.AnythingOfType("*model.Version")).Return(expectedVersion, nil)

	// Mock UpdateVersionIdById to return the updated bundle
	mockRepo.On("UpdateVersionIdById", ctx, expectedBundle.Id, mock.Anything).Return(expectedBundle, nil)

	// Execute
	result, err := service.CreateNewBundle(payload, createdBy)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBundle.Id, result.Id)
	assert.Equal(t, expectedBundle.AppId, result.AppId)
	assert.Equal(t, expectedBundle.EnvironmentId, result.EnvironmentId)
	assert.Equal(t, expectedBundle.DownloadFile, result.DownloadFile)
	assert.Equal(t, expectedBundle.Hash, result.Hash)
	assert.Equal(t, expectedBundle.Label, result.Label)
	assert.Equal(t, expectedBundle.CreatedBy, result.CreatedBy)

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestBundleService_CreateNewBundle_ExistingVersion_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	createdBy := "test-user"
	payload := &types.CreateNewBundleRequest{
		AppName:      "test-app",
		Environment:  "dev",
		DownloadFile: "test-bundle.js",
		Description:  "Test bundle",
		AppVersion:   "1.0.0",
		Size:         1024,
		Hash:         "test-hash",
	}

	app := &model.App{
		Id:   primitive.NewObjectID(),
		Name: payload.AppName,
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: app.Id,
		Name:  payload.Environment,
	}

	existingVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      payload.AppVersion,
		VersionNumber:   1.0,
		CurrentBundleId: primitive.NewObjectID(),
	}

	expectedBundle := &model.Bundle{
		Id:            primitive.NewObjectID(),
		AppId:         app.Id,
		EnvironmentId: environment.Id,
		DownloadFile:  payload.DownloadFile,
		Size:          payload.Size,
		Hash:          payload.Hash,
		Description:   payload.Description,
		IsMandatory:   false,
		Failed:        0,
		Installed:     0,
		IsValid:       false,
		Label:         "v1x2",
		CreatedBy:     createdBy,
		SequenceId:    2,
		VersionId:     existingVersion.Id,
	}

	// Mock GetAppByName to return an app
	mockAppService.On("GetAppByName", ctx, payload.AppName).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndName to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndName", ctx, app.Id, payload.Environment).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndAppVersion to return an existing version
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, payload.AppVersion).Return(existingVersion, nil)

	// Mock GetBundleByHashAndVersionId to return nil (no existing bundle with same hash)
	mockRepo.On("GetByHashAndVersionId", ctx, payload.Hash, existingVersion.Id).Return(nil, mongo.ErrNoDocuments)

	// Mock GetNextSeqByEnvironmentIdAndVersionId to return sequence ID
	mockRepo.On("GetNextSeqByEnvironmentIdAndVersionId", ctx, environment.Id, existingVersion.Id).Return(int64(2), nil)

	// Mock CreateBundle to return the created bundle
	mockRepo.On("CreateBundle", ctx, mock.AnythingOfType("*model.Bundle")).Return(expectedBundle, nil)

	// Mock UpdateVersionCurrentBundleIdByVersionId to return the updated version
	mockVersionService.On("UpdateVersionCurrentBundleIdByVersionId", ctx, existingVersion.Id, expectedBundle.Id).Return(existingVersion, nil)

	// Execute
	result, err := service.CreateNewBundle(payload, createdBy)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBundle.Id, result.Id)
	assert.Equal(t, expectedBundle.AppId, result.AppId)
	assert.Equal(t, expectedBundle.EnvironmentId, result.EnvironmentId)
	assert.Equal(t, expectedBundle.DownloadFile, result.DownloadFile)
	assert.Equal(t, expectedBundle.Hash, result.Hash)
	assert.Equal(t, expectedBundle.Label, result.Label)
	assert.Equal(t, expectedBundle.CreatedBy, result.CreatedBy)
	assert.Equal(t, expectedBundle.SequenceId, result.SequenceId)

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestBundleService_CreateNewBundle_AppNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	createdBy := "test-user"
	payload := &types.CreateNewBundleRequest{
		AppName:      "nonexistent-app",
		Environment:  "dev",
		DownloadFile: "test-bundle.js",
		Description:  "Test bundle",
		AppVersion:   "1.0.0",
		Size:         1024,
		Hash:         "test-hash",
	}

	// Mock GetAppByName to return nil (app not found)
	mockAppService.On("GetAppByName", ctx, payload.AppName).Return(nil, nil)

	// Execute
	result, err := service.CreateNewBundle(payload, createdBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "app not found", err.Error())

	mockAppService.AssertExpectations(t)
}

func TestBundleService_CreateNewBundle_EnvironmentNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	createdBy := "test-user"
	payload := &types.CreateNewBundleRequest{
		AppName:      "test-app",
		Environment:  "nonexistent-env",
		DownloadFile: "test-bundle.js",
		Description:  "Test bundle",
		AppVersion:   "1.0.0",
		Size:         1024,
		Hash:         "test-hash",
	}

	app := &model.App{
		Id:   primitive.NewObjectID(),
		Name: payload.AppName,
		OS:   "ios",
	}

	// Mock GetAppByName to return an app
	mockAppService.On("GetAppByName", ctx, payload.AppName).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndName to return nil (environment not found)
	mockEnvironmentService.On("GetEnvironmentByAppIdAndName", ctx, app.Id, payload.Environment).Return(nil, nil)

	// Execute
	result, err := service.CreateNewBundle(payload, createdBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "environment not found", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
}

func TestBundleService_CreateNewBundle_DuplicateHash(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	createdBy := "test-user"
	payload := &types.CreateNewBundleRequest{
		AppName:      "test-app",
		Environment:  "dev",
		DownloadFile: "test-bundle.js",
		Description:  "Test bundle",
		AppVersion:   "1.0.0",
		Size:         1024,
		Hash:         "existing-hash",
	}

	app := &model.App{
		Id:   primitive.NewObjectID(),
		Name: payload.AppName,
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    primitive.NewObjectID(),
		AppId: app.Id,
		Name:  payload.Environment,
	}

	existingVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      payload.AppVersion,
		VersionNumber:   1.0,
		CurrentBundleId: primitive.NewObjectID(),
	}

	existingBundle := &model.Bundle{
		Id:        primitive.NewObjectID(),
		Hash:      payload.Hash,
		VersionId: existingVersion.Id,
	}

	// Mock GetAppByName to return an app
	mockAppService.On("GetAppByName", ctx, payload.AppName).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndName to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndName", ctx, app.Id, payload.Environment).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndAppVersion to return an existing version
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, payload.AppVersion).Return(existingVersion, nil)

	// Mock GetBundleByHashAndVersionId to return an existing bundle (duplicate hash)
	mockRepo.On("GetByHashAndVersionId", ctx, payload.Hash, existingVersion.Id).Return(existingBundle, nil)

	// Execute
	result, err := service.CreateNewBundle(payload, createdBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "bundle with same hash already exists", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestBundleService_Rollback_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)
	environmentId, _ := primitive.ObjectIDFromHex(rollbackRequest.EnvironmentId)
	versionId, _ := primitive.ObjectIDFromHex(rollbackRequest.VersionId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    environmentId,
		AppId: appId,
		Name:  "dev",
	}

	version := &model.Version{
		Id:              versionId,
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	currentBundle := &model.Bundle{
		Id:         version.CurrentBundleId,
		SequenceId: 3,
	}

	rollbackBundle := &model.Bundle{
		Id:         primitive.NewObjectID(),
		SequenceId: 2,
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndVersionId to return a version
	mockVersionService.On("GetVersionByEnvironmentIdAndVersionId", ctx, versionId, environmentId).Return(version, nil)

	// Mock GetById to return the current bundle
	mockRepo.On("GetById", ctx, version.CurrentBundleId).Return(currentBundle, nil)

	// Mock GetBySequenceIdEnvironmentIdAndVersionId to return the rollback bundle
	mockRepo.On("GetBySequenceIdEnvironmentIdAndVersionId", ctx, int64(2), environmentId, versionId).Return(rollbackBundle, nil)

	// Mock UpdateVersionCurrentBundleIdByVersionId to return the updated version
	mockVersionService.On("UpdateVersionCurrentBundleIdByVersionId", ctx, versionId, rollbackBundle.Id).Return(version, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rollbackBundle.Id, result.Id)
	assert.Equal(t, rollbackBundle.SequenceId, result.SequenceId)

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestBundleService_Rollback_AppNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	// Mock GetAppById to return an error
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(nil, errors.New("app not found"))

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "app not found", err.Error())

	mockAppService.AssertExpectations(t)
}

func TestBundleService_Rollback_EnvironmentNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return nil (environment not found)
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(nil, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "environment not found", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
}

func TestBundleService_Rollback_InvalidVersionId(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     "invalid-version-id",
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    appId,
		AppId: appId,
		Name:  "dev",
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "error converting version id", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
}

func TestBundleService_Rollback_VersionNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)
	environmentId, _ := primitive.ObjectIDFromHex(rollbackRequest.EnvironmentId)
	versionId, _ := primitive.ObjectIDFromHex(rollbackRequest.VersionId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    environmentId,
		AppId: appId,
		Name:  "dev",
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndVersionId to return nil (version not found)
	mockVersionService.On("GetVersionByEnvironmentIdAndVersionId", ctx, versionId, environmentId).Return(nil, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "version not found", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
}

func TestBundleService_Rollback_NoCurrentBundle(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)
	environmentId, _ := primitive.ObjectIDFromHex(rollbackRequest.EnvironmentId)
	versionId, _ := primitive.ObjectIDFromHex(rollbackRequest.VersionId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    environmentId,
		AppId: appId,
		Name:  "dev",
	}

	version := &model.Version{
		Id:              versionId,
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NilObjectID, // No current bundle
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndVersionId to return a version
	mockVersionService.On("GetVersionByEnvironmentIdAndVersionId", ctx, versionId, environmentId).Return(version, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no bundle found", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
}

func TestBundleService_Rollback_BaseBundleCannotRollback(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)
	environmentId, _ := primitive.ObjectIDFromHex(rollbackRequest.EnvironmentId)
	versionId, _ := primitive.ObjectIDFromHex(rollbackRequest.VersionId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    environmentId,
		AppId: appId,
		Name:  "dev",
	}

	version := &model.Version{
		Id:              versionId,
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	currentBundle := &model.Bundle{
		Id:         version.CurrentBundleId,
		SequenceId: 1, // Base bundle sequence ID
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndVersionId to return a version
	mockVersionService.On("GetVersionByEnvironmentIdAndVersionId", ctx, versionId, environmentId).Return(version, nil)

	// Mock GetById to return the current bundle
	mockRepo.On("GetById", ctx, version.CurrentBundleId).Return(currentBundle, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "base bundle of a version cannot be rolled back", err.Error())

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestBundleService_Rollback_NoPreviousBundle(t *testing.T) {
	mockAppService := &MockAppService{}
	mockVersionService := &MockVersionService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockRepo := &MockBundleRepository{}
	service := NewBundleService(mockAppService, mockVersionService, mockEnvironmentService, mockRepo)

	ctx := context.Background()
	rollbackRequest := &types.RollbackRequest{
		AppId:         primitive.NewObjectID().Hex(),
		EnvironmentId: primitive.NewObjectID().Hex(),
		VersionId:     primitive.NewObjectID().Hex(),
	}

	appId, _ := primitive.ObjectIDFromHex(rollbackRequest.AppId)
	environmentId, _ := primitive.ObjectIDFromHex(rollbackRequest.EnvironmentId)
	versionId, _ := primitive.ObjectIDFromHex(rollbackRequest.VersionId)

	app := &model.App{
		Id:   appId,
		Name: "test-app",
		OS:   "ios",
	}

	environment := &model.Environment{
		Id:    environmentId,
		AppId: appId,
		Name:  "dev",
	}

	version := &model.Version{
		Id:              versionId,
		EnvironmentId:   environmentId,
		AppVersion:      "1.0.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	currentBundle := &model.Bundle{
		Id:         version.CurrentBundleId,
		SequenceId: 3,
	}

	// Mock GetAppById to return an app
	mockAppService.On("GetAppById", ctx, rollbackRequest.AppId).Return(app, nil)

	// Mock GetEnvironmentByAppIdAndEnvironmentId to return an environment
	mockEnvironmentService.On("GetEnvironmentByAppIdAndEnvironmentId", ctx, appId, rollbackRequest.EnvironmentId).Return(environment, nil)

	// Mock GetVersionByEnvironmentIdAndVersionId to return a version
	mockVersionService.On("GetVersionByEnvironmentIdAndVersionId", ctx, versionId, environmentId).Return(version, nil)

	// Mock GetById to return the current bundle
	mockRepo.On("GetById", ctx, version.CurrentBundleId).Return(currentBundle, nil)

	// Mock GetBySequenceIdEnvironmentIdAndVersionId to return nil (no previous bundle)
	mockRepo.On("GetBySequenceIdEnvironmentIdAndVersionId", ctx, int64(2), environmentId, versionId).Return(nil, nil)

	// Mock UpdateVersionCurrentBundleIdByVersionId to return the updated version
	mockVersionService.On("UpdateVersionCurrentBundleIdByVersionId", ctx, versionId, primitive.NilObjectID).Return(version, nil)

	// Execute
	result, err := service.Rollback(rollbackRequest)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result) // Should return nil when no previous bundle exists

	mockAppService.AssertExpectations(t)
	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
