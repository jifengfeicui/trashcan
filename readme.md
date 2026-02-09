# 垃圾桶管理系统

前后端分离项目，使用Go + Gin作为后端，Vue3 + Vite作为前端。

## 项目结构

- `front/` - Vue3前端项目
- `ginServer/` - Go后端服务
- `config.yml` - 配置文件

## 开发环境设置

### Go环境配置

执行以下命令，允许3.23 http访问：
```
go env -w GOPRIVATE=192.168.3.23/*
go env -w GOINSECURE=192.168.3.23/*
```

## 运行项目

### 1. 启动后端服务

```bash
go run main.go
```

后端服务默认运行在 `http://localhost:38080`

### 2. 启动前端开发服务器

```bash
cd front
pnpm install  # 首次运行需要安装依赖
pnpm dev
```

前端开发服务器默认运行在 `http://localhost:5173`

**局域网访问配置：**

前端已配置为允许局域网访问（监听 `0.0.0.0`），可以通过以下方式访问：

1. **使用localhost访问**（默认）：
   ```bash
   pnpm dev
   ```
   访问地址：`http://localhost:5173`
   后端地址：`http://localhost:38080`

2. **使用局域网IP访问**（推荐）：
   ```bash
   pnpm dev:local
   ```
   或者：
   ```bash
   vite --config vite.config.local.js
   ```
   启动时会自动检测并显示你的局域网IP地址。
   访问地址：`http://<你的局域网IP>:5173`（例如：`http://192.168.1.100:5173`）
   后端地址：`http://<你的局域网IP>:38080`
   
   注意：使用局域网IP访问时，确保后端也配置为监听 `0.0.0.0`（已在config.yml中配置）

3. **自定义后端地址**：
   ```bash
   # Windows (PowerShell)
   $env:VITE_API_URL="http://192.168.1.100:38080"; pnpm dev
   
   # Linux/Mac
   VITE_API_URL=http://192.168.1.100:38080 pnpm dev
   ```

前端通过Vite代理自动转发 `/api` 和 `/uploads` 请求到后端服务器。

## API接口

- `GET /api/trashcans/nearby` - 获取附近的垃圾桶
- `POST /api/trashcans` - 创建垃圾桶
- `GET /api/trashcans/:id` - 获取垃圾桶详情

## 生产环境部署

### 构建前端

```bash
cd front
pnpm build
```

构建产物在 `front/dist/` 目录，可以部署到静态文件服务器（如Nginx）。

### 后端配置

确保后端CORS配置允许前端域名访问。

