package service

import (
	"context"
	"errors"
	"testing"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type reviewRepositoryStub struct {
	review       *domain.Review
	getErr       error
	deleteErr    error
	deletedID    uint64
	deleteCalled bool
}

func (s *reviewRepositoryStub) Create(context.Context, *domain.Review) (*domain.Review, error) {
	return nil, errors.New("not implemented")
}

func (s *reviewRepositoryStub) GetByID(context.Context, uint64) (*domain.Review, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}

	return s.review, nil
}

func (s *reviewRepositoryStub) GetAllByPointID(context.Context, uint64) ([]*domain.Review, error) {
	return nil, errors.New("not implemented")
}

func (s *reviewRepositoryStub) DeleteByID(_ context.Context, id uint64) error {
	s.deleteCalled = true
	s.deletedID = id
	return s.deleteErr
}

func TestReviewService_Delete_OwnerCanDelete(t *testing.T) {
	authorID := uint64(7)
	repo := &reviewRepositoryStub{
		review: &domain.Review{
			ReviewID: 11,
			AuthorID: &authorID,
		},
	}
	service := NewReviewService(repo)

	err := service.Delete(context.Background(), 11, authorID)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if !repo.deleteCalled {
		t.Fatal("expected DeleteByID to be called")
	}

	if repo.deletedID != 11 {
		t.Fatalf("expected deleted id 11, got %d", repo.deletedID)
	}
}

func TestReviewService_Delete_OtherUserForbidden(t *testing.T) {
	authorID := uint64(7)
	repo := &reviewRepositoryStub{
		review: &domain.Review{
			ReviewID: 11,
			AuthorID: &authorID,
		},
	}
	service := NewReviewService(repo)

	err := service.Delete(context.Background(), 11, 8)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	if repo.deleteCalled {
		t.Fatal("did not expect DeleteByID to be called")
	}
}

func TestReviewService_Delete_ReviewNotFound(t *testing.T) {
	repo := &reviewRepositoryStub{
		getErr: repository.ErrReviewNotFound,
	}
	service := NewReviewService(repo)

	err := service.Delete(context.Background(), 11, 7)
	if !errors.Is(err, repository.ErrReviewNotFound) {
		t.Fatalf("expected ErrReviewNotFound, got %v", err)
	}
}
