package validations

import (
	"mms/common/role"
	"mms/modules/user-module/errs"
)

func validateUsername(username string) error {
	if username == "" {
		return errs.ErrUsernameIsRequired
	}
	return nil
}

func validateRole(value string) error {
	if !role.IsValid(value) {
		return errs.ErrInvalidRole
	}
	return nil
}
