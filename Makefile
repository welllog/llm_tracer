# Makefile for LLM Proxy & Tracer (llm_tracer)

.PHONY: all build build-frontend build-backend run test dev-frontend clean

# 默认构建前端并编译后端
all: build

# 1. 一键构建前端与后端
build: build-frontend build-backend

# 2. 前端打包
build-frontend:
	@echo "Building React frontend..."
	cd frontend && npm run build

# 3. 后端编译
build-backend:
	@echo "Compiling Go backend..."
	go build -o llm_tracer_web

# 4. 运行本地代理服务
run: build-backend
	@echo "Starting proxy server..."
	./llm_tracer_web -listen :1238

# 5. 运行 Go 集成测试
test:
	@echo "Running integration tests..."
	go test -v ./...

# 6. 前端 Vite 开发服务器 (支持 HMR 热更新)
dev-frontend:
	@echo "Starting Vite dev server..."
	cd frontend && npm run dev

# 7. 清理构建产物与数据库
clean:
	@echo "Cleaning up..."
	rm -f llm_tracer_web
	rm -rf static/assets
	rm -f static/index.html
	rm -f static/placeholder.txt
