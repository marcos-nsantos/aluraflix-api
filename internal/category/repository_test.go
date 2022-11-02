//go:build integration

package category

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/marcos-nsantos/aluraflix-api/internal/database"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

var dbConn *sqlx.DB

func TestMain(t *testing.M) {
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
	}
	dbConn = db.Client

	err = db.Migrate("file://../..//migrations")
	if err != nil {
		fmt.Println(err)
	}

	code := t.Run()

	os.Exit(code)
}

func TestInsertCategory(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should insert a category", func(t *testing.T) {
		category := &entity.Category{
			Title: "Gopher Blue",
			Color: "#00ADD8",
		}

		err := repo.Insert(context.Background(), category)
		assert.NoError(t, err)
		assert.NotEmpty(t, category.ID)
		assert.NotEmpty(t, category.CreatedAt)
		assert.NotEmpty(t, category.UpdatedAt)
	})
}

func TestFindAllCategories(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should find all categories", func(t *testing.T) {
		categories := []*entity.Category{
			{
				Title: "Gopher Blue",
				Color: "#00ADD8",
			},
			{
				Title: "Green",
				Color: "#00FF00",
			},
		}

		for _, category := range categories {
			err := repo.Insert(context.Background(), category)
			assert.NoError(t, err)
		}

		categories, err := repo.FindAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, categories)
		assert.GreaterOrEqual(t, len(categories), 2)
	})
}

func TestFindByIDCategory(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should find a category by id", func(t *testing.T) {
		category := &entity.Category{
			Title: "Gopher Blue",
			Color: "#00ADD8",
		}

		err := repo.Insert(context.Background(), category)
		assert.NoError(t, err)

		category, err = repo.FindByID(context.Background(), category.ID)
		assert.NoError(t, err)
		assert.Equal(t, category.Title, "Gopher Blue")
		assert.Equal(t, category.Color, "#00ADD8")
		assert.NotEmpty(t, category.CreatedAt)
		assert.NotEmpty(t, category.UpdatedAt)
	})

	t.Run("should not find a category by id", func(t *testing.T) {
		category, err := repo.FindByID(context.Background(), 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrCategoryNotFound)
		assert.Nil(t, category)
	})
}

func TestUpdateCategory(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should update a category", func(t *testing.T) {
		category := &entity.Category{
			Title: "Gopher Blue",
			Color: "#00ADD8",
		}
		err := repo.Insert(context.Background(), category)
		assert.NoError(t, err)
		categoryUpdatedAtFirst := category.UpdatedAt

		category.Title = "Gopher Green"
		category.Color = "#00FF00"

		err = repo.Update(context.Background(), category)
		assert.NoError(t, err)

		category, err = repo.FindByID(context.Background(), category.ID)
		assert.NoError(t, err)
		assert.Equal(t, category.Title, "Gopher Green")
		assert.Equal(t, category.Color, "#00FF00")
		assert.NotEmpty(t, category.CreatedAt)
		assert.NotEmpty(t, category.UpdatedAt)
		assert.NotEqual(t, category.UpdatedAt, categoryUpdatedAtFirst)
	})

	t.Run("should not update a category", func(t *testing.T) {
		category := &entity.Category{
			ID:    9999,
			Title: "Gopher Blue",
			Color: "#00ADD8",
		}
		err := repo.Update(context.Background(), category)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrCategoryNotFound)
	})
}

func TestDeleteCategory(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should delete a category", func(t *testing.T) {
		category := &entity.Category{
			Title: "Gopher Blue",
			Color: "#00ADD8",
		}
		err := repo.Insert(context.Background(), category)
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), category.ID)
		assert.NoError(t, err)

		category, err = repo.FindByID(context.Background(), category.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrCategoryNotFound)
		assert.Nil(t, category)
	})

	t.Run("should not delete a category", func(t *testing.T) {
		err := repo.Delete(context.Background(), 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrCategoryNotFound)
	})
}
