package mapper

import (
	"mms/common/database/repository"
	"mms/common/role"
	"mms/common/utils"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/gql"
)

func GQLUpdateUserToDomain(user *gql.UpdateUserInput) domain.User {
	return domain.User{
		ID:             user.ID,
		Username:       user.Username,
		AgeRestriction: user.AgeRestriction,
		PinHash:        user.PinHash,
		Role:           role.Role(user.Role),
	}
}

func DomainUserToUpdateUserParams(user domain.User) repository.UpdateUserParams {
	return repository.UpdateUserParams{
		ID:             user.ID,
		Username:       user.Username,
		AgeRestriction: int64(user.AgeRestriction),
		PinHash:        utils.ToNullString(user.PinHash),
		Role:           string(user.Role),
	}
}
