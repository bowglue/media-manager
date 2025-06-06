package mapper

import (
	"mms/common/database/repository"
	"mms/common/role"
	"mms/common/utils"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/gql"
)

func SQLUserToDomain(user *repository.User) domain.User {
	return domain.User{
		ID:             user.ID,
		Username:       user.Username,
		AgeRestriction: int(user.AgeRestriction),
		PinHash:        utils.NullStringToPtr(user.PinHash),
		Role:           role.Role(user.Role),
	}
}

func SQLUsersToDomain(users []repository.User) []*domain.User {
	domainUsers := make([]*domain.User, 0, len(users))
	for _, user := range users {
		u := SQLUserToDomain(&user)
		domainUsers = append(domainUsers, &u)
	}
	return domainUsers
}

func DomainUserToGQL(user *domain.User) *gql.User {
	return &gql.User{
		ID:             user.ID,
		Username:       user.Username,
		AgeRestriction: user.AgeRestriction,
		Role:           string(user.Role),
	}
}

func DomainUsersToGQL(users []*domain.User) []*gql.User {
	gqlUsers := make([]*gql.User, 0, len(users))
	for _, user := range users {
		gqlUsers = append(gqlUsers, DomainUserToGQL(user))
	}
	return gqlUsers
}
