.PHONY: up down

# Variables
APP_NAME_PRODUCER := producer
APP_NAME_CONSUMER := consumer
GO_FILES_PRODUCER := main_producer.go producer.go kafka.go
GO_FILES_CONSUMER := main_consumer.go consumer.go kafka.go
BUILD_DIR := ./build

# Default target
all: build

# Build the producer binary
build-producer:
	@echo "Building producer..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME_PRODUCER) $(GO_FILES_PRODUCER)

# Build the consumer binary
build-consumer:
	@echo "Building consumer..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME_CONSUMER) $(GO_FILES_CONSUMER)

# Build both binaries
build: build-producer build-consumer

# Run the producer
run-producer: build-producer
	@echo "Running producer..."
	@$(BUILD_DIR)/$(APP_NAME_PRODUCER)

# Run the consumer
run-consumer: build-consumer
	@echo "Running consumer..."
	@$(BUILD_DIR)/$(APP_NAME_CONSUMER)

# Clean up the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Build artifacts cleaned."

# Stop and remove Docker containers and volumes, then delete kafka-data folder
down:
	@echo "Stopping and removing Docker containers..."
	@docker-compose down -v
	@echo "Removing kafka-data folder..."
	@rm -rf ./kafka-data
	@echo "Containers and kafka-data folder removed."

# Start docker-compose in detached mode
up: down
	@echo "Starting Docker Compose..."
	@docker-compose up -d

# Help target
help:
	@echo "Usage:"
	@echo "  make build          - Build both producer and consumer"
	@echo "  make build-producer - Build the producer binary"
	@echo "  make build-consumer - Build the consumer binary"
	@echo "  make run-producer   - Run the producer application"
	@echo "  make run-consumer   - Run the consumer application"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make up             - Stop and remove containers, then start Docker Compose"
	@echo "  make down           - Stop and remove containers and kafka-data folder"
