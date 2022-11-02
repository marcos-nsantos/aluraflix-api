package entity

import "time"

type Video struct {
	ID          uint64
	Title       string
	Description string
	URL         string
	CategoryID  uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
