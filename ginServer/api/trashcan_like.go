package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"template/ginServer/api/common"
	"template/ginServer/model"
	"template/global"
)

// ToggleLike 点赞
// POST /api/trashcans/:id/like
func ToggleLike(c *gin.Context) {
	toggleLikeOrDislike(c, 1)
}

// ToggleDislike 点踩
// POST /api/trashcans/:id/dislike
func ToggleDislike(c *gin.Context) {
	toggleLikeOrDislike(c, -1)
}

// toggleLikeOrDislike 切换点赞/点踩状态
func toggleLikeOrDislike(c *gin.Context, likeType int8) {
	// 获取垃圾桶ID
	idStr := c.Param("id")
	trashCanID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ParamError(c)
		return
	}

	// 验证垃圾桶是否存在
	var trashCan model.TrashCan
	if err := global.DB.First(&trashCan, trashCanID).Error; err != nil {
		common.FailWithMessage("垃圾桶不存在", c)
		return
	}

	// 从中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}
	userIDUint := userID.(uint)

	// 查找用户是否已经对该垃圾桶操作过
	var existingLike model.TrashCanLike
	err = global.DB.Where("user_id = ? AND trash_can_id = ?", userIDUint, trashCanID).First(&existingLike).Error

	var userAction int8 = 0 // 0=未操作, 1=点赞, -1=点踩

	if err != nil {
		// 不存在，创建新记录
		like := model.TrashCanLike{
			UserID:     userIDUint,
			TrashCanID: uint(trashCanID),
			Type:       likeType,
		}
		if err := global.DB.Create(&like).Error; err != nil {
			global.SugarLogger.Errorf("创建点赞记录失败: %v", err)
			common.FailWithMessage("操作失败", c)
			return
		}
		userAction = likeType
	} else {
		// 已存在
		if existingLike.Type == likeType {
			// 相同操作，则取消（删除记录）
			if err := global.DB.Delete(&existingLike).Error; err != nil {
				global.SugarLogger.Errorf("删除点赞记录失败: %v", err)
				common.FailWithMessage("操作失败", c)
				return
			}
			userAction = 0 // 取消操作
		} else {
			// 不同操作，切换类型
			existingLike.Type = likeType
			if err := global.DB.Save(&existingLike).Error; err != nil {
				global.SugarLogger.Errorf("更新点赞记录失败: %v", err)
				common.FailWithMessage("操作失败", c)
				return
			}
			userAction = likeType
		}
	}

	// 重新统计点赞和点踩数量
	var likeCount int64
	var dislikeCount int64
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", trashCanID, 1).
		Count(&likeCount)
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", trashCanID, -1).
		Count(&dislikeCount)

	result := map[string]interface{}{
		"like_count":    likeCount,
		"dislike_count": dislikeCount,
		"user_action":   userAction, // 当前用户的操作：0=未操作, 1=点赞, -1=点踩
	}

	common.OkWithData(result, c)
}

