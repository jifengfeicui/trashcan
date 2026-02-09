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
- **静态文件嵌入**：Go embed（前端构建产物嵌入到二进制文件）

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
│   ├── middle/             # 中间件
│   └── static/             # 静态文件嵌入（前端构建产物）
│       ├── embed.go        # embed定义
│       └── dist/           # 前端构建输出（构建时生成）
├── initialize/             # 初始化模块
│   ├── orm/               # ORM配置
│   └── zapLog.go          # 日志初始化
├── utils/                  # 工具函数
├── config.yml             # 配置文件
├── Taskfile.yml           # Task构建配置文件
├── main.go                # 程序入口
└── go.mod                 # Go依赖管理
```

## 快速开始

### 环境要求

- Go 1.24+ 或 Go 1.22.12+
- Node.js 20.19.0+ 或 22.12.0+
- Task（构建工具，用于单二进制部署）
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

本项目使用 [Task](https://taskfile.dev/) 作为构建工具，实现前端编译后嵌入到后端二进制文件中，支持单二进制部署。

#### 安装 Task

**Windows (使用 Scoop):**
```powershell
scoop install task
```

**macOS (使用 Homebrew):**
```bash
brew install go-task/tap/go-task
```

**Linux:**
```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
```

或从 [Task 官网](https://taskfile.dev/installation/) 下载对应平台的二进制文件。

#### 使用 Task 构建

**完整构建（推荐）：**
```bash
task build
```

这会自动完成以下步骤：
1. 构建前端项目
2. 将前端构建产物复制到 embed 目录
3. 编译 Go 二进制文件（前端已嵌入）

**其他常用命令：**
```bash
# 检查构建环境
task check

# 仅构建前端
task frontend:only

# 仅构建后端（需要前端已构建）
task backend:only

# 清理构建产物
task clean

# 生成SSL证书（用于HTTPS）
task cert:generate

# 开发模式：启动前端开发服务器
task dev:frontend

# 开发模式：启动后端服务器
task dev:backend
```

#### 运行部署

构建完成后，直接运行生成的二进制文件即可：

**Windows:**
```powershell
.\trashcan.exe
```

**Linux/macOS:**
```bash
./trashcan
```

前端静态文件已嵌入到二进制文件中，无需单独部署前端文件。访问 `http://localhost:38080` 即可使用完整功能。

#### HTTPS 配置

项目支持 HTTPS 模式，可以通过配置文件启用：

1. **生成自签名证书（开发/测试用）：**
```bash
task cert:generate
```

2. **配置 HTTPS：**

编辑 `config.yml` 文件：
```yaml
gin:
  host: 0.0.0.0
  port: 38080
  enable_https: true      # 启用HTTPS
  cert_file: "cert.pem"   # SSL证书文件路径
  key_file: "key.pem"     # SSL私钥文件路径
```

3. **启动服务器：**

服务器将自动使用 HTTPS 模式启动，访问 `https://localhost:38080`。

**注意：**
- 自签名证书浏览器会显示安全警告，这是正常的
- 生产环境请使用由 CA（如 Let's Encrypt）签发的正式证书
- 如果使用正式证书，将证书和私钥文件路径配置到 `cert_file` 和 `key_file` 即可

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
- `gin.enable_https`: 是否启用HTTPS（默认：false）
- `gin.cert_file`: SSL证书文件路径（默认：cert.pem）
- `gin.key_file`: SSL私钥文件路径（默认：key.pem）
- `amap.api_key`: 高德地图 API Key
- `upload.image_dir`: 图片上传存储目录
- `jwt.secret`: JWT密钥（生产环境请修改）
- `jwt.expire_hours`: Token过期时间（小时）

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

