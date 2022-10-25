package httpTransport

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/marcos-nsantos/aluraflix-api/internal/video"
)

func HandleVideoRequests(r *chi.Mux, service video.Service) {
	r.Group(func(r chi.Router) {
		r.Post("/videos", postVideo(service))
		r.Get("/videos", getAllVideos(service))
		r.Get("/videos/{id}", getVideoByID(service))
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
