APP_NAME = trading_ace
DB_NAME = pelith
DB_PORT = 5433
DB_HOST = localhost
DOCKER_COMPOSE = docker-compose
GO_FILES = $(shell find . -type f -name '*.go')

.PHONY: all help modvendor build run test lint docker-build docker-up docker-down docker-logs docker-clean docker-db

all: help

help:
	@echo "Available make commands:"
	@echo "  make modvendor      - Run 'go mod vendor'"
	@echo "  make build          - Build the Go application"
	@echo "  make run            - Run the Go application"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start services with Docker Compose"
	@echo "  make docker-down    - Stop services with Docker Compose"
	@echo "  make docker-logs    - Show logs of running Docker services"
	@echo "  make docker-clean   - Clean up generated files and Docker containers"
	@echo "  make docker-db      - Connect to Docker PostgreSQL"

modvendor:
	go mod tidy
	go mod vendor

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) ./main.go

run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./... -v

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

docker-build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "Starting Docker services..."
	$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) down

docker-logs:
	@echo "Showing Docker logs..."
	$(DOCKER_COMPOSE) logs -f app

docker-clean:
	@echo "Cleaning up..."
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	rm -f $(APP_NAME)

docker-db:
	@echo "Connecting to Docker PostgreSQL..."
	psql -h $(DB_HOST) -p $(DB_PORT) -U admin -d $(DB_NAME)