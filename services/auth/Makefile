.PHONY: help build run test clean generate tidy

# Переменные
BINARY_NAME=auth-service
BUILD_DIR=.
API_DIR=api

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Обновить зависимости
	go mod tidy

generate: ## Сгенерировать код из OpenAPI спецификации
	cd $(API_DIR) && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml auth.yaml

generate-go: ## Сгенерировать код используя go:generate директиву
	cd $(API_DIR) && go generate

build: tidy generate ## Собрать приложение
	go build -o $(BINARY_NAME) ./cmd

run: build ## Запустить сервис
	./$(BINARY_NAME)

dev: build ## Запустить сервис в режиме разработки
	./$(BINARY_NAME) -port 8081

test: ## Запустить тесты
	go test ./...

clean: ## Очистить собранные файлы
	rm -f $(BINARY_NAME)
	rm -f $(API_DIR)/*.gen.go

install-tools: ## Установить инструменты разработки
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

docker-build: ## Собрать Docker образ
	docker build -t agromeration/auth-service .

docker-run: ## Запустить в Docker
	docker run -p 8081:8081 agromeration/auth-service

# Команды для тестирования API
test-health: ## Тестировать health endpoint
	curl -s http://localhost:8081/health | jq .

test-register: ## Тестировать регистрацию
	curl -s -X POST http://localhost:8081/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe"}' | jq .

test-login: ## Тестировать вход
	curl -s -X POST http://localhost:8081/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq .

test-all: test-health test-register test-login ## Запустить все тесты API 