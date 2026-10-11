BACKEND_DIR := backend
BIN := $(BACKEND_DIR)/bin/enterprise-admin

.PHONY: gen build run tidy help

help: ## 查看可用命令
	@grep -E '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-8s %s\n", $$1, $$2}'

gen: ## 生成 gorm 类型安全查询包（改 models/ 后重跑）
	cd $(BACKEND_DIR) && go run ./cmd/gen

build: ## 编译后端二进制到 backend/bin/enterprise-admin
	cd $(BACKEND_DIR) && go build -o bin/enterprise-admin .

run: ## 本地运行后端（默认 sqlite :8080）
	cd $(BACKEND_DIR) && go run .

tidy: ## 整理 Go 依赖
	cd $(BACKEND_DIR) && go mod tidy
