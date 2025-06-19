package service

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBundleService is a mock implementation of BundleService
type MockBundleService struct {
	mock.Mock
}

func (m *MockBundleService) UploadBundle(fileName string, file *multipart.FileHeader) error {
	args := m.Called(fileName, file)
	return args.Error(0)
}

func (m *MockBundleService) Rollback(rollbackRequest *types.RollbackRequest) (*model.Bundle, error) {
	args := m.Called(rollbackRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleService) CreateNewBundle(createNewBundleRequest *types.CreateNewBundleRequest, createdBy string) (*model.Bundle, error) {
	args := m.Called(createNewBundleRequest, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleService) GetBundleById(id primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleService) GetBundleByLabelAndEnvironmentId(label string, environmentId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(label, environmentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleService) GetBundleByHashAndVersionId(hash string, versionId primitive.ObjectID) (*model.Bundle, error) {
	args := m.Called(hash, versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bundle), args.Error(1)
}

func (m *MockBundleService) GetBundlesByVersionId(versionId primitive.ObjectID) ([]*model.Bundle, error) {
	args := m.Called(versionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Bundle), args.Error(1)
}

func (m *MockBundleService) ToggleMandatory(bundleId primitive.ObjectID) error {
	args := m.Called(bundleId)
	return args.Error(0)
}

func (m *MockBundleService) ToggleActive(bundleId primitive.ObjectID) error {
	args := m.Called(bundleId)
	return args.Error(0)
}

func (m *MockBundleService) AddActive(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleService) AddFailed(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleService) AddInstalled(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBundleService) DecrementActive(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewClientService(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}

	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	assert.NotNil(t, service)
	assert.IsType(t, &clientService{}, service)
}

func TestClientService_CheckUpdate_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "test-env-key"
	appVersion := "1.0.0"
	bundleHash := "old-hash"

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: environmentKey,
	}

	version := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      appVersion,
		CurrentBundleId: primitive.NewObjectID(),
	}

	bundle := &model.Bundle{
		Id:           version.CurrentBundleId,
		DownloadFile: "test-bundle.js",
		Hash:         "new-hash",
		Description:  "Test bundle",
		IsMandatory:  false,
		IsValid:      true,
		Size:         1024,
		Label:        "v1x1",
	}

	latestVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      "1.1.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(environment, nil)
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, appVersion).Return(version, nil)
	mockBundleService.On("GetBundleById", version.CurrentBundleId).Return(bundle, nil)
	mockVersionService.On("GetLatestVersionByEnvironmentId", ctx, environment.Id).Return(latestVersion, nil)

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bundle.Hash, result.PackageHash)
	assert.Equal(t, bundle.Description, result.Description)
	assert.Equal(t, bundle.IsMandatory, result.IsMandatory)
	assert.Equal(t, bundle.IsValid, result.IsAvailable)
	assert.Equal(t, appVersion, result.TargetBinaryRange)
	assert.False(t, result.UpdateAppVersion)

	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_CheckUpdate_EnvironmentNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "nonexistent-key"
	appVersion := "1.0.0"
	bundleHash := "old-hash"

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(nil, nil)

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockEnvironmentService.AssertExpectations(t)
}

func TestClientService_CheckUpdate_EnvironmentError(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "test-env-key"
	appVersion := "1.0.0"
	bundleHash := "old-hash"

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(nil, errors.New("database error"))

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockEnvironmentService.AssertExpectations(t)
}

func TestClientService_CheckUpdate_VersionNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "test-env-key"
	appVersion := "1.0.0"
	bundleHash := "old-hash"

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: environmentKey,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(environment, nil)
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, appVersion).Return(nil, nil)

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
}

