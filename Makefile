.PHONY: install-tools test-coverage test watch-test start-server start-seeder

install-tools:
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Installing migrate CLI..."; \
		go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		echo "migrate CLI installed successfully"; \
	else \
		echo "migrate already installed"; \
	fi

	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing Air (live reload)..."; \
		go install github.com/cosmtrek/air@latest; \
		echo "Air installed successfully"; \
		echo "To use Air, ensure you have a .air.toml configuration file in your project root."; \
	else \
		echo "air already installed"; \
	fi

	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
		echo "gotestsum installed successfully"; \
	else \
		echo "gotestsum already installed"; \
	fi

	@if ! command -v reflex >/dev/null 2>&1; then \
		echo "Installing reflex (file watcher)..."; \
		go install github.com/cespare/reflex@latest; \
		echo "reflex installed successfully"; \
	else \
		echo "reflex already installed"; \
	fi

test: install-tools
	@echo "Running tests with gotestsum..."
	@gotestsum --format=short-verbose -- ./...
	@echo "✅ Tests completed."

test-coverage: install-tools
	@echo "Running tests with coverage..."
	@gotestsum -- -coverprofile=coverage.out ./...
	@echo "Generating coverage report..."
	@go tool cover -func=coverage.out | tee coverage-summary.txt
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage summary written to coverage-summary.txt"
	@echo "✅ Coverage HTML report generated at coverage.html"

watch-test: install-tools
	@echo "Watching for changes and running tests..."
	@reflex -r '\.go$$' -s -- sh -c 'clear && gotestsum --format=short-verbose -- ./...'

start-server: install-tools
	@echo "Detecting platform and starting Docker..."

	@if [ "$$(uname)" = "Darwin" ]; then \
		if pgrep -f "OrbStack" > /dev/null; then \
			echo "Opening OrbStack..."; \
			open -a OrbStack || echo "Failed to open OrbStack"; \
		else \
			echo "Opening Docker Desktop..."; \
			open -a Docker || echo "Failed to open Docker Desktop"; \
		fi \
	else \
		echo "Running on Linux - ensure Docker is running..."; \
	fi

	@echo "Waiting for Docker to be ready..."
	@until docker info > /dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@echo "\nDocker is ready."

	@echo "Starting Docker containers (detached)..."
	@docker-compose up -d

	@echo "Waiting for MySQL to be ready..."
	@until docker exec $$(docker ps -qf "name=mysql") \
		mysqladmin ping -h"127.0.0.1" --silent; do \
		printf "."; \
		sleep 1; \
	done
	@echo "\nMySQL is ready."

	@echo "Starting server with live-reload (air)..."
	@air

start-seeder:
	@echo "Seeding the database..."
	go run ./cmd/seeder/seeder.go
	@echo "Database seeding completed"
	