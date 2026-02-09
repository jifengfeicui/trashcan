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
		model.TrashCan{},
	)
	if err != nil {
		global.SugarLogger.Error("register table failed")
		os.Exit(0)
	}
	global.SugarLogger.Info("register table success")
}
