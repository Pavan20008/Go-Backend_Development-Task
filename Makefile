.PHONY: help run build test vet tidy sqlc migrate-up migrate-down docker-up docker-down

DB_URL ?= postgres://postgres:postgres@localhost:5432/user_age_api?sslmode=disable

help:
	@echo "Targets:"
	@echo "  run          Run the API server locally"
	@echo "  build        Build the server binary into ./bin"
	@echo "  test         Run unit tests"
	@echo "  vet          Run go vet"
	@echo "  tidy         Tidy go.mod/go.sum"
	@echo "  sqlc         Regenerate the sqlc query layer"
	@echo "  migrate-up   Apply migrations using golang-migrate (requires migrate CLI)"
	@echo "  migrate-down Roll back the last migration"
	@echo "  docker-up    Start API + Postgres via docker compose"
	@echo "  docker-down  Stop docker compose stack"

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

docker-up:
	docker compose up --build

docker-down:
	docker compose down
