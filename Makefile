.PHONY: setup build test lint clean dev

# Default database URLs for local/CI environments
export DATABASE_URL ?= postgres://postgres:postgres@127.0.0.1:5432/buffalo_test_development?sslmode=disable
export TEST_DATABASE_URL ?= postgres://postgres:postgres@127.0.0.1:5432/buffalo_test_test?sslmode=disable

setup:
	@echo "=> Installing dependencies..."
	go mod download
	npm install
	@echo "=> Creating databases..."
	-soda create -e test
	-soda create -e development
	@echo "=> Running migrations..."
	soda migrate -e test
	soda migrate -e development

build:
	@echo "=> Building frontend assets..."
	npm run build
	@echo "=> Building Buffalo binary..."
	buffalo build --static -o bin/app

test:
	@echo "=> Cleaning schema.sql..."
	go run scripts/clean_schema.go
	@echo "=> Building frontend assets (required for templates)..."
	npm run build
	@echo "=> Running tests..."
	buffalo test

lint:
	@echo "=> Formatting Go code..."
	gofmt -w .

clean:
	rm -rf bin/
	rm -rf public/assets/

dev:
	@echo "=> Starting development servers..."
	# Run Buffalo and Vite concurrently
	npm run dev & buffalo dev
