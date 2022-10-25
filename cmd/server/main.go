package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	videoRepo := video.NewRepository(db.Client)
	videoServices := video.NewService(videoRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	httpServer := http.Server{Addr: ":8080", Handler: r}

	httpTransport.HandleVideoRequests(r, videoServices)

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
