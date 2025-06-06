package validations

import "mms/modules/user-module/gql"

func ValidateUpdateUserInput(input gql.UpdateUserInput) error {
	if err := validateUsername(input.Username); err != nil {
		return err
	}
	if err := validateRole(input.Role); err != nil {
		return err
	}
	return nil
}
