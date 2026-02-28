# Makefile - 后端构建脚本

.PHONY: build clean

# 默认目标
all: build

# 构建后端
build:
	@echo "构建Go二进制文件..."
	@BUILD_TIME=$$(date '+%Y-%m-%d %H:%M:%S'); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	go build -ldflags "-X 'main.BuildTime=$$BUILD_TIME' -X 'main.GitCommit=$$GIT_COMMIT'" -o trashcan main.go
	@echo "构建完成！"
	@echo "可执行文件: ./trashcan"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -f trashcan trashcan.exe
	@echo "清理完成！"

# 运行后端（开发用）
dev:
	@go run main.go
