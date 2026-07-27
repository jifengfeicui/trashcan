package Orm

import (
	"encoding/json"
	"os"
	"strings"

	"gorm.io/gorm"

	"trashcan/ginServer/model"
	"trashcan/global"
)

func RegisterTables() {
	db := global.DB
	err := db.AutoMigrate(
		model.Test{},
		model.User{},
		model.Tag{},
		model.TrashCan{},
		model.TrashCanLike{},
	)
	if err != nil {
		global.SugarLogger.Error("register table failed")
		os.Exit(0)
	}
	global.SugarLogger.Info("register table success")

	// 迁移：确保 image_path_2 和 image_path_3 列存在
	if !db.Migrator().HasColumn(&model.TrashCan{}, "ImagePath2") {
		if err := db.Migrator().AddColumn(&model.TrashCan{}, "ImagePath2"); err != nil {
			global.SugarLogger.Warnf("迁移添加 ImagePath2 列失败: %v", err)
		} else {
			global.SugarLogger.Info("迁移添加 ImagePath2 列成功")
		}
	}
	if !db.Migrator().HasColumn(&model.TrashCan{}, "ImagePath3") {
		if err := db.Migrator().AddColumn(&model.TrashCan{}, "ImagePath3"); err != nil {
			global.SugarLogger.Warnf("迁移添加 ImagePath3 列失败: %v", err)
		} else {
			global.SugarLogger.Info("迁移添加 ImagePath3 列成功")
		}
	}

	// 迁移：确保 is_admin 列存在
	if !db.Migrator().HasColumn(&model.User{}, "IsAdmin") {
		if err := db.Migrator().AddColumn(&model.User{}, "IsAdmin"); err != nil {
			global.SugarLogger.Warnf("迁移添加 IsAdmin 列失败: %v", err)
		} else {
			global.SugarLogger.Info("迁移添加 IsAdmin 列成功")
		}
	}

	// 迁移：确保 tag_id 列存在
	if !db.Migrator().HasColumn(&model.TrashCan{}, "TagID") {
		if err := db.Migrator().AddColumn(&model.TrashCan{}, "TagID"); err != nil {
			global.SugarLogger.Warnf("迁移添加 TagID 列失败: %v", err)
		} else {
			global.SugarLogger.Info("迁移添加 TagID 列成功")
		}
	}

	// 迁移：将现有的 JSON tags 数据迁移到 tags 表
	migrateTagsToTable(db)

	// 设置默认管理员（用户名匹配时）
	var adminUser model.User
	if err := db.Where("username = ?", "wechat_oaJSr4vU").First(&adminUser).Error; err == nil {
		if !adminUser.IsAdmin {
			if err := db.Model(&adminUser).Update("is_admin", true).Error; err != nil {
				global.SugarLogger.Warnf("设置管理员失败: %v", err)
			} else {
				global.SugarLogger.Info("已设置 wechat_oaJSr4vU 为管理员")
			}
		}
	}
}

func migrateTagsToTable(db *gorm.DB) {
	type TrashCanWithTags struct {
		ID   uint
		Tags string
	}

	var trashCans []TrashCanWithTags
	if err := db.Model(&model.TrashCan{}).Where("tags IS NOT NULL AND tags != ''").Find(&trashCans).Error; err != nil {
		global.SugarLogger.Warnf("查询现有 tags 失败: %v", err)
		return
	}

	if len(trashCans) == 0 {
		return
	}

	tagMap := make(map[string]*model.Tag)

	for _, tc := range trashCans {
		if tc.Tags == "" {
			continue
		}
		var tagsArray []string
		if err := json.Unmarshal([]byte(tc.Tags), &tagsArray); err != nil || len(tagsArray) == 0 {
			continue
		}

		tagName := strings.TrimSpace(tagsArray[0])
		if tagName == "" {
			continue
		}

		var tag *model.Tag
		if existing, ok := tagMap[tagName]; ok {
			tag = existing
		} else {
			var foundTag model.Tag
			if err := db.Where("name = ?", tagName).First(&foundTag).Error; err != nil {
				tag = &model.Tag{Name: tagName}
				if err := db.Create(tag).Error; err != nil {
					global.SugarLogger.Warnf("创建标签失败: %v", err)
					continue
				}
			} else {
				tag = &foundTag
			}
			tagMap[tagName] = tag
		}

		if err := db.Model(&model.TrashCan{}).Where("id = ?", tc.ID).Update("tag_id", tag.ID).Error; err != nil {
			global.SugarLogger.Warnf("更新 trashcan tag_id 失败: %v", err)
		}
	}

	global.SugarLogger.Info("标签数据迁移完成")
}