func TestClientService_CheckUpdate_BundleNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "test-env-key"
	appVersion := "1.0.0"
	bundleHash := "old-hash"

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: environmentKey,
	}

	version := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      appVersion,
		CurrentBundleId: primitive.NewObjectID(),
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(environment, nil)
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, appVersion).Return(version, nil)
	mockBundleService.On("GetBundleById", version.CurrentBundleId).Return(nil, nil)

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_CheckUpdate_NoUpdateNeeded(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	environmentKey := "test-env-key"
	appVersion := "1.0.0"
	bundleHash := "current-hash"

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: environmentKey,
	}

	version := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      appVersion,
		CurrentBundleId: primitive.NewObjectID(),
	}

	bundle := &model.Bundle{
		Id:           version.CurrentBundleId,
		DownloadFile: "test-bundle.js",
		Hash:         bundleHash, // Same hash
		Description:  "Test bundle",
		IsMandatory:  false,
		IsValid:      true,
		Size:         1024,
		Label:        "v1x1",
	}

	latestVersion := &model.Version{
		Id:              primitive.NewObjectID(),
		EnvironmentId:   environment.Id,
		AppVersion:      "1.1.0",
		CurrentBundleId: primitive.NewObjectID(),
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, environmentKey).Return(environment, nil)
	mockVersionService.On("GetVersionByEnvironmentIdAndAppVersion", ctx, environment.Id, appVersion).Return(version, nil)
	mockBundleService.On("GetBundleById", version.CurrentBundleId).Return(bundle, nil)
	mockVersionService.On("GetLatestVersionByEnvironmentId", ctx, environment.Id).Return(latestVersion, nil)

	result, err := service.CheckUpdate(environmentKey, appVersion, bundleHash)

	assert.NoError(t, err)
	assert.NotNil(t, result) // Should return an update to prompt app version update since there's a newer version
	assert.Equal(t, "1.1.0", result.TargetBinaryRange)
	assert.True(t, result.UpdateAppVersion)

	mockEnvironmentService.AssertExpectations(t)
	mockVersionService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_ReportStatusDeploy_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "test-env-key"
	label := "v1x1"

	request := &types.ReportStatusDeployRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
		Status:        "DeploymentSucceeded",
	}

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: deploymentKey,
	}

	bundle := &model.Bundle{
		Id:    primitive.NewObjectID(),
		Label: label,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(environment, nil)
	mockBundleService.On("GetBundleByLabelAndEnvironmentId", label, environment.Id).Return(bundle, nil)
	mockBundleService.On("AddActive", ctx, bundle.Id).Return(nil)

	err := service.ReportStatusDeploy(request)

	assert.NoError(t, err)

	mockEnvironmentService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_ReportStatusDeploy_EnvironmentNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "nonexistent-key"
	label := "v1x1"

	request := &types.ReportStatusDeployRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
		Status:        "DeploymentSucceeded",
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(nil, nil)

	err := service.ReportStatusDeploy(request)

	assert.Error(t, err)
	assert.Equal(t, "environment not found", err.Error())

	mockEnvironmentService.AssertExpectations(t)
}

func TestClientService_ReportStatusDeploy_BundleNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "test-env-key"
	label := "v1x1"

	request := &types.ReportStatusDeployRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
		Status:        "DeploymentSucceeded",
	}

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: deploymentKey,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(environment, nil)
	mockBundleService.On("GetBundleByLabelAndEnvironmentId", label, environment.Id).Return(nil, nil)

	err := service.ReportStatusDeploy(request)

	assert.Error(t, err)
	assert.Equal(t, "bundle not found", err.Error())

	mockEnvironmentService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_ReportStatusDeploy_DeploymentFailed(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "test-env-key"
	label := "v1x1"

	request := &types.ReportStatusDeployRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
		Status:        "DeploymentFailed",
	}

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: deploymentKey,
	}

	bundle := &model.Bundle{
		Id:    primitive.NewObjectID(),
		Label: label,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(environment, nil)
	mockBundleService.On("GetBundleByLabelAndEnvironmentId", label, environment.Id).Return(bundle, nil)
	mockBundleService.On("AddFailed", ctx, bundle.Id).Return(nil)

	err := service.ReportStatusDeploy(request)

	assert.NoError(t, err)

	mockEnvironmentService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_ReportStatusDownload_Success(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "test-env-key"
	label := "v1x1"

	request := &types.ReportStatusDownloadRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
	}

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: deploymentKey,
	}

	bundle := &model.Bundle{
		Id:    primitive.NewObjectID(),
		Label: label,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(environment, nil)
	mockBundleService.On("GetBundleByLabelAndEnvironmentId", label, environment.Id).Return(bundle, nil)
	mockBundleService.On("AddInstalled", ctx, bundle.Id).Return(nil)

	err := service.ReportStatusDownload(request)

	assert.NoError(t, err)

	mockEnvironmentService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}

func TestClientService_ReportStatusDownload_EnvironmentNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "nonexistent-key"
	label := "v1x1"

	request := &types.ReportStatusDownloadRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(nil, nil)

	err := service.ReportStatusDownload(request)

	assert.Error(t, err)
	assert.Equal(t, "environment not found", err.Error())

	mockEnvironmentService.AssertExpectations(t)
}

func TestClientService_ReportStatusDownload_BundleNotFound(t *testing.T) {
	mockAppService := &MockAppService{}
	mockEnvironmentService := &MockEnvironmentService{}
	mockBundleService := &MockBundleService{}
	mockVersionService := &MockVersionService{}
	service := NewClientService(mockAppService, mockEnvironmentService, mockBundleService, mockVersionService)

	ctx := context.Background()
	deploymentKey := "test-env-key"
	label := "v1x1"

	request := &types.ReportStatusDownloadRequest{
		DeploymentKey: deploymentKey,
		Label:         label,
	}

	environment := &model.Environment{
		Id:  primitive.NewObjectID(),
		Key: deploymentKey,
	}

	mockEnvironmentService.On("GetEnvironmentByKey", ctx, deploymentKey).Return(environment, nil)
	mockBundleService.On("GetBundleByLabelAndEnvironmentId", label, environment.Id).Return(nil, nil)

	err := service.ReportStatusDownload(request)

	assert.Error(t, err)
	assert.Equal(t, "bundle not found", err.Error())

	mockEnvironmentService.AssertExpectations(t)
	mockBundleService.AssertExpectations(t)
}
