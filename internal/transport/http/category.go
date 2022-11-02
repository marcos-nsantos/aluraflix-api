package httpTransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marcos-nsantos/aluraflix-api/internal/category"
	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
	"github.com/marcos-nsantos/aluraflix-api/internal/transport/http/presenters"
	"github.com/marcos-nsantos/aluraflix-api/internal/validator"
)

type PostCategoryRequest struct {
	Title string `json:"title" validate:"required,notblank"`
	Color string `json:"color" validate:"required,notblank,hexcolor"`
}

func convertPostCategoryRequestToCategory(category *PostCategoryRequest) entity.Category {
	return entity.Category{
		Title: category.Title,
		Color: category.Color,
	}
}

func postCategory(service category.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postCategory PostCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&postCategory); err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid request body"))
			return
		}

		if ok, err := validator.Validate(postCategory); !ok {
			presenters.JSONValidationResponse(w, err)
			return
		}

		categoryConverted := convertPostCategoryRequestToCategory(&postCategory)
		if err := service.Post(r.Context(), &categoryConverted); err != nil {
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		categoryResponse := presenters.CategoryResponse(&categoryConverted)
		presenters.JSONResponse(w, http.StatusCreated, categoryResponse)
	}
}

func getAllCategories(service category.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := service.GetAll(r.Context())
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		categoriesResponse := presenters.CategoriesResponse(categories)
		presenters.JSONResponse(w, http.StatusOK, categoriesResponse)
	}
}

func getCategoryByID(service category.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		categoryFromDB, err := service.GetByID(r.Context(), idUint)
		if err != nil {
			if errors.Is(err, entity.ErrCategoryNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		categoryResponse := presenters.CategoryResponse(categoryFromDB)
		presenters.JSONResponse(w, http.StatusOK, categoryResponse)
	}
}

func updateCategory(service category.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		var postCategory PostCategoryRequest
		if err = json.NewDecoder(r.Body).Decode(&postCategory); err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid request body"))
			return
		}

		if ok, err := validator.Validate(postCategory); !ok {
			presenters.JSONValidationResponse(w, err)
			return
		}

		categoryConverted := convertPostCategoryRequestToCategory(&postCategory)
		categoryConverted.ID = idUint
		if err = service.Update(r.Context(), &categoryConverted); err != nil {
			if errors.Is(err, entity.ErrCategoryNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		categoryResponse := presenters.CategoryResponse(&categoryConverted)
		presenters.JSONResponse(w, http.StatusOK, categoryResponse)
	}
}

func deleteCategory(service category.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			presenters.JSONErrorResponse(w, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		if err = service.Delete(r.Context(), idUint); err != nil {
			if errors.Is(err, entity.ErrCategoryNotFound) {
				presenters.JSONErrorResponse(w, http.StatusNotFound, err)
				return
			}
			presenters.JSONErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		presenters.JSONResponse(w, http.StatusNoContent, nil)
	}
}
