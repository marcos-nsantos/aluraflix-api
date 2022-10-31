package entity

import "time"

type Category struct {
	ID        uint64
	Title     string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
