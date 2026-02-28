package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

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
		"created_at": user.CreatedAt,
	}

	common.OkWithData(result, c)
}
