.PHONY: all build clean test test-coverage fmt lint install-tools mod-download mod-tidy run-local package deploy help

# Variables
BINARY_NAME=bootstrap
LAMBDA_ZIP=lambda-deployment.zip
GO=go
GOFMT=gofmt
GOVET=go vet
GOLINT=golangci-lint
GOTEST=go test
GOBUILD=go build
GOCLEAN=go clean
GOMOD=go mod
LAMBDA_RUNTIME=provided.al2023
MAIN_PATH=cmd/lambda/main.go
BUILD_DIR=build
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# AWS Lambda specific build flags
LAMBDA_BUILD_FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=arm64

# Default target
all: clean fmt lint test build

# Help target
help:
	@echo "Available targets:"
	@echo "  make all            - Run all checks and build"
	@echo "  make build          - Build the Lambda binary"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run unit tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linters"
	@echo "  make install-tools  - Install required tools"
	@echo "  make mod-download   - Download Go modules"
	@echo "  make mod-tidy       - Tidy Go modules"
	@echo "  make package        - Create Lambda deployment package"
	@echo "  make run-local      - Run Lambda locally (requires environment variables)"

# Install required tools
install-tools:
	@echo "Installing required tools..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed successfully"

# Download dependencies
mod-download:
	@echo "Downloading dependencies..."
	@$(GOMOD) download

# Tidy dependencies
mod-tidy:
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -s -w .
	@goimports -w .
	@echo "Code formatted successfully"

# Run linters
lint:
	@echo "Running linters..."
	@if ! command -v $(GOLINT) &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@$(GOLINT) run --timeout=5m ./...
	@echo "Linting completed successfully"

# Run unit tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v -race -short ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Coverage summary:"
	@$(GO) tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "Total Coverage: " $$3}'

# Verify coverage meets minimum threshold
verify-coverage: test-coverage
	@echo "Verifying coverage threshold (80%)..."
	@coverage=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage < 80" | bc -l) -eq 1 ]; then \
		echo "Coverage $$coverage% is below threshold of 80%"; \
		exit 1; \
	else \
		echo "Coverage $$coverage% meets threshold"; \
	fi

# Build the Lambda binary
build: clean mod-download
	@echo "Building Lambda function..."
	@mkdir -p $(BUILD_DIR)
	@$(LAMBDA_BUILD_FLAGS) $(GOBUILD) -ldflags="-s -w" -tags lambda.norpc -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Create Lambda deployment package
package: build
	@echo "Creating Lambda deployment package..."
	@cd $(BUILD_DIR) && zip -j ../$(LAMBDA_ZIP) $(BINARY_NAME)
	@echo "Lambda package created: $(LAMBDA_ZIP)"
	@echo "Package size: $$(du -h $(LAMBDA_ZIP) | cut -f1)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(LAMBDA_ZIP)
	@rm -f $(COVERAGE_FILE)
	@rm -f $(COVERAGE_HTML)
	@echo "Clean completed"

# Run Lambda locally (for testing)
run-local:
	@echo "Running Lambda locally..."
	@if [ -z "$$STYTCH_PROJECT_ID" ] || [ -z "$$STYTCH_SECRET" ]; then \
		echo "Error: STYTCH_PROJECT_ID and STYTCH_SECRET environment variables are required"; \
		exit 1; \
	fi
	@$(GO) run $(MAIN_PATH)

# Validate Go code
validate:
	@echo "Validating Go code..."
	@$(GOVET) ./...
	@$(GO) mod verify
	@echo "Validation completed successfully"

# Full CI pipeline
ci: clean mod-tidy fmt lint test-coverage verify-coverage build
	@echo "CI pipeline completed successfully"

# Quick check before commit
pre-commit: fmt lint test
	@echo "Pre-commit checks passed"

# Deploy to AWS Lambda (requires AWS CLI configured)
deploy: package
	@echo "Deploying to AWS Lambda..."
	@echo "Note: This requires AWS CLI configured and appropriate permissions"
	@echo "You would typically use: aws lambda update-function-code --function-name <your-function> --zip-file fileb://$(LAMBDA_ZIP)"
	@echo "Or use your infrastructure as code tool (Terraform, CloudFormation, etc.)"

# Show environment info
info:
	@echo "Environment Information:"
	@echo "  Go Version: $$(go version)"
	@echo "  GOOS: $$(go env GOOS)"
	@echo "  GOARCH: $$(go env GOARCH)"
	@echo "  CGO_ENABLED: $$(go env CGO_ENABLED)"
	@echo "  Module: $$(grep '^module' go.mod | awk '{print $$2}')"
	@echo "  Lambda Runtime: $(LAMBDA_RUNTIME)"