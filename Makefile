include .env

build:
	go build ./cmd/main.go

test:
	go test ./... -v -cover

docker-build:
	docker build --tag exchange:dev  .

run:
	go run ./cmd/main.go

lint:
	golangci-lint run ./...

generate:
	go generate ./...

migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} status -v

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${DB_CONN} down -v