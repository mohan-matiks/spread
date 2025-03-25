package types

type CreateNewBundleRequest struct {
	AppName      string `json:"appName" validate:"required"`
	Environment  string `json:"environment" validate:"required"`
	DownloadFile string `json:"downloadFile" validate:"required"`
	Description  string `json:"description"`
	AppVersion   string `json:"appVersion" validate:"required"`
	Size         int64  `json:"size" validate:"required"`
	Hash         string `json:"hash" validate:"required"`
}

type RollbackRequest struct {
	AppId         string `json:"appId" validate:"required"`
	EnvironmentId string `json:"environmentId" validate:"required"`
	VersionId     string `json:"versionId" validate:"required"`
}
