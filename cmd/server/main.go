package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcos-nsantos/aluraflix-api/internal/category"
	"github.com/marcos-nsantos/aluraflix-api/internal/database"
	httpTransport "github.com/marcos-nsantos/aluraflix-api/internal/transport/http"
	"github.com/marcos-nsantos/aluraflix-api/internal/validator"
	"github.com/marcos-nsantos/aluraflix-api/internal/video"
)

func run() error {
	fmt.Println("Starting up the application")

	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err = db.Migrate("file:///migrations"); err != nil {
		return err
	}

	if err = validator.RegisterValidators(); err != nil {
		return err
	}

	r := chi.NewRouter()
	httpServer := http.Server{Addr: ":8080", Handler: r}

	videoRepo := video.NewRepository(db.Client)
	videoServices := video.NewService(videoRepo)
	categoryRepo := category.NewRepository(db.Client)
	categoryServices := category.NewService(categoryRepo)

	httpTransport.HandleRequests(r, videoServices, categoryServices)
	if err = httpTransport.Handle(&httpServer); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
