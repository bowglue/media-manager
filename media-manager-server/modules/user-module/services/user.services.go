package services

import (
	"context"
	models "mms/common/graphql/generated"

	// "mms/modules/user-module/repository"
	"mms/common/database/repository"

	"github.com/google/uuid"
)

type UserService struct {
	queries *repository.Queries
}

func NewUserService(queries *repository.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	ctx := context.Background()
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *UserService) CreateUser(input models.CreateUserInput) (*models.User, error) {
	ctx := context.Background()
	id := uuid.NewString()

	params := repository.CreateUserParams{
		ID:       id,
		Username: input.Username,
	}

	err := s.queries.CreateUser(ctx, params)

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       id,
		Username: input.Username,
	}, nil
}

func (s *UserService) GetUsers() ([]*models.User, error) {
	return []*models.User{
		{ID: "1", Username: "user_1"},
		{ID: "2", Username: "user_2"},
	}, nil
}
