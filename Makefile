.PHONY: build run test lint clean docker-build docker-run docker-compose-up docker-compose-down

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

build:
	go build -o bin/server cmd/api/main.go

run:
	@if [ -f .env ]; then \
		echo "Loading environment from .env file"; \
		zsh -c 'source .env && go run cmd/api/main.go'; \
	else \
		echo ".env file not found. Running without it."; \
		go run cmd/api/main.go; \
	fi

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

docker-build:
	docker build -t pairs-trading-backend .

docker-run:
	docker run -p 8080:8080 pairs-trading-backend

docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down