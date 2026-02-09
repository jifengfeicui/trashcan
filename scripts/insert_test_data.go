package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"template/ginServer/model"
	"template/utils"
)

func main() {
	fmt.Println("🚀 开始插入测试数据...")

	// 连接数据库
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("无法连接到数据库: " + err.Error())
	}

	// 清空现有数据
	fmt.Println("🗑️  清空现有数据...")
	db.Exec("DELETE FROM trash_cans")
	fmt.Println("✅ 数据已清空")

	// 用户位置（上海）
	userLat := 31.19322644453637
	userLng := 121.41182831455195

	// 测试数据：在用户位置周围生成几个垃圾桶
	testData := []struct {
		lat         float64
		lng         float64
		address     string
		description string
	}{
		// 距离用户位置约100-200米的位置
		{31.1940, 121.4125, "上海市黄浦区南京东路步行街入口", "商业区垃圾桶，人流量大"},
		{31.1925, 121.4110, "上海市黄浦区人民广场地铁站出口", "地铁站附近垃圾桶"},
		{31.1935, 121.4105, "上海市黄浦区外滩观景台", "旅游景点垃圾桶"},
		{31.1920, 121.4130, "上海市黄浦区豫园商城", "商业区垃圾桶"},
		{31.1945, 121.4115, "上海市黄浦区南京路步行街中段", "步行街垃圾桶"},
		{31.1915, 121.4120, "上海市黄浦区城隍庙附近", "景区垃圾桶"},
		{31.1930, 121.4100, "上海市黄浦区外滩附近", "旅游区垃圾桶"},
		{31.1928, 121.4135, "上海市黄浦区人民广场", "广场垃圾桶"},
	}

	// 插入数据
	fmt.Println("\n📝 开始插入测试数据...")
	successCount := 0
	for i, data := range testData {
		trashCan := model.TrashCan{
			Latitude:    data.lat,
			Longitude:   data.lng,
			Address:     data.address,
			Description: data.description,
			ImagePath:   "", // 可以后续添加图片路径
		}

		// 计算距离
		distance := utils.CalculateDistance(userLat, userLng, data.lat, data.lng)

		if err := db.Create(&trashCan).Error; err != nil {
			fmt.Printf("❌ 插入第 %d 条数据失败: %v\n", i+1, err)
		} else {
			successCount++
			fmt.Printf("✅ 插入成功 #%d: %s (距离: %.2f 公里)\n",
				trashCan.ID, data.address, distance)
		}
	}

	fmt.Printf("\n📊 总共插入 %d/%d 条测试数据\n", successCount, len(testData))
	fmt.Printf("📍 用户位置: %.8f, %.8f\n", userLat, userLng)
	fmt.Println("\n✅ 测试数据已成功插入数据库！")
	fmt.Println("\n💡 提示：现在可以在前端搜索附近的垃圾桶来查看这些测试数据了。")
}
