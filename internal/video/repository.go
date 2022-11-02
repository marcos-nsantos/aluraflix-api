package video

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) Insert(ctx context.Context, video *entity.Video) error {
	query := `INSERT INTO videos (title, description, url, category_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, video.Title, video.Description, video.URL, video.CategoryID)
	if err = row.Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (r Repository) FindAll(ctx context.Context) ([]*entity.Video, error) {
	query := `SELECT id, title, description, url, category_id, created_at, updated_at FROM videos WHERE deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		var video entity.Video
		if err = rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CategoryID, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (r Repository) FindByID(ctx context.Context, id uint64) (*entity.Video, error) {
	query := `SELECT id, title, description, url, category_id, created_at, updated_at FROM videos WHERE id = $1 AND deleted_at IS NULL`

	var video entity.Video
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CategoryID, &video.CreatedAt, &video.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrVideoNotFound
		}
	}

	return &video, nil
}

func (r Repository) Update(ctx context.Context, video *entity.Video) error {
	query := `UPDATE videos SET title = $1, description = $2, url = $3, category_id = $4, updated_at = NOW() WHERE id = $5 AND deleted_at IS NULL`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, video.Title, video.Description, video.URL, video.CategoryID, video.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrVideoNotFound
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, id uint64) error {
	query := `UPDATE videos SET deleted_at = NOW() WHERE id = $1`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrVideoNotFound
	}

	return nil
}
