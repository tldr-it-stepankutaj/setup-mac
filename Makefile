# setup-mac Makefile

# Build variables
BINARY_NAME=setup-mac
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/stepankutaj/setup-mac/internal/cli.Version=$(VERSION) \
                  -X github.com/stepankutaj/setup-mac/internal/cli.Commit=$(COMMIT) \
                  -X github.com/stepankutaj/setup-mac/internal/cli.BuildDate=$(BUILD_DATE)"

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOMOD=$(GO) mod

.PHONY: all build clean test lint fmt deps run install help install-tools dist

# Default target
all: deps build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/setup-mac

# Build for release (multiple platforms)
build-release:
	@echo "Building release binaries..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/setup-mac
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/setup-mac

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ dist/
	$(GO) clean

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null 2>&1 || test -x "$$(go env GOPATH)/bin/golangci-lint" || (echo "golangci-lint not found. Run 'make install-tools' first." && exit 1)
	@if which golangci-lint > /dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		$$(go env GOPATH)/bin/golangci-lint run ./...; \
	fi

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Install globally
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) /usr/local/bin/

# Uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)

# Run dry-run
dry-run: build
	./bin/$(BINARY_NAME) install --all --dry-run

# Create distribution packages
dist: clean build-release
	@echo "Creating distribution packages..."
	@rm -rf dist/
	@for arch in amd64 arm64; do \
		DIST_NAME="$(BINARY_NAME)-darwin-$$arch"; \
		mkdir -p dist/$$DIST_NAME/bin; \
		cp bin/$(BINARY_NAME)-darwin-$$arch dist/$$DIST_NAME/bin/$(BINARY_NAME); \
		cp -r configs dist/$$DIST_NAME/; \
		cp README.md dist/$$DIST_NAME/; \
		echo '# setup-mac Installation Makefile' > dist/$$DIST_NAME/Makefile; \
		echo '' >> dist/$$DIST_NAME/Makefile; \
		echo 'BINARY_NAME=setup-mac' >> dist/$$DIST_NAME/Makefile; \
		echo 'INSTALL_DIR=/usr/local/bin' >> dist/$$DIST_NAME/Makefile; \
		echo 'CONFIG_DIR=$$(HOME)/.config/setup-mac' >> dist/$$DIST_NAME/Makefile; \
		echo '' >> dist/$$DIST_NAME/Makefile; \
		echo '.PHONY: install uninstall' >> dist/$$DIST_NAME/Makefile; \
		echo '' >> dist/$$DIST_NAME/Makefile; \
		echo 'install:' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Installing $$(BINARY_NAME)..."\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@mkdir -p $$(INSTALL_DIR)\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@cp bin/$$(BINARY_NAME) $$(INSTALL_DIR)/\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@chmod +x $$(INSTALL_DIR)/$$(BINARY_NAME)\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@mkdir -p $$(CONFIG_DIR)\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@cp configs/default.yaml $$(CONFIG_DIR)/\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Installed $$(BINARY_NAME) to $$(INSTALL_DIR)"\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Default config copied to $$(CONFIG_DIR)/default.yaml"\n' >> dist/$$DIST_NAME/Makefile; \
		echo '' >> dist/$$DIST_NAME/Makefile; \
		echo 'uninstall:' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Uninstalling $$(BINARY_NAME)..."\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@rm -f $$(INSTALL_DIR)/$$(BINARY_NAME)\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Removed $$(BINARY_NAME) from $$(INSTALL_DIR)"\n' >> dist/$$DIST_NAME/Makefile; \
		printf '\t@echo "Config files in $$(CONFIG_DIR) were not removed"\n' >> dist/$$DIST_NAME/Makefile; \
		tar -czvf dist/$$DIST_NAME.tar.gz -C dist $$DIST_NAME; \
	done
	@echo "Distribution packages created in dist/"

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Download dependencies and build (default)"
	@echo "  build         - Build the binary"
	@echo "  build-release - Build release binaries for multiple platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  fmt           - Format code"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  run           - Build and run"
	@echo "  install       - Install to /usr/local/bin"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  install-tools - Install development tools (golangci-lint)"
	@echo "  dist          - Create distribution packages (tar.gz)"
	@echo "  dry-run       - Build and run with --all --dry-run"
	@echo "  help          - Show this help"
