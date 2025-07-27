PHONY: build run test clean docker-build docker-run docker-stop deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=spotify-mcp-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) verify

# Initialize module (run this first if go.mod doesn't exist)
init:
	$(GOMOD) init github.com/your-org/spotify-mcp-server
	$(GOMOD) tidy

# Docker commands
docker-build:
	docker build -f docker/Dockerfile -t spotify-mcp-server:latest .

docker-build-no-cache:
	docker build --no-cache -f docker/Dockerfile -t spotify-mcp-server:latest .

docker-run:
	docker-compose up -d

docker-run-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

docker-stop:
	docker-compose down

docker-restart:
	docker-compose restart

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --rmi all --remove-orphans

# Docker management
docker-shell:
	docker-compose exec spotify-mcp-server sh

docker-stats:
	docker stats spotify-mcp-server

docker-inspect:
	docker inspect spotify-mcp-server

# Multi-platform build (for production)
docker-build-multi:
	docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile -t spotify-mcp-server:latest . --push

# Local registry
docker-registry:
	docker run -d -p 5000:5000 --restart=always --name registry registry:2

docker-push-local:
	docker tag spotify-mcp-server:latest localhost:5000/spotify-mcp-server:latest
	docker push localhost:5000/spotify-mcp-server:latest

# Development commands
dev:
	air -c .air.toml

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Security scan
security:
	gosec ./...