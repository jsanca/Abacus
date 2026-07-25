.PHONY: fmt lint test coverage build run-server docker-up docker-down verify

fmt:
	cd server && gofmt -w .

lint:
	cd server && golangci-lint run ./...

test:
	cd server && go test ./...

coverage:
	cd server && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

build:
	cd server && go build -o bin/api ./cmd/api

run-server:
	cd server && go run ./cmd/api

docker-up:
	docker compose up -d

docker-down:
	docker compose down

verify:
	@echo "==> Checking formatting..."
	@if [ -n "$$(cd server && gofmt -l .)" ]; then \
		echo "Format issues found:"; \
		cd server && gofmt -l .; \
		exit 1; \
	fi
	@echo "==> Linting..."
	cd server && golangci-lint run ./...
	@echo "==> Running tests with race detector..."
	cd server && go test -race ./...
	@echo "==> Building..."
	cd server && go build -o /dev/null ./cmd/api
	@echo "==> All checks passed."
