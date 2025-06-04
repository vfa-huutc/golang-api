.PHONY: install-tools test-coverage start-server

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

	@if ! command -v gocov >/dev/null 2>&1; then \
		echo "Installing gocov..."; \
		go install github.com/axw/gocov/gocov@latest; \
		echo "gocov installed successfully"; \
	else \
		echo "gocov already installed"; \
	fi

	@if ! command -v gocov-html >/dev/null 2>&1; then \
		echo "Installing gocov-html..."; \
		go install github.com/matm/gocov-html@latest; \
		echo "gocov-html installed successfully"; \
	else \
		echo "gocov-html already installed"; \
	fi

test-coverage: install-tools
	@echo "Running tests and generating coverage.json..."
	@gocov test ./... > coverage.json
	@echo "Generating HTML coverage report..."
	@gocov-html < coverage.json > coverage.html
	@echo "✅ Coverage report generated at coverage.html"

start-server: install-tools
	@echo "Starting Docker containers (detached)..."
	@docker-compose up -d
	@echo "Starting server with live-reload (air)..."
	@air