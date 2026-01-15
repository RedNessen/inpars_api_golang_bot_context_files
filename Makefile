.PHONY: help build run clean test deps install

# Переменные
APP_NAME=inpars-telegram-bot
BUILD_DIR=bin
MAIN_PATH=./cmd/bot

help: ## Показать это сообщение помощи
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## Установить зависимости
	go mod download
	go mod tidy

build: ## Собрать приложение
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: ## Запустить приложение
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	@echo "Starting $(APP_NAME)..."
	@export $$(cat .env | xargs) && go run $(MAIN_PATH)

install: deps ## Установить зависимости и собрать приложение
	$(MAKE) build

clean: ## Удалить собранные файлы
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test: ## Запустить тесты
	go test -v ./...

fmt: ## Форматировать код
	go fmt ./...

lint: ## Проверить код линтером
	golangci-lint run

docker-build: ## Собрать Docker образ
	docker build -t $(APP_NAME):latest .

docker-run: ## Запустить в Docker
	docker run --env-file .env $(APP_NAME):latest
