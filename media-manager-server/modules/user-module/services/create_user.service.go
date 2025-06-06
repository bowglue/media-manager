package services

import (
	"context"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/errs"
	"mms/modules/user-module/mapper"

	"github.com/google/uuid"
)

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	usernameTaken, err := s.queries.DoesUsernameExist(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if usernameTaken != 0 {
		return nil, errs.ErrUsernameAlreadyExists
	}

	user.ID = uuid.NewString()

	params := mapper.DomainUserToCreateUserParams(user)
	createdUser, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	domainUser := mapper.SQLUserToDomain(&createdUser)
	return &domainUser, nil
}
