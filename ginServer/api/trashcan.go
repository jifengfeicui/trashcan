package api

import (
	"os"
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

// GetUserTrashCans 获取当前用户上传的垃圾桶列表（分页）
// GET /api/users/me/trashcans?page=1&page_size=10
func GetUserTrashCans(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}
	userIDUint := userID.(uint)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大每页数量
	}

	// 查询总数
	var total int64
	if err := global.DB.Model(&model.TrashCan{}).
		Where("user_id = ?", userIDUint).
		Count(&total).Error; err != nil {
		global.SugarLogger.Errorf("查询垃圾桶总数失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	// 查询列表数据
	var trashCans []model.TrashCan
	offset := (page - 1) * pageSize
	if err := global.DB.Where("user_id = ?", userIDUint).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&trashCans).Error; err != nil {
		global.SugarLogger.Errorf("查询垃圾桶列表失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}

	// 构建返回数据
	type TrashCanItem struct {
		ID          uint    `json:"id"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Address     string  `json:"address"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
	}

	var list []TrashCanItem
	for _, tc := range trashCans {
		list = append(list, TrashCanItem{
			ID:          tc.ID,
			Latitude:    tc.Latitude,
			Longitude:   tc.Longitude,
			Address:     tc.Address,
			Description: tc.Description,
			ImageURL:    utils.GetImageURL(tc.ImagePath),
			CreatedAt:   tc.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tc.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	result := map[string]interface{}{
		"list":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}

	common.OkWithData(result, c)
}

// UpdateTrashCan 更新垃圾桶信息
// PUT /api/trashcans/:id
func UpdateTrashCan(c *gin.Context) {
	// 获取垃圾桶ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ParamError(c)
		return
	}

	// 从中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}
	userIDUint := userID.(uint)

	// 查询垃圾桶是否存在且属于当前用户
	var trashCan model.TrashCan
	if err := global.DB.Where("id = ? AND user_id = ?", id, userIDUint).First(&trashCan).Error; err != nil {
		common.FailWithMessage("垃圾桶不存在或无权限", c)
		return
	}

	// 解析表单数据
	address := c.PostForm("address")
	description := c.PostForm("description")

	// 更新地址和描述
	updates := map[string]interface{}{
		"address":     address,
		"description": description,
	}

	// 处理图片上传（如果提供了新图片）
	file, err := c.FormFile("image")
	if err == nil {
		// 有图片上传，先删除旧图片
		if trashCan.ImagePath != "" {
			if err := os.Remove(trashCan.ImagePath); err != nil {
				global.SugarLogger.Warnf("删除旧图片失败: %v", err)
				// 继续执行，不中断流程
			}
		}

		// 保存新图片
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
		imagePath, err := utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}

		updates["image_path"] = imagePath
	}

	// 更新数据库
	if err := global.DB.Model(&trashCan).Updates(updates).Error; err != nil {
		global.SugarLogger.Errorf("更新垃圾桶失败: %v", err)
		common.FailWithMessage("更新失败", c)
		return
	}

	// 重新查询以获取更新后的数据
	global.DB.First(&trashCan, id)

	// 返回更新结果
	result := map[string]interface{}{
		"id":          trashCan.ID,
		"latitude":    trashCan.Latitude,
		"longitude":   trashCan.Longitude,
		"address":     trashCan.Address,
		"description": trashCan.Description,
		"image_url":   utils.GetImageURL(trashCan.ImagePath),
		"updated_at":  trashCan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	common.OkWithDetailed(result, "更新成功", c)
}

// DeleteTrashCan 删除垃圾桶
// DELETE /api/trashcans/:id
func DeleteTrashCan(c *gin.Context) {
	// 获取垃圾桶ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ParamError(c)
		return
	}

	// 从中间件获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}
	userIDUint := userID.(uint)

	// 查询垃圾桶是否存在且属于当前用户
	var trashCan model.TrashCan
	if err := global.DB.Where("id = ? AND user_id = ?", id, userIDUint).First(&trashCan).Error; err != nil {
		common.FailWithMessage("垃圾桶不存在或无权限", c)
		return
	}

	// 删除关联的图片文件
	if trashCan.ImagePath != "" {
		if err := os.Remove(trashCan.ImagePath); err != nil {
			global.SugarLogger.Warnf("删除图片文件失败: %v", err)
			// 继续执行，不中断流程
		}
	}

	// 删除数据库记录
	if err := global.DB.Delete(&trashCan).Error; err != nil {
		global.SugarLogger.Errorf("删除垃圾桶失败: %v", err)
		common.FailWithMessage("删除失败", c)
		return
	}

	common.OkWithMessage("删除成功", c)
}
