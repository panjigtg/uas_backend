include .env

DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH=./database/migrations

dev:
	air

run:
	go run ./cmd/main.go

test:
	go test ./...

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-drop:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)