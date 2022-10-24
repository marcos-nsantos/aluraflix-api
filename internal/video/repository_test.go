//go:build integration

package video

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

func TestCreateVideo(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should create a new video", func(t *testing.T) {
		video := &entity.Video{
			Title:       "O que é e pra que serve a linguagem Go?",
			Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
			URL:         "https://youtu.be/KfCNyIrqjsg",
		}

		err := repo.Create(context.Background(), video)
		assert.NoError(t, err)
		assert.NotEmpty(t, video.ID)
		assert.NotEmpty(t, video.CreatedAt)
		assert.NotEmpty(t, video.UpdatedAt)
	})
}
