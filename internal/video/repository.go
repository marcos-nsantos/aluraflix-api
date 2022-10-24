package video

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) Create(ctx context.Context, video *entity.Video) error {
	query := `INSERT INTO videos (title, description, url) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, video.Title, video.Description, video.URL)
	if err = row.Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt); err != nil {
		return err
	}

	return nil
}