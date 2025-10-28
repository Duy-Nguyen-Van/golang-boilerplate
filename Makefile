.PHONY: lint mocks tests

bootstrap: container-up migration-up up

lint:
	golangci-lint run

mocks:
	mockery --case snake --dir ./repositories --all --output ./mocks/repositories
	mockery --case snake --dir ./adapters --all --output ./mocks/adapters

tests:
	go test -v -cover -race -timeout 300s -count=1 ./...

build:
	@cd cmd/${cmd} && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${service_name} .

dep:
	@go mod tidy

container-up:
	docker compose up -d

container-down:
	docker compose down

up:
	cd cmd/server && go run main.go

major-version-update:
	go get -u -t ./...

minor-version-update:
	go get -u ./...

migration-status:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/sql status

migration-up:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/sql up

migration-create:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/sql create ${name} sql

migration-down:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/sql down

seed-up:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/seed up

seed-down:
	set -a && source cmd/migrations/.env && set +a && goose --dir ./cmd/migrations/seed down

format:
	go fmt ./...

swagger-load:
	swag init \
		-g main.go \
		-d ./cmd/server,./internal/handlers,./internal/middlewares,./internal/services,./internal/repositories,./internal/models,./internal/utils,./internal/config,./internal/constants,./internal/integration,./internal/dtos,./internal/logger,./internal/db,./internal/cache \
		--output ./docs