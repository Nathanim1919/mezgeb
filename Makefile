.PHONY: run build migrate-up migrate-down docker-up docker-down

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/mezgeb?sslmode=disable
MIGRATE ?= $(shell go env GOPATH)/bin/migrate

run:
	go run ./cmd/bot

build:
	go build -o mezgeb ./cmd/bot

docker-up:
	docker compose up -d db
	@echo "Waiting for Postgres..."
	@sleep 2
	@echo "Ready. Run 'make migrate-up' then 'make run'"

docker-down:
	docker compose down

migrate-up:
	$(MIGRATE) -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	$(MIGRATE) -database "$(DATABASE_URL)" -path migrations down 1

migrate-create:
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)
