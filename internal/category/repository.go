package category

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

func (r *Repository) Insert(ctx context.Context, video *entity.Category) error {
	query := `INSERT INTO categories (title, color) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, video.Title, video.Color)
	if err = row.Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	query := `SELECT id, title, color, created_at, updated_at FROM categories WHERE deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		var category entity.Category
		if err = rows.Scan(&category.ID, &category.Title, &category.Color, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *Repository) FindByID(ctx context.Context, id uint64) (*entity.Category, error) {
	query := `SELECT id, title, color, created_at, updated_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	var category entity.Category
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&category.ID, &category.Title, &category.Color, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrCategoryNotFound
		}
	}

	return &category, nil
}

func (r *Repository) Update(ctx context.Context, category *entity.Category) error {
	query := `UPDATE categories SET title = $1, color = $2, updated_at = NOW() WHERE id = $3`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, category.Title, category.Color, category.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrCategoryNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := `UPDATE categories SET deleted_at = NOW() WHERE id = $1`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrCategoryNotFound
	}

	return nil
}
