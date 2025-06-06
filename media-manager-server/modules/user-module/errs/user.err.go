package errs

import (
	"mms/common/errs"
	"mms/common/role"
)

var (
	ErrUsernameAlreadyExists = errs.New("username already exists")
	ErrUsernameIsRequired    = errs.New("username is required")
	ErrInvalidRole           = errs.New("role must be one of: " + role.AllowedRolesString())
	ErrUserNotFound          = errs.New("user not found")
)
