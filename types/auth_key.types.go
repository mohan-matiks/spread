package types

type CreateAuthKeyRequest struct {
	Name string `json:"name" validate:"required"`
}
