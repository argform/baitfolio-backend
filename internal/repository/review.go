package repository

import (
	"context"

	"github.com/argform/baitfolio-backend/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) (*domain.Review, error)
	GetByID(ctx context.Context, id uint64) (*domain.Review, error)
	GetAllByPointID(ctx context.Context, pointID uint64) ([]*domain.Review, error)
	DeleteByID(ctx context.Context, id uint64) error
}
