package service

import (
	"context"
	"errors"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"go.uber.org/zap"
)

type ClientService interface {
	CheckUpdate(environmentKey string, appVersion string, bundleHash string) (*types.UpdateInfo, error)
	ReportStatusDeploy(reportStatusRequest *types.ReportStatusDeployRequest) error
	ReportStatusDownload(reportStatusRequest *types.ReportStatusDownloadRequest) error
}

type clientService struct {
	appService         AppService
	environmentService EnvironmentService
	bundleService      BundleService
	versionService     VersionService
}

func NewClientService(appService AppService, environmentService EnvironmentService, bundleService BundleService, versionService VersionService) ClientService {
	return &clientService{
		appService:         appService,
		environmentService: environmentService,
		bundleService:      bundleService,
		versionService:     versionService,
	}
}

// check for new update for a given environment and app version
// if there is a new update, return the update info
// if there is no update, return nil
func (s *clientService) CheckUpdate(environmentKey string, appVersion string, bundleHash string) (*types.UpdateInfo, error) {
	var updateInfo *types.UpdateInfo
	environment, err := s.environmentService.GetEnvironmentByKey(context.Background(), environmentKey)
	if err != nil {
		logger.L.Error("In CheckUpdate: Error getting environment by key", zap.String("environmentKey", environmentKey), zap.Error(err))
		return updateInfo, err
	}
	if environment == nil {
		logger.L.Error("In CheckUpdate: Environment not found", zap.String("environmentKey", environmentKey))
		return updateInfo, nil
	}

	version, err := s.versionService.GetVersionByEnvironmentIdAndAppVersion(context.Background(), environment.Id, appVersion)
	if err != nil {
		logger.L.Error("In CheckUpdate: Error getting version by environment id and app version", zap.String("environmentId", environment.Id.Hex()), zap.String("appVersion", appVersion), zap.Error(err))
		return updateInfo, err
	}
	if version == nil {
		logger.L.Error("In CheckUpdate: Version not found", zap.String("environmentId", environment.Id.Hex()), zap.String("appVersion", appVersion))
		return updateInfo, nil
	}

	bundle, err := s.bundleService.GetBundleById(version.CurrentBundleId)
	if err != nil {
		logger.L.Error("In CheckUpdate: Error getting bundle by id", zap.String("bundleId", version.CurrentBundleId.Hex()), zap.Error(err))
		return updateInfo, err
	}
	if bundle == nil {
		logger.L.Error("In CheckUpdate: Bundle not found", zap.String("bundleId", version.CurrentBundleId.Hex()))
		return updateInfo, nil
	}

	latestVersion, err := s.versionService.GetLatestVersionByEnvironmentId(context.Background(), environment.Id)
	if err != nil {
		logger.L.Error("In CheckUpdate: Error getting latest version by environment id", zap.String("environmentId", environment.Id.Hex()), zap.Error(err))
		return updateInfo, err
	}

	if bundle.Hash != bundleHash && appVersion == version.AppVersion {
		updateInfo = &types.UpdateInfo{
			DownloadUrl:            utils.GetBaseBucketUrl(config.ENV) + "/" + bundle.DownloadFile,
			Description:            bundle.Description,
			IsAvailable:            true,
			IsDisabled:             !bundle.IsValid,
			IsMandatory:            bundle.IsMandatory,
			TargetBinaryRange:      appVersion,
			PackageHash:            bundle.Hash,
			Label:                  bundle.Label,
			PackageSize:            bundle.Size,
			UpdateAppVersion:       false,
			ShouldRunBinaryVersion: false,
			Rollout:                100,
		}
	}

	// if no bundle is available for the version and there exisits a new version
	// send sdk to download the new version from store
	if updateInfo == nil && latestVersion != nil &&
		latestVersion.AppVersion != appVersion &&
		utils.FormatVersionStr(appVersion) < utils.FormatVersionStr(latestVersion.AppVersion) {
		updateInfo = &types.UpdateInfo{}
		updateInfo.TargetBinaryRange = latestVersion.AppVersion
		updateInfo.UpdateAppVersion = true
	}
	return updateInfo, nil
}

