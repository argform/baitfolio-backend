package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User

	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.About,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (
			username,
			email,
			password_hash,
			first_name,
			last_name,
			about
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			user_id,
			username,
			email,
			password_hash,
			first_name,
			last_name,
			about,
			created_at,
			updated_at;
	`

	created, err := scanUser(
		r.db.QueryRow(
			ctx,
			query,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.FirstName,
			user.LastName,
			user.About,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT * FROM users WHERE email = $1`, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT * FROM users WHERE user_id = $1`, id)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
