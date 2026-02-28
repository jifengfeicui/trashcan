# 微信登录实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为现有用户体系添加微信登录支持，保留原有用户名密码登录方式

**Architecture:** 
- 后端 User 模型添加 openid、session_key、nickname、avatar 字段
- 新增 `/api/users/wechat-login` 接口处理微信登录
- 小程序前端增加微信登录按钮，调用 wx.login + wx.getUserProfile

**Tech Stack:** 
- 后端: Go + Gin + GORM
- 前端: 微信小程序原生开发

---

### Task 1: 修改 User 模型，添加微信字段

**Files:**
- Modify: `ginServer/model/user.go`

**Step 1: 修改 User 模型**

```go
// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Username     string    `json:"username" gorm:"type:TEXT;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:TEXT"` // 不返回给前端
	OpenID       string    `json:"-" gorm:"type:TEXT;uniqueIndex"` // 微信 openid
	SessionKey   string    `json:"-" gorm:"type:TEXT"` // 微信 session_key
	Nickname     string    `json:"nickname" gorm:"type:TEXT"` // 微信昵称
	Avatar       string    `json:"avatar" gorm:"type:TEXT"` // 微信头像
	WechatID     string    `json:"-" gorm:"type:TEXT"` // 微信统一ID
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**Step 2: Commit**

```bash
git add ginServer/model/user.go
git commit -m "feat: 添加微信相关字段到 User 模型"
```

---

### Task 2: 新增微信登录 API

**Files:**
- Modify: `ginServer/api/user.go`
- Config: 需在配置中添加微信 AppID 和 AppSecret

**Step 1: 添加微信登录请求结构体**

```go
// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`       // wx.login 获取的 code
	Nickname string `json:"nickname"`                      // 微信昵称（可选）
	Avatar   string `json:"avatar"`                        // 微信头像 URL（可选）
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	OpenID     string `json:"open_id"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"union_id"`
}
```

**Step 2: 添加微信登录处理函数**

```go
// WechatLogin 微信登录
// POST /api/users/wechat-login
func WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c)
		return
	}

	// 调用微信 API 获取 openid 和 session_key
	appID := global.CONFIG.WechatConfig.AppID
	appSecret := global.CONFIG.WechatConfig.AppSecret
	
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, req.Code)
	
	resp, err := http.Get(url)
	if err != nil {
		global.SugarLogger.Errorf("微信登录请求失败: %v", err)
		common.FailWithMessage("微信登录失败", c)
		return
	}
	defer resp.Body.Close()

	var wechatResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wechatResp); err != nil {
		global.SugarLogger.Errorf("解析微信响应失败: %v", err)
		common.FailWithMessage("微信登录失败", c)
		return
	}

	if wechatResp.ErrCode != 0 {
		global.SugarLogger.Errorf("微信 API 返回错误: %d - %s", wechatResp.ErrCode, wechatResp.ErrMsg)
		common.FailWithMessage("微信登录失败: "+wechatResp.ErrMsg, c)
		return
	}

	if wechatResp.OpenID == "" {
		common.FailWithMessage("无法获取微信 OpenID", c)
		return
	}

	// 查询或创建用户
	var user model.User
	result := global.DB.Where("open_id = ?", wechatResp.OpenID).First(&user)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		global.SugarLogger.Errorf("查询用户失败: %v", result.Error)
		common.FailWithMessage("登录失败", c)
		return
	}

	isNewUser := result.Error == gorm.ErrRecordNotFound

	// 更新用户信息
	updates := map[string]interface{}{
		"session_key": wechatResp.SessionKey,
		"open_id":     wechatResp.OpenID,
	}

	// 如果提供了昵称和头像且用户不存在或想更新，则更新
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if wechatResp.UnionID != "" {
		updates["wechat_id"] = wechatResp.UnionID
	}

	if isNewUser {
		// 新用户：创建用户
		// 如果有昵称则用昵称，否则用微信 ID
		username := req.Nickname
		if username == "" {
			username = "wechat_" + wechatResp.OpenID[:8]
		}
		updates["username"] = username

		user = model.User{
			Username:     username,
			OpenID:       wechatResp.OpenID,
			SessionKey:   wechatResp.SessionKey,
			Nickname:     req.Nickname,
			Avatar:       req.Avatar,
			WechatID:     wechatResp.UnionID,
		}

		if err := global.DB.Create(&user).Error; err != nil {
			global.SugarLogger.Errorf("创建微信用户失败: %v", err)
			common.FailWithMessage("登录失败", c)
			return
		}
	} else {
		// 老用户：更新信息
		if err := global.DB.Model(&user).Updates(updates).Error; err != nil {
			global.SugarLogger.Errorf("更新微信用户失败: %v", err)
			common.FailWithMessage("登录失败", c)
			return
		}
		// 重新获取用户信息
		global.DB.First(&user, user.ID)
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		global.SugarLogger.Errorf("生成 token 失败: %v", err)
		common.FailWithMessage("登录失败", c)
		return
	}

	resultData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
		"is_new": isNewUser,
	}

	if isNewUser {
		common.OkWithDetailed(resultData, "注册并登录成功", c)
	} else {
		common.OkWithDetailed(resultData, "登录成功", c)
	}
}
```

**Step 3: 添加 import**

```go
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	
	"gorm.io/gorm"
	// ...existing imports
)
```

**Step 4: Commit**

```bash
git add ginServer/api/user.go
git commit -m "feat: 添加微信登录 API"
```

---

### Task 3: 添加路由配置

**Files:**
- Modify: `ginServer/router/router.go`
- Config: 需在配置结构中添加 WechatConfig

**Step 1: 添加路由**

在 `v1.POST("/users/login", api.Login)` 后添加：

```go
v1.POST("/users/wechat-login", api.WechatLogin)
```

**Step 2: Commit**

```bash
git add ginServer/router/router.go
git commit -m "feat: 添加微信登录路由"
```

---

### Task 4: 添加微信配置

**Files:**
- 需查找配置结构定义位置

**Step 1: 查找配置结构**

运行: `grep -r "type.*Config" ginServer/ --include="*.go" | head -20`

找到配置结构后，添加：

```go
// WechatConfig 微信配置
type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}
```

并在配置中添加 `WechatConfig` 字段。

**Step 2: Commit**

```bash
git add ginServer/
git commit -m "feat: 添加微信配置支持"
```

---

### Task 5: 小程序前端 - 添加微信登录 API

**Files:**
- Modify: `miniapp/utils/api/user.js`

**Step 1: 添加微信登录函数**

```javascript
function wechatLogin(code, nickname, avatar) {
  return request({
    url: "/users/wechat-login",
    method: "POST",
    data: { 
      code, 
      nickname, 
      avatar 
    }
  })
}

