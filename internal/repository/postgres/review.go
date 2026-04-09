package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type PostgresReviewRepository struct {
	db *pgxpool.Pool
}

func NewPostgresReviewRepository(db *pgxpool.Pool) *PostgresReviewRepository {
	return &PostgresReviewRepository{db: db}
}

func scanReview(row pgx.Row) (*domain.Review, error) {
	var review domain.Review

	err := row.Scan(
		&review.ReviewID,
		&review.AuthorID,
		&review.PointID,
		&review.Score,
		&review.Content,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *PostgresReviewRepository) Create(ctx context.Context, review *domain.Review) (*domain.Review, error) {
	query := `
		INSERT INTO reviews (
			author_id,
			point_id,
			score,
			content
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			review_id,
			author_id,
			point_id,
			score,
			content,
			created_at,
			updated_at
	`

	created, err := scanReview(
		r.db.QueryRow(
			ctx,
			query,
			review.AuthorID,
			review.PointID,
			review.Score,
			review.Content,
		),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("review already exists")
		}
		return nil, fmt.Errorf("create review: %w", err)
	}

	return created, nil
}

func (r *PostgresReviewRepository) GetByID(ctx context.Context, id uint64) (*domain.Review, error) {
	row := r.db.QueryRow(ctx, `SELECT * FROM reviews WHERE review_id = $1`, id)
	review, err := scanReview(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrReviewNotFound
		}
		return nil, err
	}

	return review, nil
}

func (r *PostgresReviewRepository) GetAllByPointID(ctx context.Context, pointID uint64) ([]*domain.Review, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT * FROM reviews
		WHERE point_id = $1
		ORDER BY created_at DESC, review_id DESC
		`,
		pointID,
	)
	if err != nil {
		return nil, fmt.Errorf("get reviews by point id: %w", err)
	}
	defer rows.Close()

	reviews := make([]*domain.Review, 0)
	for rows.Next() {
		review, scanErr := scanReview(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan reviews by point id: %w", scanErr)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews by point id: %w", err)
	}

	return reviews, nil
}

func (r *PostgresReviewRepository) DeleteByID(ctx context.Context, id uint64) error {
	commandTag, err := r.db.Exec(ctx, `DELETE FROM reviews WHERE review_id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete review: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return repository.ErrReviewNotFound
	}

	return nil
}
