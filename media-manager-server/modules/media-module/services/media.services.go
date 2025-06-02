package services

import (
	"mms/common/database/repository"
	models "mms/common/graphql/generated"
)

type MediaService struct {
	queries *repository.Queries
}

func NewMediaService(queries *repository.Queries) *MediaService {
	return &MediaService{
		queries: queries,
	}
}

func (s *MediaService) GetMovie(id string) (*models.Movie, error) {
	return &models.Movie{
		ID:    id,
		Title: "media_" + id,
	}, nil
}
