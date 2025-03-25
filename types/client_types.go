package types

type UpdateInfo struct {
	DownloadUrl            string `json:"download_url"`
	Description            string `json:"description"`
	IsAvailable            bool   `json:"is_available"`
	IsDisabled             bool   `json:"is_disabled"`
	TargetBinaryRange      string `json:"target_binary_range"`
	PackageHash            string `json:"package_hash"`
	Label                  string `json:"label"`
	PackageSize            int64  `json:"package_size"`
	UpdateAppVersion       bool   `json:"update_app_version"`
	ShouldRunBinaryVersion bool   `json:"should_run_binary_version"`
	IsMandatory            bool   `json:"is_mandatory"`
	Rollout                int    `json:"rollout"`
}

type ReportStatusDeployRequest struct {
	AppVersion                string  `json:"app_version"`
	DeploymentKey             string  `json:"deployment_key"`
	ClientUniqueId            string  `json:"client_unique_id"`
	Label                     string  `json:"label"`
	Status                    string  `json:"status"`
	PreviousLabelOrAppVersion *string `json:"previous_label_or_app_version"`
	PreviousDeploymentKey     *string `json:"previous_deployment_key"`
}

type ReportStatusDownloadRequest struct {
	ClientUniqueId string `json:"client_unique_id"`
	DeploymentKey  string `json:"deployment_key"`
	Label          string `json:"label"`
}
