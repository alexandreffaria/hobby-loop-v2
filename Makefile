# Makefile for Hobby Loop API

# .PHONY tells Make that these are commands, not file names
.PHONY: all build run test clean docker-up docker-down

# default target: if you just type 'make', it builds the app
all: build

# ----------------------------
# 1. Development & Running
# ----------------------------

# Run the server directly
run:
	@echo "🚀 Starting server..."
	go run cmd/api/main.go

# Build the binary into a 'bin' folder
build:
	@echo "🔨 Building binary..."
	go build -o bin/api cmd/api/main.go

# ----------------------------
# 2. Testing & Quality
# ----------------------------

# Run all tests with verbose output and no caching
test:
	@echo "🧪 Running all tests..."
	go test ./... -v -count=1

# ----------------------------
# 3. Infrastructure (Docker)
# ----------------------------

# Start the database in the background
docker-up:
	@echo "🐳 Starting database..."
	docker compose up -d

# Stop and remove the database containers
docker-down:
	@echo "🛑 Stopping database..."
	docker compose down

# ----------------------------
# 4. Utilities
# ----------------------------

# Remove build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...