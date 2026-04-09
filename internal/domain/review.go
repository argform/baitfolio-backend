package domain

import "time"

type Review struct {
	ReviewID  uint64
	AuthorID  *uint64
	PointID   uint64
	Score     int16
	Content   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
