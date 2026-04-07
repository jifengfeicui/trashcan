package api

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"template/ginServer/api/common"
	"template/ginServer/model"
	"template/global"
	"template/utils"
)

// GetNearbyTrashCans 获取附近的垃圾桶
// GET /api/trashcans/nearby?lat=39.9&lng=116.4&radius=5&limit=10&tag_id=1
func GetNearbyTrashCans(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5")
	limitStr := c.DefaultQuery("limit", "10")
	tagIDStr := c.Query("tag_id")

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

	var tagID *uint
	if tagIDStr != "" {
		if id, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			uid := uint(id)
			tagID = &uid
		}
	}

	query := global.DB.Preload("User").Preload("Tag").Model(&model.TrashCan{})
	if tagID != nil {
		query = query.Where("tag_id = ?", *tagID)
	}

	var trashCans []model.TrashCan
	if err := query.Find(&trashCans).Error; err != nil {
		global.SugarLogger.Errorf("查询垃圾桶失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}

	var trashCanIDs []uint
	for _, tc := range trashCans {
		trashCanIDs = append(trashCanIDs, tc.ID)
	}

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

	type TrashCanWithDistance struct {
		model.TrashCan
		Distance       float64 `json:"distance"`
		ImageURL       string  `json:"image_url"`
		ImageURL2      string  `json:"image_url_2"`
		ImageURL3      string  `json:"image_url_3"`
		LikeCount      int64   `json:"like_count"`
		DislikeCount   int64   `json:"dislike_count"`
		UploaderName   string  `json:"uploader_name"`
		UploaderAvatar string  `json:"uploader_avatar"`
		TagName        string  `json:"tag_name"`
	}

	var results []TrashCanWithDistance
	for _, tc := range trashCans {
		distance := utils.CalculateDistance(lat, lng, tc.Latitude, tc.Longitude)
		if distance > radius {
			continue
		}

		uploaderName := ""
		uploaderAvatar := ""
		if tc.User != nil {
			uploaderName = tc.User.Nickname
			uploaderAvatar = tc.User.Avatar
		}

		tagName := ""
		if tc.Tag != nil {
			tagName = tc.Tag.Name
		}

		results = append(results, TrashCanWithDistance{
			TrashCan:       tc,
			Distance:       distance,
			ImageURL:       utils.GetImageURL(tc.ImagePath),
			ImageURL2:      utils.GetImageURL(tc.ImagePath2),
			ImageURL3:      utils.GetImageURL(tc.ImagePath3),
			LikeCount:      likeCounts[tc.ID],
			DislikeCount:   dislikeCounts[tc.ID],
			UploaderName:   uploaderName,
			UploaderAvatar: uploaderAvatar,
			TagName:        tagName,
		})
	}

	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Distance > results[j+1].Distance {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	if len(results) > limit {
		results = results[:limit]
	}

	common.OkWithData(results, c)
}

// CreateTrashCan 创建新垃圾桶
// POST /api/trashcans
func CreateTrashCan(c *gin.Context) {
	latStr := c.PostForm("latitude")
	lngStr := c.PostForm("longitude")
	address := c.PostForm("address")
	description := c.PostForm("description")
	tagIDStr := c.PostForm("tag_id")

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

	var tagID *uint
	if tagIDStr != "" {
		if id, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			uid := uint(id)
			tagID = &uid
		}
	}

	uploadDir := global.CONFIG.UploadConfig.ImageDir
	if uploadDir == "" {
		uploadDir = "uploads/trashcans"
	}

	if err := utils.EnsureUploadDir(uploadDir); err != nil {
		global.SugarLogger.Errorf("创建上传目录失败: %v", err)
		common.FailWithMessage("创建上传目录失败", c)
		return
	}

	var imagePath, imagePath2, imagePath3 string

	if file, err := c.FormFile("image"); err == nil {
		path, err := utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}
		imagePath = path
	} else {
		if imagePathStr := c.PostForm("image"); imagePathStr != "" {
			imagePath = imagePathStr
		}
	}

	if imagePath2Str := c.PostForm("image_2"); imagePath2Str != "" {
		imagePath2 = imagePath2Str
	}
	if imagePath3Str := c.PostForm("image_3"); imagePath3Str != "" {
		imagePath3 = imagePath3Str
	}

	for i := 0; i < 3; i++ {
		key := "images"
		if i > 0 {
			key = fmt.Sprintf("images[%d]", i)
		}
		file, err := c.FormFile(key)
		if err != nil {
			continue
		}

		path, err := utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}

		switch i {
		case 0:
			imagePath = path
		case 1:
			imagePath2 = path
		case 2:
			imagePath3 = path
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		common.FailWithAuthority(c)
		return
	}

	userIDUint := userID.(uint)

	trashCan := model.TrashCan{
		UserID:      &userIDUint,
		TagID:       tagID,
		Latitude:    lat,
		Longitude:   lng,
		Address:     address,
		Description: description,
		ImagePath:   imagePath,
		ImagePath2:  imagePath2,
		ImagePath3:  imagePath3,
	}

	if err := global.DB.Create(&trashCan).Error; err != nil {
		global.SugarLogger.Errorf("创建垃圾桶失败: %v", err)
		common.FailWithMessage("创建失败", c)
		return
	}

	result := map[string]interface{}{
		"id":          trashCan.ID,
		"latitude":    trashCan.Latitude,
		"longitude":   trashCan.Longitude,
		"address":     trashCan.Address,
		"image_url":   utils.GetImageURL(trashCan.ImagePath),
		"image_url_2": utils.GetImageURL(trashCan.ImagePath2),
		"image_url_3": utils.GetImageURL(trashCan.ImagePath3),
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
	if err := global.DB.Preload("Tag").First(&trashCan, id).Error; err != nil {
		common.FailWithMessage("垃圾桶不存在", c)
		return
	}

	var likeCount int64
	var dislikeCount int64
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", id, 1).
		Count(&likeCount)
	global.DB.Model(&model.TrashCanLike{}).
		Where("trash_can_id = ? AND type = ?", id, -1).
		Count(&dislikeCount)

	var userAction int8 = 0
	userID, exists := c.Get("userID")
	if exists {
		var like model.TrashCanLike
		if err := global.DB.Where("user_id = ? AND trash_can_id = ?", userID, id).First(&like).Error; err == nil {
			userAction = like.Type
		}
	}

	tagName := ""
	tagID := uint(0)
	if trashCan.Tag != nil {
		tagName = trashCan.Tag.Name
		tagID = trashCan.Tag.ID
	}

	result := map[string]interface{}{
		"id":            trashCan.ID,
		"latitude":      trashCan.Latitude,
		"longitude":     trashCan.Longitude,
		"address":       trashCan.Address,
		"description":   trashCan.Description,
		"tag_id":        tagID,
		"tag_name":      tagName,
		"image_url":     utils.GetImageURL(trashCan.ImagePath),
		"like_count":    likeCount,
		"dislike_count": dislikeCount,
		"user_action":   userAction,
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
		ImageURL2   string  `json:"image_url_2"`
		ImageURL3   string  `json:"image_url_3"`
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
			ImageURL2:   utils.GetImageURL(tc.ImagePath2),
			ImageURL3:   utils.GetImageURL(tc.ImagePath3),
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

	// 处理多图上传
	uploadDir := global.CONFIG.UploadConfig.ImageDir
	if uploadDir == "" {
		uploadDir = "uploads/trashcans"
	}

	if err := utils.EnsureUploadDir(uploadDir); err != nil {
		global.SugarLogger.Errorf("创建上传目录失败: %v", err)
		common.FailWithMessage("创建上传目录失败", c)
		return
	}

	// 支持前端上传的 image 字段，也支持 images[0], images[1], images[2]
	// 先尝试读取 image 字段（前端 MyTrashCans.vue 使用）
	if file, err := c.FormFile("image"); err == nil {
		path, err := utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}
		updates["image_path"] = path
	} else {
		// 小程序可能传图片路径字符串
		if imagePathStr := c.PostForm("image"); imagePathStr != "" {
			updates["image_path"] = imagePathStr
		}

		// 支持 image_2, image_3 路径字符串（包括空字符串以清除图片）
		if c.PostForm("image_2") != "" {
			updates["image_path2"] = c.PostForm("image_2")
		} else if _, exists := c.GetPostForm("image_2"); exists {
			updates["image_path2"] = ""
		}
		if c.PostForm("image_3") != "" {
			updates["image_path3"] = c.PostForm("image_3")
		} else if _, exists := c.GetPostForm("image_3"); exists {
			updates["image_path3"] = ""
		}
	}

	// 支持 image_2, image_3 路径字符串
	if imagePath2Str := c.PostForm("image_2"); imagePath2Str != "" {
		updates["image_path2"] = imagePath2Str
	}
	if imagePath3Str := c.PostForm("image_3"); imagePath3Str != "" {
		updates["image_path3"] = imagePath3Str
	}

	// 支持多图上传 images[0], images[1], images[2]
	for i := 0; i < 3; i++ {
		key := "images"
		if i > 0 {
			key = fmt.Sprintf("images[%d]", i)
		}
		file, err := c.FormFile(key)
		if err != nil {
			continue
		}

		path, err := utils.SaveImage(file, uploadDir)
		if err != nil {
			global.SugarLogger.Errorf("保存图片失败: %v", err)
			common.FailWithMessage("保存图片失败: "+err.Error(), c)
			return
		}

		switch i {
		case 0:
			updates["image_path"] = path
		case 1:
			updates["image_path2"] = path
		case 2:
			updates["image_path3"] = path
		}
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
		"image_url_2": utils.GetImageURL(trashCan.ImagePath2),
		"image_url_3": utils.GetImageURL(trashCan.ImagePath3),
		"updated_at":  trashCan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	common.OkWithDetailed(result, "更新成功", c)
}

// UploadImage 上传单张图片
// POST /api/upload/image
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		common.FailWithMessage("请选择图片", c)
		return
	}

	uploadDir := global.CONFIG.UploadConfig.ImageDir
	if uploadDir == "" {
		uploadDir = "uploads/trashcans"
	}

	if err := utils.EnsureUploadDir(uploadDir); err != nil {
		global.SugarLogger.Errorf("创建上传目录失败: %v", err)
		common.FailWithMessage("创建上传目录失败", c)
		return
	}

	imagePath, err := utils.SaveImage(file, uploadDir)
	if err != nil {
		global.SugarLogger.Errorf("保存图片失败: %v", err)
		common.FailWithMessage("保存图片失败: "+err.Error(), c)
		return
	}

	common.OkWithData(map[string]string{
		"image_path": imagePath,
	}, c)
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

// AdminGetAllTrashCans 获取所有垃圾桶（管理员）
// GET /api/admin/trashcans?page=1&page_size=10&keyword=xxx
func AdminGetAllTrashCans(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := global.DB.Preload("User").Preload("Tag").Model(&model.TrashCan{})

	if keyword != "" {
		query = query.Where("address LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var trashCans []model.TrashCan
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&trashCans)

	type AdminTrashCanItem struct {
		ID           uint    `json:"id"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Address      string  `json:"address"`
		Description  string  `json:"description"`
		TagID        *uint   `json:"tag_id"`
		TagName      string  `json:"tag_name"`
		ImageURL     string  `json:"image_url"`
		ImageURL2    string  `json:"image_url_2"`
		ImageURL3    string  `json:"image_url_3"`
		UserID       *uint   `json:"user_id"`
		UploaderName string  `json:"uploader_name"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	}

	var list []AdminTrashCanItem
	for _, tc := range trashCans {
		uploaderName := ""
		if tc.User != nil {
			uploaderName = tc.User.Nickname
		}
		tagName := ""
		if tc.Tag != nil {
			tagName = tc.Tag.Name
		}
		list = append(list, AdminTrashCanItem{
			ID:           tc.ID,
			Latitude:     tc.Latitude,
			Longitude:    tc.Longitude,
			Address:      tc.Address,
			Description:  tc.Description,
			TagID:        tc.TagID,
			TagName:      tagName,
			ImageURL:     utils.GetImageURL(tc.ImagePath),
			ImageURL2:    utils.GetImageURL(tc.ImagePath2),
			ImageURL3:    utils.GetImageURL(tc.ImagePath3),
			UserID:       tc.UserID,
			UploaderName: uploaderName,
			CreatedAt:    tc.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    tc.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	common.OkWithData(map[string]interface{}{
		"list":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}, c)
}

// AdminUpdateTrashCan 更新垃圾桶（管理员）
// PUT /api/admin/trashcans/:id
func AdminUpdateTrashCan(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

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

	address := c.PostForm("address")
	description := c.PostForm("description")
	tagIDStr := c.PostForm("tag_id")

	updates := map[string]interface{}{}
	if address != "" {
		updates["address"] = address
	}
	if description != "" {
		updates["description"] = description
	}
	if tagIDStr != "" {
		if tagID, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			if tagID == 0 {
				updates["tag_id"] = nil
			} else {
				uid := uint(tagID)
				updates["tag_id"] = &uid
			}
		}
	}

	if len(updates) > 0 {
		if err := global.DB.Model(&trashCan).Updates(updates).Error; err != nil {
			global.SugarLogger.Errorf("更新垃圾桶失败: %v", err)
			common.FailWithMessage("更新失败", c)
			return
		}
	}

	global.DB.Preload("Tag").First(&trashCan, id)

	tagName := ""
	if trashCan.Tag != nil {
		tagName = trashCan.Tag.Name
	}

	result := map[string]interface{}{
		"id":          trashCan.ID,
		"latitude":    trashCan.Latitude,
		"longitude":   trashCan.Longitude,
		"address":     trashCan.Address,
		"description": trashCan.Description,
		"tag_id":      trashCan.TagID,
		"tag_name":    tagName,
		"image_url":   utils.GetImageURL(trashCan.ImagePath),
		"updated_at":  trashCan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	common.OkWithDetailed(result, "更新成功", c)
}

// AdminDeleteTrashCan 删除垃圾桶（管理员）
// DELETE /api/admin/trashcans/:id
func AdminDeleteTrashCan(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

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

	if trashCan.ImagePath != "" {
		os.Remove(trashCan.ImagePath)
	}

	global.DB.Delete(&trashCan)
	global.DB.Where("trash_can_id = ?", id).Delete(&model.TrashCanLike{})

	common.OkWithMessage("删除成功", c)
}

func getAllTagsInternal() ([]model.Tag, error) {
	var tags []model.Tag
	if err := global.DB.Order("name").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAllTags 获取所有标签
// GET /api/tags
func GetAllTags(c *gin.Context) {
	tags, err := getAllTagsInternal()
	if err != nil {
		global.SugarLogger.Errorf("查询标签失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}
	common.OkWithData(tags, c)
}

// AdminGetTagsWithCount 获取所有标签及关联地点数量（管理员）
// GET /api/admin/tags
func AdminGetTagsWithCount(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

	var tags []model.Tag
	if err := global.DB.Order("name").Find(&tags).Error; err != nil {
		global.SugarLogger.Errorf("查询标签失败: %v", err)
		common.FailWithMessage("查询失败", c)
		return
	}

	type TagItem struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}

	var tagItems []TagItem
	for _, tag := range tags {
		var count int64
		global.DB.Model(&model.TrashCan{}).Where("tag_id = ?", tag.ID).Count(&count)
		tagItems = append(tagItems, TagItem{ID: tag.ID, Name: tag.Name, Count: count})
	}

	common.OkWithData(tagItems, c)
}

// AdminDeleteTag 删除指定标签（管理员）
// DELETE /api/admin/tags?id=xxx
func AdminDeleteTag(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

	tagIDStr := c.Query("id")
	if tagIDStr == "" {
		common.ParamError(c)
		return
	}

	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		common.ParamError(c)
		return
	}

	var tag model.Tag
	if err := global.DB.First(&tag, uint(tagID)).Error; err != nil {
		common.FailWithMessage("标签不存在", c)
		return
	}

	global.DB.Model(&model.TrashCan{}).Where("tag_id = ?", tagID).Update("tag_id", nil)

	if err := global.DB.Delete(&tag).Error; err != nil {
		global.SugarLogger.Errorf("删除标签失败: %v", err)
		common.FailWithMessage("删除失败", c)
		return
	}

	common.OkWithMessage("删除成功", c)
}

// AdminCreateTag 创建标签（管理员）
// POST /api/admin/tags
func AdminCreateTag(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		common.FailWithMessage("权限不足", c)
		return
	}

	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		common.ParamError(c)
		return
	}

	var existing model.Tag
	if err := global.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		common.FailWithMessage("标签已存在", c)
		return
	}

	tag := model.Tag{Name: name}
	if err := global.DB.Create(&tag).Error; err != nil {
		global.SugarLogger.Errorf("创建标签失败: %v", err)
		common.FailWithMessage("创建失败", c)
		return
	}

	common.OkWithData(tag, c)
}
