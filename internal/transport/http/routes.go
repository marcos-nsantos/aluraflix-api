package httpTransport

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marcos-nsantos/aluraflix-api/internal/category"
	"github.com/marcos-nsantos/aluraflix-api/internal/video"
)

func HandleRequests(r *chi.Mux, vs video.Service, cs category.Service) {
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	handleVideoRequests(r, vs)
	handleCategoryRequests(r, cs)
}

func handleVideoRequests(r *chi.Mux, service video.Service) {
	r.Group(func(r chi.Router) {
		r.Post("/videos", postVideo(service))
		r.Get("/videos", getAllVideos(service))
		r.Get("/videos/{id}", getVideoByID(service))
		r.Put("/videos/{id}", updateVideo(service))
		r.Delete("/videos/{id}", deleteVideo(service))
	})
}

func handleCategoryRequests(r *chi.Mux, service category.Service) {
	r.Group(func(r chi.Router) {
		r.Post("/categories", postCategory(service))
		r.Get("/categories", getAllCategories(service))
	})
}

func Handle(h *http.Server) error {
	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
