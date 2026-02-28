# 微信小程序版（miniapp）

本目录是将原 `front/`（Vue Web）迁移出的原生微信小程序实现，使用同一套后端 API。

## 已迁移功能

- 登录 / 注册
- 首页地图定位与附近垃圾桶搜索
- 垃圾桶点赞 / 点踩
- 打开系统地图导航
- 上传垃圾桶（定位/选点、地址、描述、图片）
- 用户中心（当前用户、我的上传分页、删除）

## 目录结构

```text
miniapp/
├── app.js
├── app.json
├── app.wxss
├── config/
│   └── index.js
├── utils/
│   ├── auth.js
│   ├── request.js
│   ├── geocoder.js
│   └── api/
│       ├── user.js
│       └── trashcan.js
└── pages/
    ├── home/
    ├── upload/
    ├── profile/
    └── login/
```

## 本地联调步骤

1. 启动后端（根目录）：

```powershell
go run main.go
```

2. 打开微信开发者工具，导入 `miniapp` 目录。

3. 修改 `miniapp/config/index.js`：

- `API_ORIGIN` 改成你后端可访问地址，如 `https://192.168.1.10:38080`。
- 保持 `API_BASE_URL = ${API_ORIGIN}/api`。

4. 开发者工具中勾选：

- `不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书`

## 线上配置注意

### 第 4 步：接入真实 AppID 与 request 合法域名

1. 修改小程序 AppID：

- 文件：`miniapp/project.config.json`
- 字段：`appid`
- 当前值是占位值 `__REPLACE_WITH_REAL_MINIAPP_APPID__`，请替换为你的小程序真实 AppID（形如 `wx1234567890abcdef`）。

2. 配置线上 API 域名：

- 文件：`miniapp/config/index.js`
- `API_ORIGIN` 改为你的线上域名，例如 `https://api.example.com`
- 保持 `API_BASE_URL = ${API_ORIGIN}/api`

3. 到微信公众平台配置服务器域名：

- 入口：`微信公众平台 -> 小程序 -> 开发管理 -> 开发设置 -> 服务器域名`
- 在 `request 合法域名` 中添加与 `API_ORIGIN` 一致的域名（只填域名，不带路径）

4. 证书与协议要求：

- 必须是 `HTTPS`
- 使用受信任 CA 证书（不要使用自签名证书）
- 后端证书域名需与 `API_ORIGIN` 一致

5. 发布前检查：

- 开发者工具中关闭“`不校验合法域名...`”后仍可正常请求
- 真机预览能正常访问 `GET /api/test`、登录、上传接口

### 其他注意事项

- 若使用高德逆地理编码，需保证 `AMAP_WEB_KEY` 对应服务权限正确。

## 已知差异

- Web 端的高德 JS Map 交互已替换为小程序原生 `<map>`。
- 用户中心当前保留“删除”，暂未迁移“编辑垃圾桶”。
