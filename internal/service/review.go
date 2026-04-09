package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

type ReviewService struct {
	reviews repository.ReviewRepository
}

type CreateReviewInput struct {
	AuthorID *uint64
	PointID  uint64
	Score    int16
	Content  *string
}

func NewReviewService(reviews repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		reviews: reviews,
	}
}

func (s *ReviewService) Create(ctx context.Context, input CreateReviewInput) (*domain.Review, error) {
	review := &domain.Review{
		AuthorID: input.AuthorID,
		PointID:  input.PointID,
		Score:    input.Score,
		Content:  input.Content,
	}

	createdReview, err := s.reviews.Create(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}

	return createdReview, nil
}

func (s *ReviewService) GetAllByPointID(ctx context.Context, pointID uint64) ([]*domain.Review, error) {
	reviews, err := s.reviews.GetAllByPointID(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("get reviews by point id: %w", err)
	}

	return reviews, nil
}

func (s *ReviewService) Delete(ctx context.Context, reviewID, authorID uint64) error {
	review, err := s.reviews.GetByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("get review by id: %w", err)
	}

	if review.AuthorID == nil || *review.AuthorID != authorID {
		return ErrForbidden
	}

	if err := s.reviews.DeleteByID(ctx, reviewID); err != nil {
		return fmt.Errorf("delete review: %w", err)
	}

	return nil
}
