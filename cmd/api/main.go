package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/argform/baitfolio-backend/internal/config"
	"github.com/argform/baitfolio-backend/internal/db"

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
	
	r := httptransport.NewRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}