# Load .env for migration targets (DB_* vars)
-include .env
export

# Database URL for golang-migrate (postgres)
DB_URL ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: server build run clean migrate-up migrate-down migrate-sql-up migrate-sql-down migrate-sql-create migrate-install

server:
	nodemon --exec go run ./cmd/server/main.go --signal SIGTERM

build:
	go build -o bin/server cmd/server/main.go

run:
	./bin/server

clean:
	rm -f bin/server

# --- Schema from Go models (GORM AutoMigrate) ---
# Add a column to your model, then run: make migrate-up
# For new NOT NULL columns on existing tables, add default in gorm tag: not null;default:''
migrate-up:
	go run ./cmd/migrate/main.go

# --- Optional: file-based SQL migrations (golang-migrate) ---
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-sql-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-sql-down-all:
	migrate -path migrations -database "$(DB_URL)" down

migrate-sql-create:
	@test -n "$(name)" || (echo "Usage: make migrate-sql-create name=your_migration_name"; exit 1)
	migrate create -ext sql -dir migrations -seq $(name)