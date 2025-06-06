package mapper

import (
	"mms/common/database/repository"
	"mms/common/role"
	"mms/common/utils"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/gql"
)

func GQLCreateUserToDomain(user *gql.CreateUserInput) domain.User {
	return domain.User{
		Username:       user.Username,
		AgeRestriction: user.AgeRestriction,
		PinHash:        user.PinHash,
		Role:           role.Role(user.Role),
	}
}

func DomainUserToCreateUserParams(user domain.User) repository.CreateUserParams {
	return repository.CreateUserParams{
		ID:             user.ID,
		Username:       user.Username,
		AgeRestriction: int64(user.AgeRestriction),
		PinHash:        utils.ToNullString(user.PinHash),
		Role:           string(user.Role),
	}
}
