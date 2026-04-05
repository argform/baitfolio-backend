APP_NAME=baitfolio-backend
DB_URL=postgres://baitfolio:baitfolio@localhost:5432/baitfolio?sslmode=disable

.PHONY: run build test fmt tidy up down logs ps \
        migrate-up migrate-down migrate-version migrate-force

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) ./cmd/api/main.go

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

migrate-force:
	@echo "Usage: make migrate-force VERSION=<version>"
	migrate -path migrations -database "$(DB_URL)" force $(VERSION)