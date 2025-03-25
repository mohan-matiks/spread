package service

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/pkg"
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type BundleService interface {
	UploadBundle(fileName string, file *multipart.FileHeader) error
	Rollback(rollbackRequest *types.RollbackRequest) (*model.Bundle, error)
	CreateNewBundle(createNewBundleRequest *types.CreateNewBundleRequest, createdBy string) (*model.Bundle, error)
	GetBundleById(id primitive.ObjectID) (*model.Bundle, error)
	GetBundleByLabel(label string) (*model.Bundle, error)
	GetBundlesByVersionId(versionId primitive.ObjectID) ([]*model.Bundle, error)
	ToggleMandatory(bundleId primitive.ObjectID) error
	ToggleActive(bundleId primitive.ObjectID) error
	AddActive(ctx context.Context, id primitive.ObjectID) error
	AddFailed(ctx context.Context, id primitive.ObjectID) error
	AddInstalled(ctx context.Context, id primitive.ObjectID) error
	DecrementActive(ctx context.Context, id primitive.ObjectID) error
}

type bundleService struct {
	appService         AppService
	versionService     VersionService
	environmentService EnvironmentService

	bundleRepository repository.BundleRepository
}

func NewBundleService(appService AppService, versionService VersionService, environmentService EnvironmentService, bundleRepository repository.BundleRepository) BundleService {
	return &bundleService{appService: appService, versionService: versionService, environmentService: environmentService, bundleRepository: bundleRepository}
}

func (bundleService *bundleService) UploadBundle(fileName string, file *multipart.FileHeader) error {
	r2Service, err := pkg.NewR2Service()
	if err != nil {
		return err
	}

	fileBytes, err := file.Open()
	if err != nil {
		return err
	}
	defer fileBytes.Close()

	// Read file into byte slice
	buffer, err := io.ReadAll(fileBytes)
	if err != nil {
		return err
	}

	err = r2Service.UploadFileToR2(context.Background(), fileName, buffer)
	if err != nil {
		return err
	}
	return nil
}

