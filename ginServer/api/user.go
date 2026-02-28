package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"template/ginServer/api/common"
	"template/ginServer/model"
	"template/global"
	"template/utils"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"` // wx.login 获取的 code
	Nickname string `json:"nickname"`                // 微信昵称（可选）
	Avatar   string `json:"avatar"`                  // 微信头像 URL（可选）
}

// Register 用户注册
// POST /api/users/register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c)
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := global.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		common.FailWithMessage("用户名已存在", c)
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		global.SugarLogger.Errorf("密码加密失败: %v", err)
		common.FailWithMessage("注册失败", c)
		return
	}

	// 创建用户
	user := model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := global.DB.Create(&user).Error; err != nil {
		global.SugarLogger.Errorf("创建用户失败: %v", err)
		common.FailWithMessage("注册失败", c)
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		global.SugarLogger.Errorf("生成token失败: %v", err)
		common.FailWithMessage("注册成功，但登录失败", c)
		return
	}

	result := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	}

	common.OkWithDetailed(result, "注册成功", c)
}

// Login 用户登录
// POST /api/users/login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c)
		return
	}

	// 查找用户
	var user model.User
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		common.FailWithMessage("用户名或密码错误", c)
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		common.FailWithMessage("用户名或密码错误", c)
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		global.SugarLogger.Errorf("生成token失败: %v", err)
		common.FailWithMessage("登录失败", c)
		return
	}

	result := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	}

	common.OkWithDetailed(result, "登录成功", c)
}

// GetCurrentUser 获取当前用户信息
// GET /api/users/me
func GetCurrentUser(c *gin.Context) {
	// 从中间件中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}

	var user model.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		common.FailWithMessage("用户不存在", c)
		return
	}

	result := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
	}

	common.OkWithData(result, c)
}

// WechatLogin 微信登录
// POST /api/users/wechat-login
func WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c)
		return
	}

	if global.CONFIG.WechatConfig == nil {
		common.FailWithMessage("微信配置未启用", c)
		return
	}

	appID := global.CONFIG.WechatConfig.AppID
	appSecret := global.CONFIG.WechatConfig.AppSecret

	if appID == "" || appSecret == "" {
		common.FailWithMessage("微信配置未正确设置", c)
		return
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, req.Code)

	resp, err := http.Get(url)
	if err != nil {
		global.SugarLogger.Errorf("微信登录请求失败: %v", err)
		common.FailWithMessage("微信登录失败", c)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var wechatResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &wechatResp); err != nil {
		global.SugarLogger.Errorf("解析微信响应失败: %v", err)
		global.SugarLogger.Errorf("微信响应内容: %s", string(body))
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

	var user model.User
	result := global.DB.Where("open_id = ?", wechatResp.OpenID).First(&user)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		global.SugarLogger.Errorf("查询用户失败: %v", result.Error)
		common.FailWithMessage("登录失败", c)
		return
	}

	isNewUser := result.Error == gorm.ErrRecordNotFound

	if isNewUser {
		username := req.Nickname
		if username == "" {
			username = "wechat_" + wechatResp.OpenID[:8]
		}

		user = model.User{
			Username:   username,
			OpenID:     wechatResp.OpenID,
			SessionKey: wechatResp.SessionKey,
			Nickname:   req.Nickname,
			Avatar:     req.Avatar,
			WechatID:   wechatResp.UnionID,
		}

		if err := global.DB.Create(&user).Error; err != nil {
			global.SugarLogger.Errorf("创建微信用户失败: %v", err)
			common.FailWithMessage("登录失败", c)
			return
		}
	} else {
		updates := map[string]interface{}{
			"session_key": wechatResp.SessionKey,
		}
		if req.Nickname != "" {
			updates["nickname"] = req.Nickname
		}
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}
		if wechatResp.UnionID != "" {
			updates["wechat_id"] = wechatResp.UnionID
		}

		if err := global.DB.Model(&user).Updates(updates).Error; err != nil {
			global.SugarLogger.Errorf("更新微信用户失败: %v", err)
			common.FailWithMessage("登录失败", c)
			return
		}
		global.DB.First(&user, user.ID)
	}

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
