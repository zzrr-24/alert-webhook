.PHONY: build test run clean docker-build docker-run

BINARY_NAME=alert-webhook
BUILD_DIR=build
DOCKER_IMAGE=registry.cn-hangzhou.aliyuncs.com/zzrr_images/adapter
DOCKER_TAG=20260326

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/webhook
	@cp configs/config.yaml.example $(BUILD_DIR)/config.yaml
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	@go test -v -cover ./...

run:
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd/webhook

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete"

docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker build complete"

docker-run:
	@echo "Running Docker container..."
	@docker run -d -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)
