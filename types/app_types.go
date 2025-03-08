package types

type CreateAppRequest struct {
	AppName string `json:"appName" validate:"required"`
	OS      string `json:"os" validate:"required"`
}
