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

func TestFindAllVideos(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should return all videos", func(t *testing.T) {
		videos := []*entity.Video{
			{
				Title:       "O que é e pra que serve a linguagem Go?",
				Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
				URL:         "https://youtu.be/KfCNyIrqjsg",
			},
			{
				Title: "Grandes sistemas em PHP com Vinicius Dias",
				Description: "Os Grandes Sistemas em que o PHP foi usado, desde o Wordpress ou Magento, até a sua " +
					"evolução com PHP 7 e PHP 8 aumentando a performance do código em comparação com a Hack " +
					"Language, esta última criada pela Meta (Facebook) que começou utilizando PHP.",
				URL: "https://youtu.be/arZCoJMSTlI",
			},
			{
				Title: "Os MELHORES livros de tecnologia para ler em Programação com Roberta Arcoverde",
				Description: "Conheça os melhores livros de tecnologia para se aprender computação ou programação, " +
					"sejam iniciantes ou avançados, para base acadêmica, prática de programar e/ou de carreira " +
					"no mundo do desenvolvimento.",
				URL: "https://youtu.be/RvWQQRjz1Pw",
			},
		}

		for _, video := range videos {
			err := repo.Create(context.Background(), video)
			assert.NoError(t, err)
		}

		videos, err := repo.FindAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, videos)
		assert.GreaterOrEqual(t, len(videos), 3)
	})
}

func TestFindVideoByID(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should return a video by ID", func(t *testing.T) {
		video := &entity.Video{
			Title:       "O que é e pra que serve a linguagem Go?",
			Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
			URL:         "https://youtu.be/KfCNyIrqjsg",
		}

		err := repo.Create(context.Background(), video)
		assert.NoError(t, err)

		video, err = repo.FindByID(context.Background(), video.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, video)
		assert.NotEmpty(t, video.ID)
		assert.NotEmpty(t, video.CreatedAt)
		assert.NotEmpty(t, video.UpdatedAt)
		assert.Equal(t, "O que é e pra que serve a linguagem Go?", video.Title)
		assert.Equal(t, "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?", video.Description)
		assert.Equal(t, "https://youtu.be/KfCNyIrqjsg", video.URL)
	})
}

func TestUpdateVideo(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should update a video", func(t *testing.T) {
		video := &entity.Video{
			Title:       "O que é e pra que serve a linguagem Go?",
			Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
			URL:         "https://youtu.be/KfCNyIrqjsg",
		}

		err := repo.Create(context.Background(), video)
		assert.NoError(t, err)

		video.Title = "O que é e pra que serve a linguagem Python?"
		video.Description = "Você provavelmente já ouviu falar da linguagem de programação Python. Mas qual o propósito dela?"
		video.URL = "https://youtu.be/KfCNyIrqjsg"

		err = repo.Update(context.Background(), video)
		assert.NoError(t, err)

		video, err = repo.FindByID(context.Background(), video.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, video)
		assert.NotEmpty(t, video.ID)
		assert.NotEmpty(t, video.CreatedAt)
		assert.NotEmpty(t, video.UpdatedAt)
		assert.Equal(t, "O que é e pra que serve a linguagem Python?", video.Title)
		assert.Equal(t, "Você provavelmente já ouviu falar da linguagem de programação Python. Mas qual o propósito dela?", video.Description)
	})

	t.Run("should return error when video does not exist", func(t *testing.T) {
		video := &entity.Video{
			ID:          9999,
			Title:       "O que é e pra que serve a linguagem Go?",
			Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
			URL:         "https://youtu.be/KfCNyIrqjsg",
		}

		err := repo.Update(context.Background(), video)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrVideoNotFound)
	})
}

func TestDeleteVideo(t *testing.T) {
	repo := NewRepository(dbConn)

	t.Run("should delete a video", func(t *testing.T) {
		video := &entity.Video{
			Title:       "O que é e pra que serve a linguagem Go?",
			Description: "Você provavelmente já ouviu falar da linguagem de programação Go. Mas qual o propósito dela?",
			URL:         "https://youtu.be/KfCNyIrqjsg",
		}

		err := repo.Create(context.Background(), video)
		assert.NoError(t, err)

		err = repo.Delete(context.Background(), video.ID)
		assert.NoError(t, err)

		video, err = repo.FindByID(context.Background(), video.ID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrVideoNotFound)
		assert.Empty(t, video)
	})

	t.Run("should return error when video does not exist", func(t *testing.T) {
		err := repo.Delete(context.Background(), 9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrVideoNotFound)
	})
}
