package category

import (
	"context"

	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

type Writer interface {
	Insert(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint64) error
}

type Reader interface {
	FindAll(ctx context.Context) ([]*entity.Category, error)
	FindByID(ctx context.Context, id uint64) (*entity.Category, error)
}

type Repor interface {
	Reader
	Writer
}
