package service

import (
	"context"

	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"go.uber.org/zap"
)

type ClientService interface {
	CheckUpdate(environmentKey string, appVersion string, bundleHash string) (*types.UpdateInfo, error)
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

	latestVersion, _ := s.versionService.GetLatestVersionByEnvironmentId(context.Background(), environment.Id)

	if bundle.Hash != bundleHash && appVersion == version.AppVersion {
		updateInfo = &types.UpdateInfo{
			DownloadUrl:            utils.BASE_BUCKET_URL + "/" + bundle.DownloadFile,
			Description:            bundle.Description,
			IsAvailable:            true,
			IsDisabled:             false,
			IsMandatory:            false,
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
