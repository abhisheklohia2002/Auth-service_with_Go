MIGRATIONS_DIR=migrations
GOOSE_DRIVER=postgres
APP_BIN=bin/api

run:
	go run ./cmd/api

build:
	go build -o $(APP_BIN) ./cmd/api

migrate-up:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DATABASE_URL)" status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql