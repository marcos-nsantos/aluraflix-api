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

func (r Repository) FindAll(ctx context.Context) ([]*entity.Video, error) {
	query := `SELECT id, title, description, url, created_at, updated_at FROM videos`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		var video entity.Video
		if err = rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (r Repository) FindByID(ctx context.Context, id uint64) (*entity.Video, error) {
	query := `SELECT id, title, description, url, created_at, updated_at FROM videos WHERE id = $1`

	var video entity.Video
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CreatedAt, &video.UpdatedAt); err != nil {
		return nil, err
	}

	return &video, nil
}
