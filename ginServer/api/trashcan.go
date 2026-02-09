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

	// 计算距离并筛选
	type TrashCanWithDistance struct {
		model.TrashCan
		Distance float64 `json:"distance"`  // 距离（公里）
		ImageURL string  `json:"image_url"` // 图片URL
	}

	var results []TrashCanWithDistance
	for _, tc := range trashCans {
		distance := utils.CalculateDistance(lat, lng, tc.Latitude, tc.Longitude)
		if distance <= radius {
			results = append(results, TrashCanWithDistance{
				TrashCan: tc,
				Distance: distance,
				ImageURL: utils.GetImageURL(tc.ImagePath),
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

	// 创建垃圾桶记录
	trashCan := model.TrashCan{
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

	// 构建返回数据
	result := map[string]interface{}{
		"id":          trashCan.ID,
		"latitude":    trashCan.Latitude,
		"longitude":   trashCan.Longitude,
		"address":     trashCan.Address,
		"description": trashCan.Description,
		"image_url":   utils.GetImageURL(trashCan.ImagePath),
		"created_at":  trashCan.CreatedAt,
		"updated_at":  trashCan.UpdatedAt,
	}

	common.OkWithData(result, c)
}
