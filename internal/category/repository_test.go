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

	err = db.UndoMigrations("file://../..//migrations")
	if err != nil {
		fmt.Println(err)
	}

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
