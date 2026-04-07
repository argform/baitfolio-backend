package repository

import (
	"context"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/geo"
)

type PointRepository interface {
	Create(ctx context.Context, point *domain.Point) (*domain.Point, error)
	GetByID(ctx context.Context, id uint64) (*domain.Point, error)
	GetAllInsideTile(ctx context.Context, t geo.Tile) ([]*domain.Point, error)
}