.PHONY: build test fmt clean run

BUILD_DIR := bin
BINARY_NAME := ue
MAIN_PATH := ./cmd/ue/main.go

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test:
	@echo "Running tests..."
	go test ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
