#!/bin/bash
# Shell构建脚本 - 用于Linux/Mac
# 构建前端并嵌入到Go二进制文件中

set -e

echo "开始构建..."

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: 未找到Node.js，请先安装Node.js"
    exit 1
fi

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go，请先安装Go"
    exit 1
fi

# 进入前端目录
cd front

# 安装依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    if command -v pnpm &> /dev/null; then
        pnpm install
    elif command -v npm &> /dev/null; then
        npm install
    else
        echo "错误: 未找到pnpm或npm"
        exit 1
    fi
fi

# 构建前端
echo "构建前端..."
if command -v pnpm &> /dev/null; then
    pnpm build
else
    npm run build
fi

# 返回项目根目录
cd ..

# 检查dist目录是否存在
if [ ! -d "front/dist" ]; then
    echo "错误: 前端构建输出目录 front/dist 不存在"
    exit 1
fi

# 将dist目录复制到ginServer/static/dist（embed需要）
echo "复制前端构建文件到embed目录..."
if [ -d "ginServer/static/dist" ]; then
    rm -rf ginServer/static/dist
fi
cp -r front/dist ginServer/static/dist

# 构建Go二进制文件
echo "构建Go二进制文件..."
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

go build -ldflags "-X 'main.BuildTime=$BUILD_TIME' -X 'main.GitCommit=$GIT_COMMIT'" -o trashcan main.go

echo "构建完成！"
echo "可执行文件: ./trashcan"

