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
	AppName     string `json:"appName" validate:"required"`
	Environment string `json:"environment" validate:"required"`
	AppVersion  string `json:"appVersion" validate:"required"`
}
