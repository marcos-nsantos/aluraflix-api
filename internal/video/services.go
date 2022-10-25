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

func (s *service) GetByID(ctx context.Context, id uint64) (*entity.Video, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, video *entity.Video, id uint64) error {
	video.ID = id
	return s.repo.Update(ctx, video)
}