// we retrive bundle using the labelId, during checkupdate we send label as bundleId which
// the SDK sends back to report status. Use the label to fetch bundle and react to the status
func (s *clientService) ReportStatusDeploy(reportStatusRequest *types.ReportStatusDeployRequest) error {

	bundle, err := s.bundleService.GetBundleByLabel(reportStatusRequest.Label)
	if err != nil {
		logger.L.Error("In ReportStatusDeploy: Error getting bundle by id", zap.String("App Version", reportStatusRequest.AppVersion), zap.String("Deployment Key", reportStatusRequest.DeploymentKey), zap.String("Client Unique Id", reportStatusRequest.ClientUniqueId), zap.String("Label", reportStatusRequest.Label), zap.Error(err))
		return err
	}
	if bundle == nil {
		logger.L.Error("In ReportStatusDeploy: Bundle not found", zap.String("bundleLabel", reportStatusRequest.Label))
		return errors.New("bundle not found")
	}

	if reportStatusRequest.Status == "DeploymentSucceeded" {
		err = s.bundleService.AddActive(context.Background(), bundle.Id)
		if err != nil {
			logger.L.Error("In ReportStatusDeploy: Error adding active", zap.String("bundleId", bundle.Id.Hex()), zap.Error(err))
			return err
		}
	}

	if reportStatusRequest.Status == "DeploymentFailed" {
		err = s.bundleService.AddFailed(context.Background(), bundle.Id)
		if err != nil {
			logger.L.Error("In ReportStatusDeploy: Error adding failed", zap.String("bundleId", bundle.Id.Hex()), zap.Error(err))
			return err
		}
	}

	// if the previous label or app version is not nil, we need to decrement the active count of the previous bundle
	if reportStatusRequest.Status == "DeploymentSucceeded" && reportStatusRequest.PreviousLabelOrAppVersion != nil && reportStatusRequest.PreviousDeploymentKey != nil {
		logger.L.Info("In ReportStatusDeploy: Decrementing active count of previous bundle", zap.String("previousDeploymentKey", *reportStatusRequest.PreviousDeploymentKey), zap.String("previousLabelOrAppVersion", *reportStatusRequest.PreviousLabelOrAppVersion))
		previousEnvironment, err := s.environmentService.GetEnvironmentByKey(context.Background(), *reportStatusRequest.PreviousDeploymentKey)
		if err != nil {
			logger.L.Error("In ReportStatusDeploy: Error getting environment by key", zap.String("deploymentKey", *reportStatusRequest.PreviousDeploymentKey), zap.Error(err))
			return nil
		}
		if previousEnvironment == nil {
			logger.L.Error("In ReportStatusDeploy: Previous environment not found", zap.String("previousDeploymentKey", *reportStatusRequest.PreviousDeploymentKey))
			return nil
		}
		previousBundle, err := s.bundleService.GetBundleByLabel(*reportStatusRequest.PreviousLabelOrAppVersion)
		if err != nil {
			logger.L.Error("In ReportStatusDeploy: Error getting bundle by label", zap.String("bundleId", *reportStatusRequest.PreviousLabelOrAppVersion), zap.Error(err))
			return err
		}
		if previousBundle == nil {
			logger.L.Error("In ReportStatusDeploy: Previous bundle not found", zap.String("bundleLabel", *reportStatusRequest.PreviousLabelOrAppVersion))
			return errors.New("previous bundle not found")
		}
		// if previous bundle exist, then we decrement the active count of the previous bundle
		err = s.bundleService.DecrementActive(context.Background(), previousBundle.Id)
		if err != nil {
			logger.L.Error("In ReportStatusDeploy: Error decrementing active", zap.String("bundleId", previousBundle.Id.Hex()), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *clientService) ReportStatusDownload(reportStatusRequest *types.ReportStatusDownloadRequest) error {
	bundle, err := s.bundleService.GetBundleByLabel(*reportStatusRequest.Label)
	if err != nil {
		logger.L.Error("In ReportStatusDownload: Error getting bundle by id", zap.String("bundleId", *reportStatusRequest.Label), zap.Error(err))
		return err
	}
	if bundle == nil {
		logger.L.Error("In ReportStatusDownload: Bundle not found", zap.String("bundleLabel", *reportStatusRequest.Label))
		return errors.New("bundle not found")
	}
	err = s.bundleService.AddInstalled(context.Background(), bundle.Id)
	if err != nil {
		logger.L.Error("In ReportStatusDownload: Error adding install", zap.String("bundleId", bundle.Id.Hex()), zap.Error(err))
		return err
	}
	return nil
}
