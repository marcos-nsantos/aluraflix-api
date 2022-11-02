package video

import (
	"context"

	"github.com/marcos-nsantos/aluraflix-api/internal/category"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

func NewService(videoRepor Repor, categoryRepor category.Repor) Service {
	return &service{
		videoRepor:    videoRepor,
		categoryRepor: categoryRepor,
	}
}

type service struct {
	videoRepor    Repor
	categoryRepor category.Repor
}

func (s *service) Post(ctx context.Context, video *entity.Video) error {
	if _, err := s.categoryRepor.FindByID(ctx, video.CategoryID); err != nil {
		return err
	}
	return s.videoRepor.Insert(ctx, video)
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Video, error) {
	return s.videoRepor.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id uint64) (*entity.Video, error) {
	return s.videoRepor.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, video *entity.Video, id uint64) error {
	if _, err := s.categoryRepor.FindByID(ctx, video.CategoryID); err != nil {
		return err
	}
	video.ID = id
	return s.videoRepor.Update(ctx, video)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.videoRepor.Delete(ctx, id)
}
