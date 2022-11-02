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
	handleCategoryRequests(r, cs, vs)
}

func handleVideoRequests(r *chi.Mux, vs video.Service) {
	r.Group(func(r chi.Router) {
		r.Post("/videos", postVideo(vs))
		r.Get("/videos", getAllVideos(vs))
		r.Get("/videos/{id}", getVideoByID(vs))
		r.Put("/videos/{id}", updateVideo(vs))
		r.Delete("/videos/{id}", deleteVideo(vs))
	})
}

func handleCategoryRequests(r *chi.Mux, cs category.Service, vs video.Service) {
	r.Group(func(r chi.Router) {
		r.Post("/categories", postCategory(cs))
		r.Get("/categories", getAllCategories(cs))
		r.Get("/categories/{id}", getCategoryByID(cs))
		r.Put("/categories/{id}", updateCategory(cs))
		r.Delete("/categories/{id}", deleteCategory(cs))
		r.Get("/categories/{id}/videos", getAllVideosByCategory(vs))
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
