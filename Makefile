.PHONY: dev build run test lint tidy docker-up docker-down migrate

APP_NAME=expense-tracker
BUILD_DIR=./bin

## dev: Start the API with hot-reload (requires Air)
dev:
	air -c .air.toml

## build: Compile the binary
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/server ./cmd/server

## run: Build and run without hot-reload
run: build
	$(BUILD_DIR)/server

## test: Run all tests with race detection
test:
	go test -race -cover ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## tidy: Tidy and verify Go modules
tidy:
	go mod tidy
	go mod verify

## docker-up: Start all services
docker-up:
	docker compose up --build

## docker-down: Stop all services
docker-down:
	docker compose down

## docker-logs: Tail API logs
docker-logs:
	docker compose logs -f api

## docker-build: Build the production image
docker-build:
	docker build --target production -t $(APP_NAME):latest .

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //'
