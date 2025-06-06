package services

import (
	"context"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/mapper"
)

func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	domainUsers := mapper.SQLUsersToDomain(users)
	return domainUsers, nil
}
