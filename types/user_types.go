package types

type CreateUserRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
