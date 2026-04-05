package main

import (
	"log"

	httptransport "github.com/argform/baitfolio-backend/internal/transport/http"
)

func main() {
	r := httptransport.NewRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}