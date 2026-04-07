package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/geo"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type PostgresPointRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPointRepository(db *pgxpool.Pool) *PostgresPointRepository {
	return &PostgresPointRepository{db: db}
}

func scanPoint(row pgx.Row) (*domain.Point, error) {
	var point domain.Point

	err := row.Scan(
		&point.PointID,
		&point.CreatedBy,
		&point.Name,
		&point.Description,
		&point.Lat,
		&point.Lon,
		&point.CreatedAt,
		&point.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &point, err
}

func (r *PostgresPointRepository) Create(ctx context.Context, point *domain.Point) (*domain.Point, error) {
	query := `
		INSERT INTO points (
			created_by,
			name,
			description,
			lat,
			lon
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			point_id,
			created_by,
			name,
			description,
			lat,
			lon,
			created_at,
			updated_at
	`

	created, err := scanPoint(
		r.db.QueryRow(
			ctx,
			query,
			point.PointID,
			point.CreatedBy,
			point.Name,
			point.Description,
			point.Lat,
			point.Lon,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create point: %w", err)
	}

	return created, nil
}

func (r *PostgresPointRepository) GetByID(ctx context.Context, id uint64) (*domain.Point, error) {
	row := r.db.QueryRow(ctx, `SELECT * FROM points WHERE id = $1`, id)
	point, err := scanPoint(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrPointNotFound
		}
		return nil, err
	}
	return point, nil
}

func (r *PostgresPointRepository) GetAllInsideTile(ctx context.Context, t geo.Tile) ([]*domain.Point, error) {
	bbox := geo.TileToBBox(t)
	query := `
		SELECT * FROM points
		WHERE (lat BETWEEN $1 AND $2) 
			AND (lon BETWEEN $3 AND $4)
	`
	rows, err := r.db.Query(
		ctx,
		query,
		bbox.South,
		bbox.North,
		bbox.West,
		bbox.East,
	)
	if err != nil {
		return nil, fmt.Errorf("get points inside tile: %w", err)
	}
	
	defer rows.Close()

	points := make([]*domain.Point, 0)

	for rows.Next() {
		point, err := scanPoint(rows)
		if err != nil {
			return nil, fmt.Errorf("scan point inside tile: %w", err)
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate points inside tile: %w", err)
	}

	return points, nil
}