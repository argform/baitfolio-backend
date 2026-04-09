package service

import (
	"context"
	"fmt"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/geo"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type PointService struct {
	points repository.PointRepository
}

type CreatePointInput struct {
	PointID              uint64
	CreatedBy            *uint64
	Name                 string
	Description          *string
	Lat                  float64
	Lon                  float64
	WaterbodyHydrologyID *int32
	ShoreTypeID          *int16
	AccessTypeID         *int16
}

func NewPointService(points repository.PointRepository) *PointService {
	return &PointService{
		points: points,
	}
}

func (s *PointService) Create(ctx context.Context, input CreatePointInput) (*domain.Point, error) {
	point := &domain.Point{
		CreatedBy:            input.CreatedBy,
		Name:                 input.Name,
		Description:          input.Description,
		Lat:                  input.Lat,
		Lon:                  input.Lon,
		WaterbodyHydrologyID: input.WaterbodyHydrologyID,
		ShoreTypeID:          input.ShoreTypeID,
		AccessTypeID:         input.AccessTypeID,
	}

	createdPoint, err := s.points.Create(ctx, point)
	if err != nil {
		return nil, fmt.Errorf("create point: %w", err)
	}

	return createdPoint, nil
}

func (s *PointService) GetByID(ctx context.Context, id uint64) (*domain.Point, error) {
	point, err := s.points.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get point by id: %w", err)
	}

	return point, nil
}

func (s *PointService) GetAllInsideTile(ctx context.Context, tile geo.Tile) ([]*domain.Point, error) {
	points, err := s.points.GetAllInsideTile(ctx, tile)
	if err != nil {
		return nil, fmt.Errorf("get points inside tile: %w", err)
	}

	return points, nil
}
