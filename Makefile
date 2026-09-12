.PHONY: help build run test clean migrate-up migrate-down migrate-create docker-up docker-down

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 编译后端应用
	@echo "Building backend..."
	@go build -o bin/api cmd/api/main.go

run: ## 运行后端应用
	@echo "Running backend..."
	@go run cmd/api/main.go

test: ## 运行测试
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## 运行测试并显示覆盖率
	@go tool cover -html=coverage.out

clean: ## 清理构建产物
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

migrate-up: ## 执行数据库迁移
	@echo "Running migrations..."
	@migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

migrate-create: ## 创建新的迁移文件 (usage: make migrate-create NAME=create_users_table)
	@migrate create -ext sql -dir migrations -seq $(NAME)

docker-up: ## 启动 Docker 容器
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## 停止 Docker 容器
	@echo "Stopping Docker containers..."
	@docker-compose down

lint: ## 运行代码检查
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...
