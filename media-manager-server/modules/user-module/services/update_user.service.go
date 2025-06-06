package services

import (
	"context"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/mapper"
)

func (s *UserService) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	_, err := s.GetUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	params := mapper.DomainUserToUpdateUserParams(user)
	updatedUser, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	domainUser := mapper.SQLUserToDomain(&updatedUser)
	return &domainUser, nil
}
