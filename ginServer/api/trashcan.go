package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"template/ginServer/api/common"
	"template/ginServer/model"
	"template/global"
	"template/utils"
)

// GetNearbyTrashCans 获取附近的垃圾桶
// GET /api/trashcans/nearby?lat=39.9&lng=116.4&radius=5&limit=10
func GetNearbyTrashCans(c *gin.Context) {
	// 获取查询参数
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5") // 默认搜索半径5公里
	limitStr := c.DefaultQuery("limit", "10")  // 默认返回10个

	// 解析参数
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		common.ParamError(c)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		common.ParamError(c)
		return
	}

	radius, _ := strconv.ParseFloat(radiusStr, 64)
	limit, _ := strconv.Atoi(limitStr)

	// 查询所有垃圾桶
	var trashCans []model.TrashCan
	if err := global.DB.Find(&trashCans).Error; err != nil {
		global.SugarLogger.Errorf("查询垃圾桶失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}

	// 获取所有垃圾桶的ID列表
	var trashCanIDs []uint
	for _, tc := range trashCans {
		trashCanIDs = append(trashCanIDs, tc.ID)
	}

	// 统计每个垃圾桶的点赞和点踩数量
	likeCounts := make(map[uint]int64)
	dislikeCounts := make(map[uint]int64)
	if len(trashCanIDs) > 0 {
		var likes []struct {
			TrashCanID uint
			Count      int64
		}
		global.DB.Model(&model.TrashCanLike{}).
			Where("trash_can_id IN ? AND type = ?", trashCanIDs, 1).
			Select("trash_can_id, COUNT(*) as count").
			Group("trash_can_id").
			Scan(&likes)
		for _, l := range likes {
			likeCounts[l.TrashCanID] = l.Count
		}

		var dislikes []struct {
			TrashCanID uint
			Count      int64
		}
		global.DB.Model(&model.TrashCanLike{}).
			Where("trash_can_id IN ? AND type = ?", trashCanIDs, -1).
			Select("trash_can_id, COUNT(*) as count").
			Group("trash_can_id").
			Scan(&dislikes)
		for _, d := range dislikes {
			dislikeCounts[d.TrashCanID] = d.Count
		}
	}

	// 计算距离并筛选
	type TrashCanWithDistance struct {
		model.TrashCan
		Distance     float64 `json:"distance"`      // 距离（公里）
		ImageURL     string  `json:"image_url"`     // 图片URL
		LikeCount    int64   `json:"like_count"`    // 点赞数
		DislikeCount int64   `json:"dislike_count"` // 点踩数
	}

	var results []TrashCanWithDistance
	for _, tc := range trashCans {
		distance := utils.CalculateDistance(lat, lng, tc.Latitude, tc.Longitude)
		if distance <= radius {
			results = append(results, TrashCanWithDistance{
				TrashCan:     tc,
				Distance:     distance,
				ImageURL:     utils.GetImageURL(tc.ImagePath),
				LikeCount:    likeCounts[tc.ID],
				DislikeCount: dislikeCounts[tc.ID],
			})
		}
	}

	// 按距离排序（冒泡排序，简单实现）
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Distance > results[j+1].Distance {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// 限制返回数量
	if len(results) > limit {
		results = results[:limit]
	}

	common.OkWithData(results, c)
}

// CreateTrashCan 创建新垃圾桶
// POST /api/trashcans
func CreateTrashCan(c *gin.Context) {
	// 解析表单数据
	latStr := c.PostForm("latitude")
	lngStr := c.PostForm("longitude")
	address := c.PostForm("address")
	description := c.PostForm("description")

	// 验证必填参数
	if latStr == "" || lngStr == "" {
		common.ParamError(c)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		common.ParamError(c)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		common.ParamError(c)
		return
	}

	// 处理图片上传
	var imagePath string
	file, err := c.FormFile("image")
	if err == nil {
		// 有图片上传
		uploadDir := global.CONFIG.UploadConfig.ImageDir
		if uploadDir == "" {
			uploadDir = "uploads/trashcans"
		}

		// 确保上传目录存在
		if err := utils.EnsureUploadDir(uploadDir); err != nil {
			global.SugarLogger.Errorf("创建上传目录失败: %v", err)
			common.FailWithMessage("创建上传目录失败", c)
			return
		}

		// 保存图片
		imagePath, err = utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}
	}

	// 从中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}

	userIDUint := userID.(uint)

	// 创建垃圾桶记录
	trashCan := model.TrashCan{
		UserID:      &userIDUint,
		Latitude:    lat,
		Longitude:   lng,
		Address:     address,
		Description: description,
		ImagePath:   imagePath,
	}

	if err := global.DB.Create(&trashCan).Error; err != nil {
		global.SugarLogger.Errorf("创建垃圾桶失败: %v", err)
		common.FailWithMessage("创建失败", c)
		return
	}

	// 返回创建结果
	result := map[string]interface{}{
		"id":        trashCan.ID,
		"latitude":  trashCan.Latitude,
		"longitude": trashCan.Longitude,
		"address":   trashCan.Address,
		"image_url": utils.GetImageURL(trashCan.ImagePath),
	}

	common.OkWithDetailed(result, "创建成功", c)
}

// GetTrashCanDetail 获取垃圾桶详情
// GET /api/trashcans/:id
func GetTrashCanDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ParamError(c)
		return
	}

	var trashCan model.TrashCan
	if err := global.DB.First(&trashCan, id).Error; err != nil {
		common.FailWithMessage("垃圾桶不存在", c)
		return
	}

	// 统计点赞和点踩数量
	var likeCount int64
	var dislikeCount int64
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", id, 1).
		Count(&likeCount)
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", id, -1).
		Count(&dislikeCount)

	// 获取当前用户的操作状态（如果已登录）
	var userAction int8 = 0 // 0=未操作, 1=点赞, -1=点踩
	userID, exists := c.Get("userID")
	if exists {
		var like model.TrashCanLike
		if err := global.DB.Where("user_id = ? AND trash_can_id = ?", userID, id).First(&like).Error; err == nil {
			userAction = like.Type
		}
	}

	// 构建返回数据
	result := map[string]interface{}{
		"id":            trashCan.ID,
		"latitude":      trashCan.Latitude,
		"longitude":     trashCan.Longitude,
		"address":       trashCan.Address,
		"description":   trashCan.Description,
		"image_url":     utils.GetImageURL(trashCan.ImagePath),
		"like_count":    likeCount,
		"dislike_count": dislikeCount,
		"user_action":   userAction, // 当前用户的操作：0=未操作, 1=点赞, -1=点踩
		"created_at":    trashCan.CreatedAt,
		"updated_at":    trashCan.UpdatedAt,
	}

	common.OkWithData(result, c)
}
