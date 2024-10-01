include .env

build:
	docker compose build

test:
	go test ./... -v -cover

run:
	docker compose up -d

lint:
	golangci-lint run ./...

migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} status -v

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} down -v