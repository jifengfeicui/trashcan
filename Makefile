# Makefile - 跨平台构建脚本

.PHONY: build clean frontend backend

# 默认目标
all: build

# 构建前端和后端
build: frontend backend

# 构建前端
frontend:
	@echo "构建前端..."
	@cd front && \
	if command -v pnpm >/dev/null 2>&1; then \
		pnpm install && pnpm build; \
	elif command -v npm >/dev/null 2>&1; then \
		npm install && npm run build; \
	else \
		echo "错误: 未找到pnpm或npm"; \
		exit 1; \
	fi
	@if [ ! -d "front/dist" ]; then \
		echo "错误: 前端构建输出目录 front/dist 不存在"; \
		exit 1; \
	fi
	@echo "复制前端构建文件到embed目录..."
	@rm -rf ginServer/static/dist
	@cp -r front/dist ginServer/static/dist

# 构建后端
backend:
	@echo "构建Go二进制文件..."
	@BUILD_TIME=$$(date '+%Y-%m-%d %H:%M:%S'); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	go build -ldflags "-X 'main.BuildTime=$$BUILD_TIME' -X 'main.GitCommit=$$GIT_COMMIT'" -o trashcan main.go
	@echo "构建完成！"
	@echo "可执行文件: ./trashcan"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf front/dist
	@rm -rf ginServer/static/dist
	@rm -f trashcan trashcan.exe
	@echo "清理完成！"

# 仅构建前端（开发用）
dev-frontend:
	@cd front && \
	if command -v pnpm >/dev/null 2>&1; then \
		pnpm dev; \
	else \
		npm run dev; \
	fi

# 仅运行后端（开发用）
dev-backend:
	@go run main.go