module.exports = {
  register,
  login,
  getCurrentUser,
  wechatLogin
}
```

**Step 2: Commit**

```bash
git add miniapp/utils/api/user.js
git commit -m "feat: 添加小程序微信登录 API"
```

---

### Task 6: 小程序前端 - 登录页添加微信登录按钮

**Files:**
- Modify: `miniapp/pages/login/index.wxml`
- Modify: `miniapp/pages/login/index.wxss`
- Modify: `miniapp/pages/login/index.js`

**Step 1: 添加按钮到 wxml**

在表单提交按钮后添加：

```html
<button 
  class="wechat-btn" 
  bindtap="handleWechatLogin" 
  loading="{{wechatLoading}}"
>
  <image class="wechat-icon" src="/assets/wechat.png"></image>
  微信一键登录
</button>
```

**Step 2: 添加样式**

```css
.wechat-btn {
  background-color: #07C160;
  color: white;
  margin-top: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
}

.wechat-icon {
  width: 40rpx;
  height: 40rpx;
}
```

**Step 3: 添加处理函数**

```javascript
const { wechatLogin } = require("../../utils/api/user")

Page({
  data: {
    // ... existing data
    wechatLoading: false
  },

  handleWechatLogin() {
    if (this.data.wechatLoading) {
      return
    }

    this.setData({ wechatLoading: true, error: "" })

    // 第一步：wx.login 获取 code
    wx.login({
      success: (loginRes) => {
        // 第二步：获取用户信息（需要用户授权）
        wx.getUserProfile({
          desc: '用于完善用户资料',
          success: (profileRes) => {
            const userInfo = profileRes.userInfo || {}
            
            // 第三步：发送到后端
            wechatLogin(loginRes.code, userInfo.nickName, userInfo.avatarUrl)
              .then((res) => {
                const payload = res.data || {}
                setToken(payload.token || "")
                setUserInfo(payload.user || null)
                wx.showToast({ 
                  title: payload.is_new ? "注册成功" : "登录成功", 
                  icon: "success" 
                })
                this.redirectAfterAuth()
              })
              .catch((err) => {
                this.setData({ error: err.message || "微信登录失败" })
              })
              .finally(() => {
                this.setData({ wechatLoading: false })
              })
          },
          fail: (profileErr) => {
            // 用户拒绝授权，仍然可以用 code 登录（不传昵称头像）
            console.log('用户拒绝授权:', profileErr)
            wechatLogin(loginRes.code, "", "")
              .then((res) => {
                const payload = res.data || {}
                setToken(payload.token || "")
                setUserInfo(payload.user || null)
                wx.showToast({ 
                  title: payload.is_new ? "注册成功" : "登录成功", 
                  icon: "success" 
                })
                this.redirectAfterAuth()
              })
              .catch((err) => {
                this.setData({ error: err.message || "微信登录失败" })
              })
              .finally(() => {
                this.setData({ wechatLoading: false })
              })
          }
        })
      },
      fail: (loginErr) => {
        this.setData({ error: "微信登录失败", wechatLoading: false })
      }
    })
  }
})
```

**Step 4: Commit**

```bash
git add miniapp/pages/login/
git commit -m "feat: 小程序登录页添加微信登录按钮"
```

---

### Task 7: 整体测试

**Step 1: 启动后端服务**

```bash
cd ginServer && go run server.go
```

**Step 2: 使用微信开发者工具测试登录流程**

1. 点击「微信一键登录」按钮
2. 授权后应能成功登录
3. 再次登录应能识别为老用户

**Step 3: 验证数据**

检查数据库中用户表是否有 open_id、nickname、avatar 等字段。

---

### 总结

| Task | 描述 |
|------|------|
| 1 | User 模型添加微信字段 |
| 2 | 微信登录 API |
| 3 | 路由配置 |
| 4 | 微信配置 |
| 5 | 小程序 API |
| 6 | 登录页 UI + 逻辑 |
| 7 | 整体测试 |
