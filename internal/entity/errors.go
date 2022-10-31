package entity

import "errors"

var (
	ErrVideoNotFound    = errors.New("video not found")
	ErrCategoryNotFound = errors.New("category not found")
)
