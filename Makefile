.PHONY: build run test lint clean docker-build docker-run docker-compose-up docker-compose-down

# Import environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Build the application
build:
	go build -o bin/server cmd/api/main.go

# Run the application with environment variables
run:
	@if [ -f .env ]; then \
		echo "Loading environment from .env file"; \
		zsh -c 'source .env && go run cmd/api/main.go'; \
	else \
		echo ".env file not found. Running without it."; \
		go run cmd/api/main.go; \
	fi

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Clean built binaries
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t pairs-trading-backend .

# Run Docker container
docker-run:
	docker run -p 8080:8080 pairs-trading-backend

# Run the application using docker-compose
docker-compose-up:
	docker-compose up --build

# Stop docker-compose services
docker-compose-down:
	docker-compose down