# Build variables
BINARY_NAME=csv-tools
MODULE_NAME=github.com/3-2-1-contact/radio-scan-list

# Git information
VERSION ?= $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "dev")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-X $(MODULE_NAME)/internal/version.Version=$(VERSION) \
                  -X $(MODULE_NAME)/internal/version.CommitHash=$(COMMIT_HASH) \
                  -X $(MODULE_NAME)/internal/version.BuildTime=$(BUILD_TIME)"

.PHONY: all clean build

all: build

build:
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	@echo "Commit: $(COMMIT_HASH)"
	@echo "Build time: $(BUILD_TIME)"
	go build $(LDFLAGS) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_HASH)"
	@echo "Build time: $(BUILD_TIME)"

# Optional: add install target
install:
	go install $(LDFLAGS)

# Optional: add test target
test:
	go test ./...
