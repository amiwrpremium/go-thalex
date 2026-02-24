.PHONY: all build test test-integration test-verbose coverage lint fmt vet tidy check clean examples help version release

# Project metadata
MODULE := github.com/amiwrpremium/go-thalex
GO     := go
GOVET  := $(GO) vet
GOFMT  := gofmt
GOTEST := $(GO) test

# Default target
all: check

## build: Compile all packages
build:
	$(GO) build ./...

## test: Run tests with race detector and coverage
test:
	$(GOTEST) ./... -race -count=1 -coverprofile=coverage.out -covermode=atomic
	@echo "Coverage report: coverage.out"

## test-integration: Run integration tests (requires THALEX_PEM_PATH and THALEX_KEY_ID)
test-integration:
	$(GOTEST) -tags=integration ./... -race -count=1 -v -timeout=120s

## test-verbose: Run unit tests with verbose output
test-verbose:
	$(GOTEST) ./... -race -count=1 -v

## coverage: Open HTML coverage report in browser
coverage: test
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage HTML: coverage.html"
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || true

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format all Go source files
fmt:
	$(GOFMT) -w .
	goimports -w -local $(MODULE) .

## vet: Run go vet
vet:
	$(GOVET) ./...

## tidy: Tidy and verify module dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## check: Run all checks (format, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed."

## clean: Remove build artifacts and coverage files
clean:
	rm -f coverage.out coverage.html
	rm -rf dist/ bin/
	$(GO) clean ./...

## examples: Build all examples
examples:
	$(GO) build ./examples/...

# Version management
VERSION_FILE := config/network.go
CURRENT_VERSION := $(shell sed -n 's/^const Version = "\(.*\)"/\1/p' $(VERSION_FILE))

## version: Show current SDK version
version:
	@echo "$(CURRENT_VERSION)"

## release: Create a versioned release (usage: make release VERSION=0.2.0)
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=x.y.z"; \
		echo "Current version: $(CURRENT_VERSION)"; \
		exit 1; \
	fi
	@if ! echo "$(VERSION)" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$$'; then \
		echo "Error: VERSION must be in semver format (e.g., 0.2.0, 1.0.0-beta.1)"; \
		exit 1; \
	fi
	@echo "Releasing v$(VERSION) (current: $(CURRENT_VERSION))"
	@echo ""
	@echo "Updating $(VERSION_FILE)..."
	@sed -i.bak 's/const Version = "$(CURRENT_VERSION)"/const Version = "$(VERSION)"/' $(VERSION_FILE)
	@rm -f $(VERSION_FILE).bak
	@echo "Running checks..."
	@$(GO) build ./...
	@$(GOTEST) ./... -race -count=1
	@echo ""
	@echo "Committing version bump..."
	@git add $(VERSION_FILE)
	@git commit -m "chore: bump version to v$(VERSION)"
	@echo ""
	@echo "Creating tag v$(VERSION)..."
	@git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	@echo ""
	@echo "Done! To publish:"
	@echo "  git push origin master --tags"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