// we check if a version (0.0.1) exists, if it does then we create a new bundle and set the version id to the bundle
// if it doesn't exist then we create a new version and set the bundle id to the version
func (bundleService *bundleService) CreateNewBundle(payload *types.CreateNewBundleRequest, createdBy string) (*model.Bundle, error) {
	// Retrieve the app by name
	app, err := bundleService.appService.GetAppByName(context.Background(), payload.AppName)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("app not found")
	}
	// Retrieve the environment by app ID and name
	environment, err := bundleService.environmentService.GetEnvironmentByAppIdAndName(context.Background(), app.Id, payload.Environment)
	if err != nil {
		return nil, err
	}
	if environment == nil {
		return nil, errors.New("environment not found")
	}
	// Retrieve the version by environment ID and app version
	version, err := bundleService.versionService.GetVersionByEnvironmentIdAndAppVersion(context.Background(), environment.Id, payload.AppVersion)
	if err != nil {
		return nil, err
	}
	logger.L.Info("In CreateNewBundle: Version found", zap.Any("version", version))
	// If version is not found, create a new bundle and version
	if version == nil {
		versionNumber := utils.FormatVersionStr(payload.AppVersion)
		bundle := &model.Bundle{
			AppId:         app.Id,
			EnvironmentId: environment.Id,
			DownloadFile:  payload.DownloadFile,
			Size:          payload.Size,
			Hash:          payload.Hash,
			Description:   payload.Description,
			IsMandatory:   false,
			Failed:        0,
			Installed:     0,
			IsValid:       true,
			Label:         "v" + strconv.Itoa(int(versionNumber)) + ":" + strconv.Itoa(1),
			CreatedBy:     createdBy,
			SequenceId:    1,
		}

		bundle, err = bundleService.bundleRepository.CreateBundle(context.Background(), bundle)
		if err != nil {
			return nil, err
		}
		version = &model.Version{
			EnvironmentId:   environment.Id,
			AppVersion:      payload.AppVersion,
			VersionNumber:   versionNumber,
			CurrentBundleId: bundle.Id,
		}
		_, err = bundleService.versionService.CreateVersion(context.Background(), version)
		if err != nil {
			return nil, err
		}
		// Set the version ID to the bundle
		bundle.VersionId = version.Id
		_, err = bundleService.bundleRepository.UpdateVersionIdById(context.Background(), bundle.Id, version.Id)
		if err != nil {
			return nil, err
		}
		return bundle, nil
	}
	// If version exists, check if a bundle with the same hash already exists
	existingBundle, err := bundleService.GetBundleByHash(payload.Hash)
	if err != nil {
		return nil, err
	}
	if existingBundle != nil {
		return nil, errors.New("bundle with same hash already exists")
	}

	sequenceId, err := bundleService.bundleRepository.GetNextSeqByEnvironmentIdAndVersionId(context.Background(), environment.Id, version.Id)
	if err != nil {
		return nil, err
	}
	// Create a new bundle and set it to the version
	bundle := &model.Bundle{
		AppId:         app.Id,
		EnvironmentId: environment.Id,
		DownloadFile:  payload.DownloadFile,
		Size:          payload.Size,
		Hash:          payload.Hash,
		Description:   payload.Description,
		SequenceId:    sequenceId,
		VersionId:     version.Id,
		IsMandatory:   false,
		Failed:        0,
		Installed:     0,
		Label:         "v" + strconv.Itoa(int(version.VersionNumber)) + ":" + strconv.Itoa(int(sequenceId)),
		IsValid:       true,
	}
	bundle, err = bundleService.bundleRepository.CreateBundle(context.Background(), bundle)
	if err != nil {
		return nil, err
	}
	version.CurrentBundleId = bundle.Id
	_, err = bundleService.versionService.UpdateVersionCurrentBundleIdByVersionId(context.Background(), version.Id, bundle.Id)
	if err != nil {
		return nil, err
	}
	return bundle, nil
}

// Rollback is essentially changing the bundle of a version to the previous bundle if any exists
func (bundleService *bundleService) Rollback(rollbackRequest *types.RollbackRequest) (*model.Bundle, error) {
	app, err := bundleService.appService.GetAppById(context.Background(), rollbackRequest.AppId)
	if err != nil {
		return nil, err
	}
	logger.L.Info("In Rollback: App found", zap.Any("app", app))
	environment, err := bundleService.environmentService.GetEnvironmentByAppIdAndEnvironmentId(context.Background(), app.Id, rollbackRequest.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if environment == nil {
		return nil, errors.New("environment not found")
	}
	logger.L.Info("In Rollback: Environment found", zap.Any("environment", environment))
	versionId, err := primitive.ObjectIDFromHex(rollbackRequest.VersionId)
	if err != nil {
		return nil, errors.New("error converting version id")
	}
	version, err := bundleService.versionService.GetVersionByEnvironmentIdAndVersionId(context.Background(), versionId, environment.Id)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("version not found")
	}
	logger.L.Info("In Rollback: Version found", zap.Any("version", version))
	// if no current bundle is set for the version, then return an error
	if version.CurrentBundleId == primitive.NilObjectID {
		return nil, errors.New("no bundle found")
	}
	logger.L.Info("In Rollback: Version found", zap.Any("version", version))
	bundle, err := bundleService.bundleRepository.GetById(context.Background(), version.CurrentBundleId)
	if err != nil {
		logger.L.Error("In Rollback: Error getting current bundle", zap.Error(err))
		return nil, err
	}
	// using the current bundle of version, get the previous bundle
	// every bundle of a version has a sequence id, so we get the previous bundle by subtracting 1 from the current bundle's sequence id
	rollbackBundle, err := bundleService.bundleRepository.GetBySequenceIdEnvironmentIdAndVersionId(context.Background(), bundle.SequenceId-1, environment.Id, version.Id)
	if err != nil {
		logger.L.Error("In Rollback: Error getting previous bundle", zap.Error(err))
		return nil, err
	}
	// if reollback bundle which is the previous bundle to rollback to is not present
	if rollbackBundle == nil {
		version.CurrentBundleId = primitive.NilObjectID
		_, err = bundleService.versionService.UpdateVersionCurrentBundleIdByVersionId(context.Background(), version.Id, primitive.NilObjectID)
		if err != nil {
			logger.L.Error("In Rollback: Error updating version current bundle id with nil", zap.Error(err))
			return nil, err
		}
		return nil, nil
	}
	version.CurrentBundleId = rollbackBundle.Id
	_, err = bundleService.versionService.UpdateVersionCurrentBundleIdByVersionId(context.Background(), version.Id, rollbackBundle.Id)
	if err != nil {
		logger.L.Error("In Rollback: Error updating version current bundle id", zap.Error(err))
		return nil, err
	}
	return rollbackBundle, nil
}

