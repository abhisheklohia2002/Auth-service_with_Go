DATABASE_URL=postgresql://gorm_demo_user:rTDC3Or8XAb6Nr8731YDBhN40cQGJvRH@dpg-d8cgmbh9rddc73d9o8f0-a/gorm_demo
MIGRATIONS_DIR=migrations

run:
	go run ./cmd/api

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql