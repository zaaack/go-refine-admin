BACKEND_DIR := backend
FRONTEND_DIR := frontend
BIN := $(BACKEND_DIR)/bin/enterprise-admin

.PHONY: gen build run tidy dev-backend dev-frontend dev help

help: ## 查看可用命令
	@grep -E '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-13s %s\n", $$1, $$2}'

gen: ## 生成 gorm 类型安全查询包（改 models/ 后重跑）
	cd $(BACKEND_DIR) && go run ./cmd/gen

build: ## 编译后端二进制到 backend/bin/enterprise-admin
	cd $(BACKEND_DIR) && go build -o bin/enterprise-admin .

run: ## 本地运行后端（默认 sqlite :8080）
	cd $(BACKEND_DIR) && go run .

tidy: ## 整理 Go 依赖
	cd $(BACKEND_DIR) && go mod tidy

dev-backend: ## 启动后端开发服务（默认 sqlite :8080）
	cd $(BACKEND_DIR) && go run .

dev-frontend: ## 启动前端开发服务（:5173，代理 /api 到 :8080）
	cd $(FRONTEND_DIR) && pnpm dev

dev: ## 同时启动前后端开发服务
	$(MAKE) -j2 dev-backend dev-frontend
