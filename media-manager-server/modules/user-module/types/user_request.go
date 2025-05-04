package types

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
}
