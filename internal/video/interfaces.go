package video

import (
	"context"

	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

type Writer interface {
	Insert(ctx context.Context, video *entity.Video) error
	Update(ctx context.Context, video *entity.Video) error
	Delete(ctx context.Context, id uint64) error
}

type Reader interface {
	FindAll(ctx context.Context) ([]*entity.Video, error)
	FindByID(ctx context.Context, id uint64) (*entity.Video, error)
}

type Repor interface {
	Writer
	Reader
}

type Service interface {
	Post(ctx context.Context, video *entity.Video) error
	GetAll(ctx context.Context) ([]*entity.Video, error)
}
