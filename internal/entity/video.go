package entity

import "time"

type Video struct {
	ID          uint
	Title       string
	Description string
	Url         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
