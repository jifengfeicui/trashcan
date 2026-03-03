package Orm

import (
	"os"
	"template/ginServer/model"
	"template/global"
)

func RegisterTables() {
	db := global.DB
	err := db.AutoMigrate(
		model.Test{},
		model.User{},
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
}
