package services

import (
	"context"
	"database/sql"
	"errors"
	"mms/modules/user-module/domain"
	"mms/modules/user-module/errs"
	"mms/modules/user-module/mapper"
)

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	domainUser := mapper.SQLUserToDomain(&user)
	return &domainUser, nil
}
