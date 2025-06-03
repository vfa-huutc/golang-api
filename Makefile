.PHONY: test-coverage install-tools

GOCOV := $(shell command -v gocov 2> /dev/null)
GOCOV_HTML := $(shell command -v gocov-html 2> /dev/null)

install-tools:
ifndef GOCOV
	@echo "Installing gocov..."
	@go install github.com/axw/gocov/gocov@latest
endif
ifndef GOCOV_HTML
	@echo "Installing gocov-html..."
	@go install github.com/matm/gocov-html@latest
endif

test-coverage: install-tools
	@echo "Running tests and generating coverage.json..."
	@gocov test ./... > coverage.json
	@echo "Generating HTML coverage report..."
	@gocov-html < coverage.json > coverage.html
	@echo "✅ Coverage report generated at coverage.html"

