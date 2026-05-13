.PHONY: run build migrate-up migrate-down docker-up docker-down

DATABASE_URL ?= postgres://mezgeb:mezgeb@localhost:5432/mezgeb?sslmode=disable

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
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
