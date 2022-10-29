package httpTransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
	"github.com/marcos-nsantos/aluraflix-api/internal/transport/http/presenters"
	"github.com/marcos-nsantos/aluraflix-api/internal/validator"
	"github.com/marcos-nsantos/aluraflix-api/internal/video"
)

type PostVideoRequest struct {
	Title       string `json:"title" validate:"required,notblank"`
	Description string `json:"description" validate:"required,notblank"`
	URL         string `json:"url" validate:"required,url"`
}

func convertPostVideoRequestToVideo(video *PostVideoRequest) entity.Video {
	return entity.Video{
		Title:       video.Title,
		Description: video.Description,
		URL:         video.URL,
	}
}

func postVideo(service video.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postVideo PostVideoRequest
		if err := json.NewDecoder(r.Body).Decode(&postVideo); err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid request body"))
			return
		}

		if ok, err := validator.Validate(postVideo); !ok {
			presenters.JSONValidationResponse(w, err)
			return
		}

		videoConverted := convertPostVideoRequestToVideo(&postVideo)
		if err := service.Post(r.Context(), &videoConverted); err != nil {
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		videoResponse := presenters.VideoResponse(&videoConverted)
		presenters.JSONResponse(w, http.StatusCreated, videoResponse)
	}
}

func getAllVideos(service video.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videos, err := service.GetAll(r.Context())
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		videosResponse := presenters.VideosResponse(videos)
		presenters.JSONResponse(w, http.StatusOK, videosResponse)
	}
}

func getVideoByID(service video.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		videoFromDB, err := service.GetByID(r.Context(), idUint)
		if err != nil {
			if errors.Is(err, entity.ErrVideoNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		videoResponse := presenters.VideoResponse(videoFromDB)
		presenters.JSONResponse(w, http.StatusOK, videoResponse)
	}
}

func updateVideo(service video.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		var postVideo PostVideoRequest
		if err = json.NewDecoder(r.Body).Decode(&postVideo); err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid request body"))
			return
		}

		if ok, errs := validator.Validate(postVideo); !ok {
			presenters.JSONValidationResponse(w, errs)
			return
		}

		videoConverted := convertPostVideoRequestToVideo(&postVideo)
		if err = service.Update(r.Context(), &videoConverted, idUint); err != nil {
			if errors.Is(err, entity.ErrVideoNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		videoResponse := presenters.VideoResponse(&videoConverted)
		presenters.JSONResponse(w, http.StatusOK, videoResponse)
	}
}

func deleteVideo(service video.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		if err = service.Delete(r.Context(), idUint); err != nil {
			if errors.Is(err, entity.ErrVideoNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		presenters.JSONResponse(w, http.StatusOK, "video deleted")
	}
}
