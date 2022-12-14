package presenters

import "github.com/marcos-nsantos/aluraflix-api/internal/entity"

type Category struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func CategoryResponse(category *entity.Category) *Category {
	return &Category{
		ID:        category.ID,
		Title:     category.Title,
		Color:     category.Color,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func CategoriesResponse(categories []*entity.Category) []*Category {
	var categoriesResponse []*Category
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, CategoryResponse(category))
	}
	return categoriesResponse
}
