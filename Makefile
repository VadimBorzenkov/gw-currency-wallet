# Variables
APP_NAME = grpc-exchange
DOCKER_IMAGE = grpc-exchange:latest
BUILD_DIR = build
GO_FILES = $(shell find . -name '*.go' -not -path "./vendor/*")
PROTO_FILES = $(shell find ./proto -name '*.proto')

# Commands
.PHONY: all build run test docker-build docker-run clean proto

# Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

# Run the application locally
run:
	@echo "Running $(APP_NAME)..."
	go run ./cmd/main.go

# Test the application
test:
	@echo "Running tests..."
	go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)


# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run the Docker container
docker-run:
	@echo "Running Docker container..."
	docker-compose up

# Stop and remove Docker containers
docker-clean:
	@echo "Stopping and removing Docker containers..."
	docker-compose down
