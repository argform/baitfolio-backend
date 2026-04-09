package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"

	"github.com/argform/baitfolio-backend/internal/auth"
	"github.com/argform/baitfolio-backend/internal/config"
	"github.com/argform/baitfolio-backend/internal/db"
	"github.com/argform/baitfolio-backend/internal/repository/postgres"
	"github.com/argform/baitfolio-backend/internal/service"
	httptransport "github.com/argform/baitfolio-backend/internal/transport/http"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := db.NewPostgresPool(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := postgres.NewPostgresUserRepository(pool)
	pointRepo := postgres.NewPostgresPointRepository(pool)
	reviewRepo := postgres.NewPostgresReviewRepository(pool)
	jwtManager := auth.NewJWTManager("super-secret-key", 7*24*time.Hour)
	authService := service.NewAuthService(userRepo, jwtManager)
	pointService := service.NewPointService(pointRepo)
	reviewService := service.NewReviewService(reviewRepo)

	router := httptransport.NewRouter(httptransport.Dependencies{
		AuthService:   authService,
		PointService:  pointService,
		ReviewService: reviewService,
		JWTManager:    jwtManager,
	})

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
