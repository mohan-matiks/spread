package types

type CreateEnvironmentRequest struct {
	EnvironmentName string `json:"environmentName" validate:"required"`
	AppName         string `json:"appName" validate:"required"`
}
