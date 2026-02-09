# 垃圾桶定位系统

一个基于高德地图的垃圾桶定位系统，帮助用户快速查找附近的垃圾桶位置，支持查看详细信息、导航功能，以及用户上传新的垃圾桶位置和图片。

## 功能特性

- 🗺️ **地图展示**：基于高德地图，直观展示垃圾桶位置
- 🔍 **附近搜索**：根据当前位置搜索附近的垃圾桶，支持自定义搜索半径和返回数量
- 📍 **位置定位**：自动定位用户当前位置
- 📸 **图片展示**：支持查看垃圾桶的实景图片
- ➕ **位置上传**：用户可以在地图上点击添加新的垃圾桶位置，并上传图片
- 🧭 **导航功能**：集成高德地图导航，快速到达目标垃圾桶
- 📱 **响应式设计**：支持桌面端和移动端访问

## 技术栈

### 后端
- **语言**：Go 1.24+
- **框架**：Gin
- **数据库**：SQLite (GORM)
- **日志**：Zap
- **配置管理**：Viper

### 前端
- **框架**：Vue 3
- **构建工具**：Vite
- **状态管理**：Pinia
- **路由**：Vue Router
- **HTTP客户端**：Axios
- **地图服务**：高德地图 API

## 项目结构

```
trashcan/
├── front/                    # 前端项目
│   ├── src/
│   │   ├── api/            # API接口定义
│   │   ├── components/     # Vue组件
│   │   ├── views/          # 页面视图
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia状态管理
│   │   └── utils/          # 工具函数
│   ├── package.json
│   └── vite.config.js
├── ginServer/               # 后端服务
│   ├── api/                # API处理器
│   ├── model/              # 数据模型
│   ├── router/             # 路由定义
│   └── middle/             # 中间件
├── initialize/             # 初始化模块
│   ├── orm/               # ORM配置
│   └── zapLog.go          # 日志初始化
├── utils/                  # 工具函数
├── config.yml             # 配置文件
├── main.go                # 程序入口
└── go.mod                 # Go依赖管理
```

## 快速开始

### 环境要求

- Go 1.24+ 或 Go 1.22.12+
- Node.js 20.19.0+ 或 22.12.0+
- 高德地图 API Key

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd trashcan
```

2. **配置后端**

编辑 `config.yml` 文件，配置高德地图 API Key：

```yaml
gin:
  host: 0.0.0.0
  port: 38080

amap:
  api_key: "你的高德地图API Key"  # 请替换为你的高德地图API Key

upload:
  image_dir: "uploads/trashcans"  # 图片上传目录
```

3. **安装后端依赖**
```bash
go mod download
```

4. **安装前端依赖**
```bash
cd front
npm install
# 或
pnpm install
```

5. **初始化数据库**

数据库会在首次运行时自动创建。如果需要初始化测试数据，可以运行：

```bash
go run scripts/init_db.go
go run scripts/insert_test_data.go
```

6. **启动后端服务**
```bash
go run main.go
```

后端服务将在 `http://localhost:38080` 启动。

7. **启动前端开发服务器**
```bash
cd front
npm run dev
# 或
pnpm dev
```

前端开发服务器通常会在 `http://localhost:5173` 启动（具体端口以实际输出为准）。

### 生产环境部署

1. **构建前端**
```bash
cd front
npm run build
# 或
pnpm build
```

2. **配置后端静态文件服务**

确保后端配置了静态文件服务，将前端构建产物部署到后端可访问的目录。

3. **运行后端**
```bash
go build -o trashcan main.go
./trashcan
```

## API 接口

### 获取附近垃圾桶
```
GET /api/trashcans/nearby
参数：
  - lat: 纬度（必填）
  - lng: 经度（必填）
  - radius: 搜索半径（公里，默认5）
  - limit: 返回数量限制（默认10）
```

### 创建垃圾桶
```
POST /api/trashcans
表单数据：
  - latitude: 纬度（必填）
  - longitude: 经度（必填）
  - address: 地址（可选）
  - description: 描述（可选）
  - image: 图片文件（可选）
```

### 获取垃圾桶详情
```
GET /api/trashcans/:id
```

## 配置说明

### 后端配置 (config.yml)

- `gin.host`: 服务器监听地址
- `gin.port`: 服务器端口
- `amap.api_key`: 高德地图 API Key
- `upload.image_dir`: 图片上传存储目录

### 前端配置

前端的高德地图 API Key 配置在 `front/vite.config.js` 中，需要与后端配置保持一致。

## 开发说明

### 添加新的垃圾桶

1. 访问上传页面
2. 在地图上点击选择位置，或手动输入经纬度
3. 填写地址和描述（可选）
4. 上传垃圾桶图片（可选）
5. 提交表单

### 搜索附近垃圾桶

1. 点击"定位到我的位置"按钮获取当前位置
2. 设置搜索半径和返回数量
3. 点击"搜索附近垃圾桶"
4. 在地图上查看结果，点击标记查看详细信息
5. 使用导航功能规划路线

## 注意事项

- 使用本系统需要申请高德地图 API Key，并配置相应的服务权限（Web服务API、Web端JS API）
- 图片上传功能需要确保 `uploads/trashcans` 目录有写入权限
- 生产环境部署时，建议配置 HTTPS 以支持浏览器定位功能
- SQLite 数据库文件 `sqlite.db` 会在首次运行时自动创建

## 许可证

[添加许可证信息]

## 贡献

欢迎提交 Issue 和 Pull Request！