func (bundleService *bundleService) GetBundleByLabel(label string) (*model.Bundle, error) {
	bundle, err := bundleService.bundleRepository.GetByLabel(context.Background(), label)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bundle, nil
}

func (bundleService *bundleService) GetBundlesByVersionId(versionId primitive.ObjectID) ([]*model.Bundle, error) {
	bundles, err := bundleService.bundleRepository.GetAllByVersionId(context.Background(), versionId)
	if err != nil {
		return nil, err
	}
	// loop through bundles append base bucket url to downloadUrl
	for i, bundle := range bundles {
		bundles[i].DownloadFile = utils.GetBaseBucketUrl(config.ENV) + "/" + bundle.DownloadFile
	}
	return bundles, nil
}

func (bundleService *bundleService) ToggleMandatory(bundleId primitive.ObjectID) error {
	bundle, err := bundleService.bundleRepository.GetById(context.Background(), bundleId)
	if err != nil {
		return err
	}
	bundle.IsMandatory = !bundle.IsMandatory
	_, err = bundleService.bundleRepository.UpdateIsMandatoryById(context.Background(), bundleId, bundle.IsMandatory)
	if err != nil {
		return err
	}
	return nil
}

func (bundleService *bundleService) ToggleActive(bundleId primitive.ObjectID) error {
	bundle, err := bundleService.bundleRepository.GetById(context.Background(), bundleId)
	if err != nil {
		return err
	}
	bundle.IsValid = !bundle.IsValid
	_, err = bundleService.bundleRepository.UpdateIsValid(context.Background(), bundleId, bundle.IsValid)
	if err != nil {
		return err
	}
	return nil
}

func (bundleService *bundleService) GetBundleByHash(hash string) (*model.Bundle, error) {
	bundle, err := bundleService.bundleRepository.GetByHash(context.Background(), hash)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bundle, nil
}

func (bundleService *bundleService) GetBundleById(id primitive.ObjectID) (*model.Bundle, error) {
	bundle, err := bundleService.bundleRepository.GetById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return bundle, nil
}

func (bundleService *bundleService) AddActive(ctx context.Context, id primitive.ObjectID) error {
	return bundleService.bundleRepository.AddActive(ctx, id)
}

func (bundleService *bundleService) AddFailed(ctx context.Context, id primitive.ObjectID) error {
	return bundleService.bundleRepository.AddFailed(ctx, id)
}

func (bundleService *bundleService) AddInstalled(ctx context.Context, id primitive.ObjectID) error {
	return bundleService.bundleRepository.AddInstalled(ctx, id)
}

func (bundleService *bundleService) DecrementActive(ctx context.Context, id primitive.ObjectID) error {
	return bundleService.bundleRepository.DecrementActive(ctx, id)
}
