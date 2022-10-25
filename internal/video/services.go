package video

import (
	"context"

	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

func NewService(repo Repor) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repor
}

func (s *service) Post(ctx context.Context, video *entity.Video) error {
	return s.repo.Insert(ctx, video)
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Video, error) {
	return s.repo.FindAll(ctx)
}
