package services

import (
	"mms/common/database/repository"
)

type MediaService struct {
	queries *repository.Queries
}

func NewMediaService(queries *repository.Queries) *MediaService {
	return &MediaService{
		queries: queries,
	}
}

// func (s *MediaService) GetMovie(id string) (*gql.Movie, error) {
// 	return &gql.Movie{
// 		ID:    id,
// 		Title: "media_" + id,
// 	}, nil
// }
