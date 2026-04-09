package domain

import (
	"time"
)

type Point struct {
	PointID              uint64
	CreatedBy            *uint64
	Name                 string
	Description          *string
	Lat                  float64
	Lon                  float64
	WaterbodyHydrologyID *int32
	ShoreTypeID          *int16
	AccessTypeID         *int16
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
