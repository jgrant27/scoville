BINARY_NAME=scoville
BUILD_DIR=bin
.PHONY: all build run run-live install clean

install:
	@echo "📦 Installing Dependencies..."
	go mod download
	go mod tidy

build: install
	@echo "🌶️  Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/scoville/main.go

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

run-live: build
	@echo "🌶️  RUNNING IN LIVE MODE! 🌶️"
	PAPER_MODE=false ./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR) internal/state/scoville_progress.json
